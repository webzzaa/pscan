package core

import (
	"net"
	"testing"
	"time"
)

// BenchmarkTCPDial 测试原始 TCP 连接性能（本地回环）
func BenchmarkTCPDial(b *testing.B) {
	// 本地监听一个端口
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		b.Skip("无法创建监听器")
	}
	defer listener.Close()

	addr := listener.Addr().String()

	// 后台接受连接
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
		if err == nil {
			conn.Close()
		}
	}
}

// BenchmarkResultCollectorAdd 测试结果收集器添加性能
func BenchmarkResultCollectorAdd(b *testing.B) {
	collector := &resultCollector{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector.Add("192.168.1.1:80")
	}
}

// BenchmarkResultCollectorAddParallel 测试结果收集器并发添加性能
func BenchmarkResultCollectorAddParallel(b *testing.B) {
	collector := &resultCollector{}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			collector.Add("192.168.1.1:80")
		}
	})
}

// BenchmarkResultCollectorGetAll 测试结果收集器获取全部性能
func BenchmarkResultCollectorGetAll(b *testing.B) {
	collector := &resultCollector{}
	// 预填充数据
	for i := 0; i < 1000; i++ {
		collector.Add("192.168.1.1:80")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = collector.GetAll()
	}
}

// BenchmarkFailedPortCollectorAdd 测试失败端口收集器添加性能
func BenchmarkFailedPortCollectorAdd(b *testing.B) {
	collector := &failedPortCollector{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector.Add("192.168.1.1", 80, "192.168.1.1:80")
	}
}

