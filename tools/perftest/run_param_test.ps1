# fscan 参数性能测试脚本
# 测试 -time 和 -mt 参数对性能的影响

param(
    [Parameter(Mandatory=$true)]
    [string]$Target,

    [string]$Ports = "22,80,443,3389,8080",
    [int]$Threads = 600,
    [string]$OutputDir = "results/param_tests"
)

$fscanPath = Join-Path $PSScriptRoot "..\..\fscan.exe"
if (-not (Test-Path $fscanPath)) {
    $fscanPath = "fscan.exe"
}

# 创建输出目录
$fullOutputDir = Join-Path $PSScriptRoot $OutputDir
New-Item -ItemType Directory -Force -Path $fullOutputDir | Out-Null

Write-Host "=== fscan 参数性能测试 ===" -ForegroundColor Cyan
Write-Host "目标: $Target"
Write-Host "端口: $Ports"
Write-Host "基准线程: $Threads"
Write-Host ""

# ==================== 测试1: -time 超时参数 ====================
Write-Host ">>> 测试1: -time 超时参数 <<<" -ForegroundColor Yellow
$timeValues = @(1, 2, 3, 5)
$timeResults = @()

foreach ($t in $timeValues) {
    Write-Host "[time=$t] " -NoNewline

    $output = & $fscanPath -h $Target -p $Ports -t $Threads -time $t -np -nopoc -perf -no 2>&1 | Out-String

    if ($output -match '\[PERF_STATS_JSON\](.*?)\[/PERF_STATS_JSON\]') {
        $stats = $Matches[1] | ConvertFrom-Json
        $timeResults += [PSCustomObject]@{
            time_seconds = $t
            duration_ms = $stats.scan_duration_ms
            packets_per_sec = [math]::Round($stats.packets_per_second, 2)
            tcp_success = $stats.tcp_success
            tcp_failed = $stats.tcp_failed
            success_rate = [math]::Round($stats.success_rate, 2)
        }
        Write-Host "耗时: $([math]::Round($stats.scan_duration_ms/1000, 2))s, $([math]::Round($stats.packets_per_second, 1)) pps, 成功率: $([math]::Round($stats.success_rate, 1))%" -ForegroundColor Green
    } else {
        Write-Host "解析失败" -ForegroundColor Red
    }
}

$timeResults | Export-Csv -Path (Join-Path $fullOutputDir "time_results.csv") -NoTypeInformation
Write-Host ""

# ==================== 测试2: -mt 模块线程参数 ====================
Write-Host ">>> 测试2: -mt 模块线程参数 <<<" -ForegroundColor Yellow
$mtValues = @(5, 10, 20, 50, 100)
$mtResults = @()

foreach ($mt in $mtValues) {
    Write-Host "[mt=$mt] " -NoNewline

    # 注意: -mt 主要影响服务识别和POC，这里不用 -nopoc 来观察效果
    $output = & $fscanPath -h $Target -p $Ports -t $Threads -mt $mt -np -nopoc -perf -no 2>&1 | Out-String

    if ($output -match '\[PERF_STATS_JSON\](.*?)\[/PERF_STATS_JSON\]') {
        $stats = $Matches[1] | ConvertFrom-Json
        $mtResults += [PSCustomObject]@{
            module_threads = $mt
            duration_ms = $stats.scan_duration_ms
            packets_per_sec = [math]::Round($stats.packets_per_second, 2)
            tcp_success = $stats.tcp_success
            tcp_failed = $stats.tcp_failed
            success_rate = [math]::Round($stats.success_rate, 2)
        }
        Write-Host "耗时: $([math]::Round($stats.scan_duration_ms/1000, 2))s, $([math]::Round($stats.packets_per_second, 1)) pps" -ForegroundColor Green
    } else {
        Write-Host "解析失败" -ForegroundColor Red
    }
}

$mtResults | Export-Csv -Path (Join-Path $fullOutputDir "mt_results.csv") -NoTypeInformation

Write-Host ""
Write-Host "=== 测试完成 ===" -ForegroundColor Cyan
Write-Host "结果保存到: $fullOutputDir" -ForegroundColor Yellow

# 打印汇总
Write-Host ""
Write-Host ">>> -time 测试结果 <<<" -ForegroundColor Cyan
$timeResults | Format-Table -AutoSize

Write-Host ">>> -mt 测试结果 <<<" -ForegroundColor Cyan
$mtResults | Format-Table -AutoSize
