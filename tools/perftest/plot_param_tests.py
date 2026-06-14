#!/usr/bin/env python3
"""
fscan 参数测试结果可视化
"""

import csv
import matplotlib.pyplot as plt
import numpy as np
import os

# 设置中文字体
plt.rcParams['font.sans-serif'] = ['Microsoft YaHei', 'SimHei', 'DejaVu Sans']
plt.rcParams['axes.unicode_minus'] = False

def read_csv(filepath):
    results = []
    with open(filepath, 'r', encoding='utf-8') as f:
        reader = csv.DictReader(f)
        for row in reader:
            results.append(row)
    return results

def create_charts(output_dir):
    # 读取数据
    time_file = os.path.join(output_dir, "time_results.csv")
    mt_file = os.path.join(output_dir, "mt_results.csv")

    time_data = read_csv(time_file)
    mt_data = read_csv(mt_file)

    # 创建图表
    fig, axes = plt.subplots(1, 2, figsize=(14, 5))
    fig.suptitle('fscan 参数性能测试 (目标: 1.1.1.0/24)', fontsize=14, fontweight='bold')

    # ==================== 图1: -time 超时参数 ====================
    ax1 = axes[0]

    times = [int(d['time_seconds']) for d in time_data]
    pps = [float(d['packets_per_sec']) for d in time_data]
    duration = [float(d['duration_ms'])/1000 for d in time_data]

    # 双Y轴
    color1 = '#2ecc71'
    ax1.set_xlabel('超时时间 -time (秒)', fontsize=11)
    ax1.set_ylabel('吞吐量 (pps)', color=color1, fontsize=11)
    bars = ax1.bar([x - 0.2 for x in range(len(times))], pps, 0.4, color=color1, alpha=0.7, label='吞吐量')
    ax1.tick_params(axis='y', labelcolor=color1)
    ax1.set_xticks(range(len(times)))
    ax1.set_xticklabels([f'{t}s' for t in times])

    # 标注最优值
    max_idx = np.argmax(pps)
    ax1.bar(max_idx - 0.2, pps[max_idx], 0.4, color='#27ae60', alpha=0.9, edgecolor='red', linewidth=2)

    ax1_twin = ax1.twinx()
    color2 = '#e74c3c'
    ax1_twin.set_ylabel('扫描耗时 (秒)', color=color2, fontsize=11)
    ax1_twin.bar([x + 0.2 for x in range(len(times))], duration, 0.4, color=color2, alpha=0.7, label='耗时')
    ax1_twin.tick_params(axis='y', labelcolor=color2)

    # 添加数值标签
    for i, (p, d) in enumerate(zip(pps, duration)):
        ax1.text(i - 0.2, p + 2, f'{p:.1f}', ha='center', fontsize=9, color=color1)
        ax1_twin.text(i + 0.2, d + 0.3, f'{d:.1f}s', ha='center', fontsize=9, color=color2)

    ax1.set_title('-time 超时参数影响\n(默认值: 3秒)', fontsize=12)
    ax1.axhline(y=pps[2], color='gray', linestyle='--', alpha=0.5, label='默认值基准')

    # 计算相对于默认值的提升
    default_pps = pps[2]  # time=3 是默认值
    improvement = [(p - default_pps) / default_pps * 100 for p in pps]

    # ==================== 图2: -mt 模块线程参数 ====================
    ax2 = axes[1]

    mts = [int(d['module_threads']) for d in mt_data]
    mt_pps = [float(d['packets_per_sec']) for d in mt_data]
    mt_duration = [float(d['duration_ms'])/1000 for d in mt_data]

    ax2.bar(range(len(mts)), mt_pps, color='#3498db', alpha=0.7)
    ax2.set_xlabel('模块线程数 -mt', fontsize=11)
    ax2.set_ylabel('吞吐量 (pps)', fontsize=11)
    ax2.set_xticks(range(len(mts)))
    ax2.set_xticklabels(mts)
    ax2.set_title('-mt 模块线程参数影响\n(默认值: 20, 测试时禁用POC)', fontsize=12)

    # 添加数值标签
    for i, p in enumerate(mt_pps):
        ax2.text(i, p + 0.5, f'{p:.1f}', ha='center', fontsize=9)

    # 计算变化范围
    mt_range = max(mt_pps) - min(mt_pps)
    ax2.set_ylim(min(mt_pps) - 5, max(mt_pps) + 5)

    # 添加注释
    ax2.text(0.5, 0.95, f'变化幅度: {mt_range:.2f} pps (可忽略)',
             transform=ax2.transAxes, ha='center', fontsize=10,
             bbox=dict(boxstyle='round', facecolor='wheat', alpha=0.5))

    plt.tight_layout()

    output_path = os.path.join(output_dir, 'param_test_chart.png')
    plt.savefig(output_path, dpi=150, bbox_inches='tight')
    print(f"图表已保存: {output_path}")
    plt.close()

    # ==================== 生成详细分析图 ====================
    fig2, ax = plt.subplots(figsize=(10, 6))

    x = np.arange(len(times))
    width = 0.35

    # 性能提升百分比
    colors = ['#e74c3c' if imp < 0 else '#2ecc71' for imp in improvement]
    bars = ax.bar(x, improvement, width, color=colors, alpha=0.8)

    ax.axhline(y=0, color='black', linestyle='-', linewidth=0.5)
    ax.set_xlabel('超时时间 -time (秒)', fontsize=12)
    ax.set_ylabel('相对默认值(3秒)的性能变化 (%)', fontsize=12)
    ax.set_title('-time 参数优化效果分析', fontsize=14, fontweight='bold')
    ax.set_xticks(x)
    ax.set_xticklabels([f'{t}s' for t in times])

    # 添加数值标签
    for i, (bar, imp) in enumerate(zip(bars, improvement)):
        height = bar.get_height()
        ax.text(bar.get_x() + bar.get_width()/2., height + (1 if height >= 0 else -3),
                f'{imp:+.1f}%', ha='center', va='bottom' if height >= 0 else 'top',
                fontsize=11, fontweight='bold')

    # 添加建议
    ax.text(0.02, 0.98,
            '建议:\n• 内网环境: -time 1 或 2\n• 公网环境: -time 3 (默认)\n• 高延迟网络: -time 5+',
            transform=ax.transAxes, fontsize=10, verticalalignment='top',
            bbox=dict(boxstyle='round', facecolor='lightyellow', alpha=0.8))

    plt.tight_layout()

    output_path2 = os.path.join(output_dir, 'time_optimization_chart.png')
    plt.savefig(output_path2, dpi=150, bbox_inches='tight')
    print(f"图表已保存: {output_path2}")
    plt.close()

def main():
    output_dir = "results/param_tests"
    create_charts(output_dir)

if __name__ == "__main__":
    main()
