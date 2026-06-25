#!/usr/bin/env python3
"""
fscan 内部指标可视化脚本
生成线程数-性能关系图
"""

import csv
import matplotlib.pyplot as plt
import numpy as np
import sys
import os

# 设置中文字体
plt.rcParams['font.sans-serif'] = ['Microsoft YaHei', 'SimHei', 'DejaVu Sans']
plt.rcParams['axes.unicode_minus'] = False

def read_csv(filepath):
    """读取 CSV 结果文件"""
    results = []
    with open(filepath, 'r', encoding='utf-8') as f:
        reader = csv.DictReader(f)
        for row in reader:
            results.append({
                'threads': int(row['threads']),
                'duration_ms': float(row['duration_ms']),
                'pps': float(row['packets_per_sec']),
                'total': int(row['total_packets']),
                'success': int(row['tcp_success']),
                'failed': int(row['tcp_failed']),
                'success_rate': float(row['success_rate'])
            })
    return results

def create_charts(results, output_dir):
    """创建可视化图表"""

    threads = [r['threads'] for r in results]
    pps = [r['pps'] for r in results]
    duration = [r['duration_ms']/1000 for r in results]  # 转换为秒

    # 创建 2x2 子图
    fig, axes = plt.subplots(2, 2, figsize=(12, 10))
    fig.suptitle('fscan 内部指标性能分析 (目标: 1.1.1.0/24)', fontsize=14, fontweight='bold')

    # 子图1: 吞吐量 vs 线程数
    ax1 = axes[0, 0]
    ax1.plot(threads, pps, 'o-', color='#2ecc71', linewidth=2, markersize=8, label='实测吞吐量')

    # 找到最优点
    max_pps_idx = np.argmax(pps)
    ax1.axvline(x=threads[max_pps_idx], color='red', linestyle='--', alpha=0.7, label=f'最优线程数: {threads[max_pps_idx]}')
    ax1.scatter([threads[max_pps_idx]], [pps[max_pps_idx]], color='red', s=150, zorder=5, marker='*')

    ax1.set_xlabel('线程数')
    ax1.set_ylabel('吞吐量 (packets/s)')
    ax1.set_title('线程数 vs 吞吐量')
    ax1.legend()
    ax1.grid(True, alpha=0.3)

    # 子图2: 扫描耗时 vs 线程数
    ax2 = axes[0, 1]
    ax2.plot(threads, duration, 's-', color='#e74c3c', linewidth=2, markersize=8)

    ax2.set_xlabel('线程数')
    ax2.set_ylabel('扫描耗时 (秒)')
    ax2.set_title('线程数 vs 扫描耗时')
    ax2.grid(True, alpha=0.3)

    # 添加耗时标签
    for t, d in zip(threads, duration):
        ax2.annotate(f'{d:.1f}s', (t, d), textcoords="offset points",
                    xytext=(0, 10), ha='center', fontsize=8)

    # 子图3: 效率分析 (吞吐量/线程数)
    ax3 = axes[1, 0]
    efficiency = [p/t*100 for p, t in zip(pps, threads)]  # 每100线程的吞吐量
    ax3.bar(range(len(threads)), efficiency, color='#3498db', alpha=0.7)
    ax3.set_xticks(range(len(threads)))
    ax3.set_xticklabels(threads)
    ax3.set_xlabel('线程数')
    ax3.set_ylabel('效率 (pps/100线程)')
    ax3.set_title('线程效率分析')

    # 添加数值标签
    for i, e in enumerate(efficiency):
        ax3.text(i, e + 0.5, f'{e:.1f}', ha='center', fontsize=9)

    # 子图4: 加速比分析
    ax4 = axes[1, 1]
    base_pps = pps[0]  # 200线程作为基准
    speedup = [p/base_pps for p in pps]
    ideal_speedup = [t/threads[0] for t in threads]  # 理想线性加速

    ax4.plot(threads, speedup, 'o-', color='#2ecc71', linewidth=2, markersize=8, label='实际加速比')
    ax4.plot(threads, ideal_speedup, '--', color='#95a5a6', linewidth=1.5, label='理想线性加速')

    ax4.set_xlabel('线程数')
    ax4.set_ylabel('加速比 (相对于200线程)')
    ax4.set_title('可扩展性分析')
    ax4.legend()
    ax4.grid(True, alpha=0.3)

    plt.tight_layout()

    output_path = os.path.join(output_dir, 'internal_metrics_chart.png')
    plt.savefig(output_path, dpi=150, bbox_inches='tight')
    print(f"图表已保存: {output_path}")
    plt.close()

    # 单独生成一张主要图表
    fig2, ax = plt.subplots(figsize=(10, 6))

    # 双Y轴
    ax.set_xlabel('线程数', fontsize=12)
    ax.set_ylabel('吞吐量 (packets/s)', color='#2ecc71', fontsize=12)
    line1 = ax.plot(threads, pps, 'o-', color='#2ecc71', linewidth=2.5, markersize=10, label='吞吐量')
    ax.tick_params(axis='y', labelcolor='#2ecc71')
    ax.axvline(x=threads[max_pps_idx], color='red', linestyle='--', alpha=0.5)
    ax.scatter([threads[max_pps_idx]], [pps[max_pps_idx]], color='red', s=200, zorder=5, marker='*')

    ax2 = ax.twinx()
    ax2.set_ylabel('扫描耗时 (秒)', color='#e74c3c', fontsize=12)
    line2 = ax2.plot(threads, duration, 's--', color='#e74c3c', linewidth=2, markersize=8, label='耗时')
    ax2.tick_params(axis='y', labelcolor='#e74c3c')

    # 合并图例
    lines = line1 + line2
    labels = [l.get_label() for l in lines]
    ax.legend(lines, labels, loc='center right')

    ax.set_title('fscan 内部指标: 线程数 vs 性能\n(目标: 1.1.1.0/24, 端口: 22,80,443,3389,8080)', fontsize=13)
    ax.grid(True, alpha=0.3)

    # 添加最优点标注
    ax.annotate(f'最优: {threads[max_pps_idx]}线程\n{pps[max_pps_idx]:.1f} pps',
                xy=(threads[max_pps_idx], pps[max_pps_idx]),
                xytext=(threads[max_pps_idx]+200, pps[max_pps_idx]-10),
                arrowprops=dict(arrowstyle='->', color='red'),
                fontsize=10, color='red')

    output_path2 = os.path.join(output_dir, 'scalability_chart.png')
    plt.savefig(output_path2, dpi=150, bbox_inches='tight')
    print(f"图表已保存: {output_path2}")
    plt.close()

def main():
    if len(sys.argv) < 2:
        input_file = "results/internal_metrics/precise_results.csv"
        output_dir = "results/internal_metrics"
    else:
        input_file = sys.argv[1]
        output_dir = os.path.dirname(input_file) or "."

    if not os.path.exists(input_file):
        print(f"文件不存在: {input_file}")
        sys.exit(1)

    print(f"读取数据: {input_file}")
    results = read_csv(input_file)
    print(f"找到 {len(results)} 条记录")

    for r in results:
        print(f"  线程={r['threads']}: {r['pps']:.1f} pps, 耗时={r['duration_ms']/1000:.1f}s")

    create_charts(results, output_dir)

if __name__ == "__main__":
    main()
