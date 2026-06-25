@echo off
setlocal enabledelayedexpansion

:: fscan 可扩展性测试脚本
:: 用法: run_perftest.bat 192.168.1.0/24

set TARGET=%1
if "%TARGET%"=="" (
    echo 用法: run_perftest.bat ^<target^>
    echo 示例: run_perftest.bat 192.168.1.0/24
    exit /b 1
)

set PORTS=22,80,443,3389,8080
set THREADS=100 200 400 600 800 1000 1500 2000
set OUTPUT=perf_results.csv

echo threads,duration_sec,timestamp > %OUTPUT%

echo === fscan 可扩展性测试 ===
echo 目标: %TARGET%
echo 端口: %PORTS%

for %%t in (%THREADS%) do (
    echo.
    echo [测试] 线程数=%%t

    :: 记录开始时间
    set START=%time%

    :: 运行 fscan
    fscan.exe -h %TARGET% -p %PORTS% -t %%t -np -nopoc -o NUL 2>NUL

    :: 记录结束时间并计算耗时
    set END=%time%

    :: 简单输出（实际耗时需要手动计算或用 PowerShell）
    echo %%t,%START%-%END%,%date% >> %OUTPUT%
    echo 完成: %%t 线程
)

echo.
echo 结果已保存到: %OUTPUT%
echo 使用 plot_results.py 绘图
