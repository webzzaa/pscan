# fscan 可扩展性测试 PowerShell 脚本
# 用法: .\run_perftest.ps1 -Target 192.168.1.0/24

param(
    [Parameter(Mandatory=$true)]
    [string]$Target,

    [string]$Ports = "22,80,443,3389,8080",
    [int[]]$Threads = @(100, 200, 400, 600, 800, 1000, 1500),
    [int]$Repeat = 3,
    [string]$Output = "perf_results.csv"
)

$fscanPath = Join-Path $PSScriptRoot "..\..\fscan.exe"
if (-not (Test-Path $fscanPath)) {
    $fscanPath = "fscan.exe"
}

Write-Host "=== fscan 可扩展性测试 ===" -ForegroundColor Cyan
Write-Host "目标: $Target"
Write-Host "端口: $Ports"
Write-Host "线程数: $($Threads -join ', ')"
Write-Host "重复次数: $Repeat"
Write-Host ""

$results = @()

foreach ($t in $Threads) {
    Write-Host "[线程=$t] " -NoNewline

    $durations = @()

    for ($i = 1; $i -le $Repeat; $i++) {
        Write-Host "." -NoNewline

        $stopwatch = [System.Diagnostics.Stopwatch]::StartNew()

        $proc = Start-Process -FilePath $fscanPath -ArgumentList @(
            "-h", $Target,
            "-p", $Ports,
            "-t", $t,
            "-np",
            "-nopoc",
            "-o", "NUL"
        ) -NoNewWindow -Wait -PassThru

        $stopwatch.Stop()
        $durations += $stopwatch.Elapsed.TotalSeconds
    }

    $avgDuration = ($durations | Measure-Object -Average).Average

    # 估算端口扫描数
    $portCount = ($Ports -split ',').Count
    if ($Target -match '/24') { $ipCount = 254 }
    elseif ($Target -match '/16') { $ipCount = 65534 }
    else { $ipCount = 1 }

    $totalPorts = $ipCount * $portCount
    $portsPerSec = if ($avgDuration -gt 0) { $totalPorts / $avgDuration } else { 0 }

    $results += [PSCustomObject]@{
        threads = $t
        duration_sec = [math]::Round($avgDuration, 3)
        ports_per_sec = [math]::Round($portsPerSec, 1)
        total_ports = $totalPorts
    }

    Write-Host " 平均: $([math]::Round($avgDuration, 2))s, $([math]::Round($portsPerSec, 0)) ports/sec" -ForegroundColor Green
}

# 导出 CSV
$results | Export-Csv -Path $Output -NoTypeInformation
Write-Host "`n结果已保存到: $Output" -ForegroundColor Yellow

# 打印绘图命令
Write-Host "`n=== 绘图命令 ===" -ForegroundColor Cyan
Write-Host "python plot_results.py $Output"
