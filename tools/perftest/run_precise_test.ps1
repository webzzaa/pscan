# fscan 精确性能测试脚本 (使用内部指标)
# 用法: .\run_precise_test.ps1 -Target 1.1.1.0/24

param(
    [Parameter(Mandatory=$true)]
    [string]$Target,

    [string]$Ports = "22,80,443,3389,8080",
    [int[]]$Threads = @(200, 400, 600, 800, 1000, 1500, 2000),
    [int]$Repeat = 1,
    [string]$Output = "precise_results.csv"
)

$fscanPath = Join-Path $PSScriptRoot "..\..\fscan.exe"
if (-not (Test-Path $fscanPath)) {
    $fscanPath = "fscan.exe"
}

Write-Host "=== fscan 精确性能测试 (内部指标) ===" -ForegroundColor Cyan
Write-Host "目标: $Target"
Write-Host "端口: $Ports"
Write-Host "线程数: $($Threads -join ', ')"
Write-Host ""

$results = @()

foreach ($t in $Threads) {
    Write-Host "[线程=$t] " -NoNewline

    $allStats = @()

    for ($i = 1; $i -le $Repeat; $i++) {
        Write-Host "." -NoNewline

        # 运行 fscan 并捕获输出
        $output = & $fscanPath -h $Target -p $Ports -t $t -np -nopoc -perf -no 2>&1 | Out-String

        # 提取 JSON
        if ($output -match '\[PERF_STATS_JSON\](.*?)\[/PERF_STATS_JSON\]') {
            $jsonStr = $Matches[1]
            $stats = $jsonStr | ConvertFrom-Json
            $allStats += $stats
        }
    }

    if ($allStats.Count -gt 0) {
        # 计算平均值
        $avgDuration = ($allStats | Measure-Object -Property scan_duration_ms -Average).Average
        $avgPPS = ($allStats | Measure-Object -Property packets_per_second -Average).Average
        $avgTotal = ($allStats | Measure-Object -Property total_packets -Average).Average
        $avgSuccess = ($allStats | Measure-Object -Property tcp_success -Average).Average
        $avgFailed = ($allStats | Measure-Object -Property tcp_failed -Average).Average
        $avgSuccessRate = ($allStats | Measure-Object -Property success_rate -Average).Average

        $results += [PSCustomObject]@{
            threads = $t
            duration_ms = [math]::Round($avgDuration, 0)
            packets_per_sec = [math]::Round($avgPPS, 2)
            total_packets = [math]::Round($avgTotal, 0)
            tcp_success = [math]::Round($avgSuccess, 0)
            tcp_failed = [math]::Round($avgFailed, 0)
            success_rate = [math]::Round($avgSuccessRate, 2)
        }

        Write-Host " 耗时: $([math]::Round($avgDuration/1000, 2))s, $([math]::Round($avgPPS, 1)) pkt/s, 成功率: $([math]::Round($avgSuccessRate, 1))%" -ForegroundColor Green
    } else {
        Write-Host " 解析失败" -ForegroundColor Red
    }
}

# 导出 CSV
$results | Export-Csv -Path $Output -NoTypeInformation
Write-Host "`n结果已保存到: $Output" -ForegroundColor Yellow

# 打印数据表格
Write-Host "`n=== 测试结果 ===" -ForegroundColor Cyan
$results | Format-Table -AutoSize
