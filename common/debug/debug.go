//go:build debug
// +build debug

package debug

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
)

var (
	cpuProfile   *os.File
	traceFile    *os.File
	profilesPath = "./profiles"
)

func Start() {
	if err := os.MkdirAll(profilesPath, 0755); err != nil {
		fmt.Printf("[DEBUG] 创建 profiles 目录失败: %v\n", err)
		return
	}

	var err error
	cpuProfile, err = os.Create(profilesPath + "/cpu.prof")
	if err != nil {
		fmt.Printf("[DEBUG] 创建 CPU profile 失败: %v\n", err)
	} else {
		if err := pprof.StartCPUProfile(cpuProfile); err != nil {
			fmt.Printf("[DEBUG] 启动 CPU profile 失败: %v\n", err)
			cpuProfile.Close()
			cpuProfile = nil
		} else {
			fmt.Printf("[DEBUG] CPU profiling 已启动 -> %s/cpu.prof\n", profilesPath)
		}
	}

	traceFile, err = os.Create(profilesPath + "/trace.out")
	if err != nil {
		fmt.Printf("[DEBUG] 创建 trace 文件失败: %v\n", err)
	} else {
		if err := trace.Start(traceFile); err != nil {
			fmt.Printf("[DEBUG] 启动 trace 失败: %v\n", err)
			traceFile.Close()
			traceFile = nil
		} else {
			fmt.Printf("[DEBUG] Execution trace 已启动 -> %s/trace.out\n", profilesPath)
		}
	}

	fmt.Printf("[DEBUG] 性能分析已启动，程序结束时自动保存到 %s/\n", profilesPath)
}

func Stop() {
	if cpuProfile != nil {
		pprof.StopCPUProfile()
		cpuProfile.Close()
		fmt.Printf("[DEBUG] CPU profile 已保存\n")
	}

	if traceFile != nil {
		trace.Stop()
		traceFile.Close()
		fmt.Printf("[DEBUG] Trace 已保存\n")
	}

	memProfile, err := os.Create(profilesPath + "/mem.prof")
	if err != nil {
		fmt.Printf("[DEBUG] 创建内存 profile 失败: %v\n", err)
	} else {
		runtime.GC()
		if err := pprof.WriteHeapProfile(memProfile); err != nil {
			fmt.Printf("[DEBUG] 写入内存 profile 失败: %v\n", err)
		} else {
			fmt.Printf("[DEBUG] 内存 profile 已保存 -> %s/mem.prof\n", profilesPath)
		}
		memProfile.Close()
	}

	goroutineProfile, err := os.Create(profilesPath + "/goroutine.prof")
	if err != nil {
		fmt.Printf("[DEBUG] 创建 goroutine profile 失败: %v\n", err)
	} else {
		if err := pprof.Lookup("goroutine").WriteTo(goroutineProfile, 0); err != nil {
			fmt.Printf("[DEBUG] 写入 goroutine profile 失败: %v\n", err)
		} else {
			fmt.Printf("[DEBUG] Goroutine profile 已保存 -> %s/goroutine.prof\n", profilesPath)
		}
		goroutineProfile.Close()
	}

	fmt.Printf("\n[DEBUG] 所有性能分析文件已保存到 %s/\n", profilesPath)
	fmt.Printf("[DEBUG] 查看方法:\n")
	fmt.Printf("  CPU 火焰图:      go tool pprof -http=:8081 %s/cpu.prof\n", profilesPath)
	fmt.Printf("  内存火焰图:      go tool pprof -http=:8081 %s/mem.prof\n", profilesPath)
	fmt.Printf("  协程分析:        go tool pprof -http=:8081 %s/goroutine.prof\n", profilesPath)
	fmt.Printf("  执行时间线:      go tool trace %s/trace.out\n", profilesPath)
}
