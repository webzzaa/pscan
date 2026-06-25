package common

import (
	"sync"
	"testing"
)

/*
state_test.go - State 并发安全测试

测试重点：
1. 并发安全性 - 多goroutine同时操作计数器
2. 原子操作一致性 - 增减计数正确
3. Reset功能 - 重置后计数器归零

不测试：
- 限速器（需要复杂的时间模拟）
- 简单getter/setter
*/

// TestState_ConcurrentPacketCount 测试并发包计数
func TestState_ConcurrentPacketCount(t *testing.T) {
	s := NewState()

	const goroutines = 100
	const incrementsPerGoroutine = 1000

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				s.IncrementPacketCount()
			}
		}()
	}

	wg.Wait()

	expected := int64(goroutines * incrementsPerGoroutine)
	actual := s.GetPacketCount()

	if actual != expected {
		t.Errorf("并发计数不一致: 期望 %d, 实际 %d", expected, actual)
	}
}

// TestState_ConcurrentTCPCount 测试并发TCP计数
func TestState_ConcurrentTCPCount(t *testing.T) {
	s := NewState()

	const goroutines = 50
	const operationsPerGoroutine = 500

	var wg sync.WaitGroup
	wg.Add(goroutines * 2) // 成功和失败各一半

	// 成功连接
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				s.IncrementTCPSuccessPacketCount()
			}
		}()
	}

	// 失败连接
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				s.IncrementTCPFailedPacketCount()
			}
		}()
	}

	wg.Wait()

	expectedTotal := int64(goroutines * operationsPerGoroutine * 2)
	expectedSuccess := int64(goroutines * operationsPerGoroutine)
	expectedFailed := int64(goroutines * operationsPerGoroutine)

	if s.GetPacketCount() != expectedTotal {
		t.Errorf("总包计数不一致: 期望 %d, 实际 %d", expectedTotal, s.GetPacketCount())
	}
	if s.GetTCPPacketCount() != expectedTotal {
		t.Errorf("TCP包计数不一致: 期望 %d, 实际 %d", expectedTotal, s.GetTCPPacketCount())
	}
	if s.GetTCPSuccessPacketCount() != expectedSuccess {
		t.Errorf("TCP成功计数不一致: 期望 %d, 实际 %d", expectedSuccess, s.GetTCPSuccessPacketCount())
	}
	if s.GetTCPFailedPacketCount() != expectedFailed {
		t.Errorf("TCP失败计数不一致: 期望 %d, 实际 %d", expectedFailed, s.GetTCPFailedPacketCount())
	}
}

// TestState_Reset 测试重置功能
func TestState_Reset(t *testing.T) {
	s := NewState()

	// 增加一些计数
	for i := 0; i < 100; i++ {
		s.IncrementTCPSuccessPacketCount()
		s.IncrementTCPFailedPacketCount()
		s.IncrementUDPPacketCount()
		s.IncrementHTTPPacketCount()
		s.IncrementResourceExhaustedCount()
	}

	// 验证有值
	if s.GetPacketCount() == 0 {
		t.Fatal("重置前计数应该非零")
	}

	// 重置
	s.ResetPacketCounters()

	// 验证全部归零
	if s.GetPacketCount() != 0 {
		t.Errorf("重置后PacketCount应该为0, 实际 %d", s.GetPacketCount())
	}
	if s.GetTCPPacketCount() != 0 {
		t.Errorf("重置后TCPPacketCount应该为0, 实际 %d", s.GetTCPPacketCount())
	}
	if s.GetTCPSuccessPacketCount() != 0 {
		t.Errorf("重置后TCPSuccessPacketCount应该为0, 实际 %d", s.GetTCPSuccessPacketCount())
	}
	if s.GetTCPFailedPacketCount() != 0 {
		t.Errorf("重置后TCPFailedPacketCount应该为0, 实际 %d", s.GetTCPFailedPacketCount())
	}
	if s.GetUDPPacketCount() != 0 {
		t.Errorf("重置后UDPPacketCount应该为0, 实际 %d", s.GetUDPPacketCount())
	}
	if s.GetHTTPPacketCount() != 0 {
		t.Errorf("重置后HTTPPacketCount应该为0, 实际 %d", s.GetHTTPPacketCount())
	}
	if s.GetResourceExhaustedCount() != 0 {
		t.Errorf("重置后ResourceExhaustedCount应该为0, 实际 %d", s.GetResourceExhaustedCount())
	}
}

// TestState_TaskCounters 测试任务计数器
func TestState_TaskCounters(t *testing.T) {
	s := NewState()

	// 初始值应该为0
	if s.GetEnd() != 0 || s.GetNum() != 0 {
		t.Error("初始任务计数器应该为0")
	}

	// 设置值
	s.SetEnd(100)
	s.SetNum(50)

	if s.GetEnd() != 100 {
		t.Errorf("End应该为100, 实际 %d", s.GetEnd())
	}
	if s.GetNum() != 50 {
		t.Errorf("Num应该为50, 实际 %d", s.GetNum())
	}

	// 增加值
	s.IncrementEnd()
	s.IncrementNum()

	if s.GetEnd() != 101 {
		t.Errorf("IncrementEnd后应该为101, 实际 %d", s.GetEnd())
	}
	if s.GetNum() != 51 {
		t.Errorf("IncrementNum后应该为51, 实际 %d", s.GetNum())
	}
}

// TestState_ConcurrentTaskCounters 测试并发任务计数
func TestState_ConcurrentTaskCounters(t *testing.T) {
	s := NewState()

	const goroutines = 100
	const incrementsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	// 并发增加End
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				s.IncrementEnd()
			}
		}()
	}

	// 并发增加Num
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				s.IncrementNum()
			}
		}()
	}

	wg.Wait()

	expected := int64(goroutines * incrementsPerGoroutine)
	if s.GetEnd() != expected {
		t.Errorf("End并发计数不一致: 期望 %d, 实际 %d", expected, s.GetEnd())
	}
	if s.GetNum() != expected {
		t.Errorf("Num并发计数不一致: 期望 %d, 实际 %d", expected, s.GetNum())
	}
}

// TestState_OutputMutex 测试输出互斥锁
func TestState_OutputMutex(t *testing.T) {
	s := NewState()

	counter := 0
	const goroutines = 100
	const incrementsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				s.LockOutput()
				counter++
				s.UnlockOutput()
			}
		}()
	}

	wg.Wait()

	expected := goroutines * incrementsPerGoroutine
	if counter != expected {
		t.Errorf("输出互斥锁保护失败: 期望 %d, 实际 %d", expected, counter)
	}
}
