#!/usr/bin/env python3
"""
Go Benchmark 结果可视化脚本
生成柱状图展示各函数的性能指标
"""

import re
import matplotlib.pyplot as plt
import numpy as np
import sys
import os

# 设置中文字体
plt.rcParams['font.sans-serif'] = ['Microsoft YaHei', 'SimHei', 'DejaVu Sans']
plt.rcParams['axes.unicode_minus'] = False

def parse_benchmark_results(filepath):
    """解析 Go benchmark 输出"""
    results = []

    with open(filepath, 'r', encoding='utf-8') as f:
        for line in f:
            # 匹配格式: BenchmarkXxx-24    123456    1.234 ns/op    123 B/op    12 allocs/op
            match = re.match(
                r'(Benchmark\w+)-\d+\s+(\d+)\s+([\d.]+)\s+(ns|µs|ms)/op(?:\s+([\d.]+)\s+B/op)?(?:\s+(\d+)\s+allocs/op)?',
                line.strip()
            )
            if match:
                name = match.group(1).replace('Benchmark', '')
                ops = int(match.group(2))
                time_val = float(match.group(3))
                time_unit = match.group(4)
                bytes_op = float(match.group(5)) if match.group(5) else 0
                allocs_op = int(match.group(6)) if match.group(6) else 0

                # 统一转换为 ns
                if time_unit == 'µs':
                    time_ns = time_val * 1000
                elif time_unit == 'ms':
                    time_ns = time_val * 1000000
                else:
                    time_ns = time_val

                results.append({
                    'name': name,
                    'ops': ops,
                    'time_ns': time_ns,
                    'time_val': time_val,
                    'time_unit': time_unit,
                    'bytes': bytes_op,
                    'allocs': allocs_op
                })

    return results

def format_time(ns):
    """格式化时间显示"""
    if ns >= 1000000:
        return f"{ns/1000000:.1f}ms"
    elif ns >= 1000:
        return f"{ns/1000:.1f}µs"
    else:
        return f"{ns:.1f}ns"

def format_bytes(b):
    """格式化内存显示"""
    if b >= 1024*1024:
        return f"{b/1024/1024:.1f}MB"
    elif b >= 1024:
        return f"{b/1024:.1f}KB"
    else:
        return f"{b:.0f}B"

