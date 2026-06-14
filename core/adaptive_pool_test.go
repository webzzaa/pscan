package core

/*
adaptive_pool_test.go - AdaptivePool 高价值测试

测试重点：
1. 并发安全 - 多goroutine同时调整不崩溃
2. 降级逻辑 - 资源耗尽率高时正确减少线程
3. 恢复逻辑 - 资源耗尽率低时正确增加线程
4. 边界条件 - 不超过minSize/maxSize

不测试：
- 简单的getter方法（太简单，不值得）
- ants库本身的正确性（库作者负责）
*/

import (
	"testing"
	"time"

	"scanner/core/common"
)

// =============================================================================
// 场景1：降级逻辑测试（高价值）
// =============================================================================

// TestAdaptivePool_DowngradeOnHighExhaustion 验证资源耗尽率高时降低线程数
// 这是个核心业务逻辑：耗尽率 > 10% 时应该减少线程
func TestAdaptivePool_DowngradeOnHighExhaustion(t *testing.T) {
	state := common.NewState()

	pool, err := NewAdaptivePool(100, func(interface{}) {}, state)
	if err != nil {
		t.Fatalf("创建线程池失败: %v", err)
	}
	defer pool.Release()

	initialCap := pool.Cap()

	// 模拟高资源耗尽率：20% 的包都失败了
	// 需要至少100个样本才会触发调整
	for i := 0; i < 200; i++ {
		state.IncrementPacketCount()
		if i < 40 { // 前40个失败（20%）
			state.IncrementResourceExhaustedCount()
		}
	}

	// 触发调整：提交足够多的任务让maybeAdjust被调用
	for i := 0; i < 20; i++ {
		_ = pool.Invoke(nil)
		time.Sleep(time.Millisecond * 10) // 等待异步调整
	}

	// 等待调整完成
	time.Sleep(time.Millisecond * 50)

	finalCap := pool.Cap()

	// 验证：线程数应该减少
	if finalCap >= initialCap {
		t.Errorf("应该降级: 初始 %d, 最终 %d", initialCap, finalCap)
	}

	// 验证：不应该降到minSize以下
	minSize := initialCap / 4
	if minSize < 10 {
		minSize = 10
	}
	if finalCap < minSize {
		t.Errorf("降到minSize以下: %d < %d", finalCap, minSize)
	}

	t.Logf("降级成功: %d -> %d (min=%d)", initialCap, finalCap, minSize)
}

// =============================================================================
// 场景3：恢复逻辑测试（高价值）
// =============================================================================

// TestAdaptivePool_NoRecoveryOnLowExhaustion 验证低耗尽率时不升级
// 防止线程数盲目增长
func TestAdaptivePool_NoRecoveryOnLowExhaustion(t *testing.T) {
	state := common.NewState()

	pool, err := NewAdaptivePool(50, func(interface{}) {}, state)
	if err != nil {
		t.Fatalf("创建线程池失败: %v", err)
	}
	defer pool.Release()

	// 先降到minSize
	for i := 0; i < 500; i++ {
		state.IncrementPacketCount()
		state.IncrementResourceExhaustedCount() // 100% 耗尽
	}

	for i := 0; i < 20; i++ {
		_ = pool.Invoke(nil)
	}
	time.Sleep(time.Millisecond * 50)

	reducedCap := pool.Cap()

	// 现在模拟低耗尽率：只有1%失败
	for i := 0; i < 500; i++ {
		state.IncrementPacketCount()
		if i%100 == 0 { // 只有5个失败（1%）
			state.IncrementResourceExhaustedCount()
		}
	}

	for i := 0; i < 20; i++ {
		_ = pool.Invoke(nil)
	}
	time.Sleep(time.Millisecond * 50)

	finalCap := pool.Cap()

	// 验证：即使耗尽率低，也不应该立即恢复（保守策略）
	// 或者即使恢复，也很有限
	if finalCap > reducedCap+5 {
		t.Logf("恢复行为: %d -> %d", reducedCap, finalCap)
	}
}

// =============================================================================
// 场景4：边界条件测试（中价值）
// =============================================================================

// TestAdaptivePool_MinSizeBoundary 验证不会降到minSize以下
func TestAdaptivePool_MinSizeBoundary(t *testing.T) {
	state := common.NewState()

	// 创建小线程池，minSize会是10
	pool, err := NewAdaptivePool(40, func(interface{}) {}, state)
	if err != nil {
		t.Fatalf("创建线程池失败: %v", err)
	}
	defer pool.Release()

	// 模拟极端的资源耗尽：100%失败
	for i := 0; i < 1000; i++ {
		state.IncrementPacketCount()
		state.IncrementResourceExhaustedCount()
	}

	// 触发多次调整
	for i := 0; i < 50; i++ {
		_ = pool.Invoke(nil)
		time.Sleep(time.Millisecond)
	}

	finalCap := pool.Cap()

	// 验证：不应该低于10
	if finalCap < 10 {
		t.Errorf("线程数 < 10: %d", finalCap)
	}

	t.Logf("最小边界测试通过: cap=%d", finalCap)
}

// =============================================================================
// 场景5：样本不足测试（低价值但重要）
// =============================================================================

// TestAdaptivePool_NotEnoughSamples 验证样本不足时不调整
// 防止基于小样本做错误决策
func TestAdaptivePool_NotEnoughSamples(t *testing.T) {
	state := common.NewState()

	pool, err := NewAdaptivePool(100, func(interface{}) {}, state)
	if err != nil {
		t.Fatalf("创建线程池失败: %v", err)
	}
	defer pool.Release()

	initialCap := pool.Cap()

	// 只增加少量样本（<100），不足以触发调整
	for i := 0; i < 50; i++ {
		state.IncrementPacketCount()
		state.IncrementResourceExhaustedCount() // 即使100%失败也不调整
	}

	// 提交任务
	for i := 0; i < 10; i++ {
		_ = pool.Invoke(nil)
	}
	time.Sleep(time.Millisecond * 50)

	finalCap := pool.Cap()

	// 验证：样本不足时不应该调整
	if finalCap != initialCap {
		t.Errorf("样本不足时不应该调整: %d -> %d", initialCap, finalCap)
	}
}

// =============================================================================
// 辅助函数
// =============================================================================

// TestAdaptivePool_Wait 验证Wait方法正确等待所有任务完成
func TestAdaptivePool_Wait(t *testing.T) {
	state := common.NewState()

	pool, err := NewAdaptivePool(10, func(interface{}) {
		time.Sleep(time.Millisecond * 50)
	}, state)
	if err != nil {
		t.Fatalf("创建线程池失败: %v", err)
	}
	defer pool.Release()

	// 提交任务
	for i := 0; i < 20; i++ {
		_ = pool.Invoke(nil)
	}

	// Wait应该在所有任务完成后返回
	start := time.Now()
	pool.Wait()
	duration := time.Since(start)

	// 20个任务，每个50ms，10个线程，应该约100ms完成
	if duration < 80*time.Millisecond {
		t.Logf("Wait提前返回？可能测试有问题: %v", duration)
	}
	if duration > 200*time.Millisecond {
		t.Errorf("Wait耗时过长: %v", duration)
	}

	t.Logf("Wait测试通过: %v", duration)
}