def create_benchmark_charts(results, output_dir):
    """创建 benchmark 可视化图表"""

    if not results:
        print("没有找到 benchmark 结果")
        return

    # 按模块分组
    core_funcs = [r for r in results if any(x in r['name'] for x in
                  ['CheckSum', 'TCPDial', 'ResultCollector', 'FailedPort', 'Estimate', 'Calculate', 'BuildExclude', 'ArrayCount'])]
    parser_funcs = [r for r in results if 'ParseIP' in r['name'] or 'ParsePort' in r['name']]
    finger_funcs = [r for r in results if 'DecodePattern' in r['name']]

    # 图1: 执行时间对比 (对数刻度)
    fig, axes = plt.subplots(2, 2, figsize=(14, 10))
    fig.suptitle('Go Benchmark 性能分析', fontsize=14, fontweight='bold')

    # 子图1: 所有函数执行时间
    ax1 = axes[0, 0]
    names = [r['name'][:20] for r in results]
    times = [r['time_ns'] for r in results]
    colors = ['#2ecc71' if t < 1000 else '#f39c12' if t < 100000 else '#e74c3c' for t in times]

    bars = ax1.barh(names, times, color=colors)
    ax1.set_xscale('log')
    ax1.set_xlabel('执行时间 (ns, 对数刻度)')
    ax1.set_title('各函数执行时间')

    # 添加数值标签
    for bar, t in zip(bars, times):
        ax1.text(t * 1.5, bar.get_y() + bar.get_height()/2,
                format_time(t), va='center', fontsize=8)

    # 子图2: 内存分配
    ax2 = axes[0, 1]
    mem_results = [r for r in results if r['bytes'] > 0]
    if mem_results:
        names = [r['name'][:20] for r in mem_results]
        mem = [r['bytes'] for r in mem_results]
        colors = ['#3498db' if m < 1024 else '#9b59b6' if m < 100000 else '#e74c3c' for m in mem]

        bars = ax2.barh(names, mem, color=colors)
        ax2.set_xscale('log')
        ax2.set_xlabel('内存分配 (B, 对数刻度)')
        ax2.set_title('各函数内存分配')

        for bar, m in zip(bars, mem):
            ax2.text(m * 1.5, bar.get_y() + bar.get_height()/2,
                    format_bytes(m), va='center', fontsize=8)

    # 子图3: 核心模块详细对比
    ax3 = axes[1, 0]
    if core_funcs:
        names = [r['name'][:18] for r in core_funcs]
        times = [r['time_ns'] for r in core_funcs]

        x = np.arange(len(names))
        bars = ax3.bar(x, times, color='#3498db')
        ax3.set_xticks(x)
        ax3.set_xticklabels(names, rotation=45, ha='right', fontsize=8)
        ax3.set_ylabel('执行时间 (ns)')
        ax3.set_title('核心模块 (core) 性能')
        ax3.set_yscale('log')

        for bar, t in zip(bars, times):
            ax3.text(bar.get_x() + bar.get_width()/2, t * 1.2,
                    format_time(t), ha='center', fontsize=7)

    # 子图4: 解析器模块详细对比
    ax4 = axes[1, 1]
    if parser_funcs:
        names = [r['name'].replace('ParseIP', 'IP').replace('ParsePort', 'Port')[:15] for r in parser_funcs]
        times = [r['time_ns'] for r in parser_funcs]

        x = np.arange(len(names))
        bars = ax4.bar(x, times, color='#e74c3c')
        ax4.set_xticks(x)
        ax4.set_xticklabels(names, rotation=45, ha='right', fontsize=8)
        ax4.set_ylabel('执行时间 (ns)')
        ax4.set_title('解析器模块 (parsers) 性能')
        ax4.set_yscale('log')

        for bar, t in zip(bars, times):
            ax4.text(bar.get_x() + bar.get_width()/2, t * 1.2,
                    format_time(t), ha='center', fontsize=7)

    plt.tight_layout()

    output_path = os.path.join(output_dir, 'benchmark_chart.png')
    plt.savefig(output_path, dpi=150, bbox_inches='tight')
    print(f"图表已保存: {output_path}")
    plt.close()

    # 图2: 性能热力图 - 时间 vs 内存
    fig2, ax = plt.subplots(figsize=(10, 6))

    # 筛选有内存分配的结果
    valid_results = [r for r in results if r['bytes'] > 0]
    if valid_results:
        times = [r['time_ns'] for r in valid_results]
        mems = [r['bytes'] for r in valid_results]
        names = [r['name'][:15] for r in valid_results]

        scatter = ax.scatter(times, mems, s=100, c=range(len(valid_results)),
                            cmap='viridis', alpha=0.7, edgecolors='black')

        ax.set_xscale('log')
        ax.set_yscale('log')
        ax.set_xlabel('执行时间 (ns)')
        ax.set_ylabel('内存分配 (B)')
        ax.set_title('性能-内存权衡分析')

        # 添加标签
        for i, (t, m, n) in enumerate(zip(times, mems, names)):
            ax.annotate(n, (t, m), textcoords="offset points",
                       xytext=(5, 5), fontsize=7)

        # 添加参考线
        ax.axhline(y=1024, color='orange', linestyle='--', alpha=0.5, label='1KB')
        ax.axhline(y=1024*1024, color='red', linestyle='--', alpha=0.5, label='1MB')
        ax.axvline(x=1000, color='green', linestyle='--', alpha=0.5, label='1µs')
        ax.axvline(x=1000000, color='purple', linestyle='--', alpha=0.5, label='1ms')
        ax.legend(loc='upper left', fontsize=8)

        output_path2 = os.path.join(output_dir, 'benchmark_tradeoff.png')
        plt.savefig(output_path2, dpi=150, bbox_inches='tight')
        print(f"图表已保存: {output_path2}")
        plt.close()

def main():
    if len(sys.argv) < 2:
        # 默认路径
        input_file = "results/benchmarks/benchmark_results.txt"
        output_dir = "results/benchmarks"
    else:
        input_file = sys.argv[1]
        output_dir = os.path.dirname(input_file) or "."

    if not os.path.exists(input_file):
        print(f"文件不存在: {input_file}")
        sys.exit(1)

    print(f"解析 benchmark 结果: {input_file}")
    results = parse_benchmark_results(input_file)
    print(f"找到 {len(results)} 个 benchmark 结果")

    for r in results:
        print(f"  - {r['name']}: {format_time(r['time_ns'])}, {format_bytes(r['bytes'])}")

    create_benchmark_charts(results, output_dir)

if __name__ == "__main__":
    main()
