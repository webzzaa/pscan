#!/usr/bin/env python3
"""
fscan 可扩展性图表生成工具

用法:
    python plot_results.py perf_results.csv
    python plot_results.py perf_results.csv -o my_chart.png
"""

import sys
import argparse

def main():
    parser = argparse.ArgumentParser(description='绘制 fscan 可扩展性图表')
    parser.add_argument('csv_file', help='CSV 数据文件')
    parser.add_argument('-o', '--output', default='scalability.png', help='输出图片文件')
    parser.add_argument('--style', choices=['default', 'dark', 'minimal'], default='default')
    args = parser.parse_args()

    try:
        import pandas as pd
        import matplotlib.pyplot as plt
        import matplotlib.ticker as ticker
        import matplotlib as mpl
    except ImportError:
        print("需要安装依赖: pip install pandas matplotlib")
        sys.exit(1)

    # 设置中文字体
    plt.rcParams['font.sans-serif'] = ['Microsoft YaHei', 'SimHei', 'DejaVu Sans']
    plt.rcParams['axes.unicode_minus'] = False  # 解决负号显示问题

    # 读取数据
    df = pd.read_csv(args.csv_file)

    # 创建图表
    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(14, 5))

    # 样式设置
    if args.style == 'dark':
        plt.style.use('dark_background')
        color1, color2 = '#00ff88', '#ff6b6b'
    else:
        color1, color2 = '#2563eb', '#dc2626'

    # 图1: 吞吐量 vs 线程数
    ax1.plot(df['threads'], df['ports_per_sec'], 'o-',
             color=color1, linewidth=2.5, markersize=10, label='实测吞吐量')

    # 理想线性扩展线（以第一个点为基准）
    if len(df) > 0:
        base_rate = df['ports_per_sec'].iloc[0]
        base_threads = df['threads'].iloc[0]
        ideal = [base_rate * (t / base_threads) for t in df['threads']]
        ax1.plot(df['threads'], ideal, '--', color='gray', alpha=0.5, label='理想线性扩展')

    ax1.set_xlabel('线程数', fontsize=12)
    ax1.set_ylabel('扫描速率 (端口/秒)', fontsize=12)
    ax1.set_title('fscan 可扩展性曲线', fontsize=14, fontweight='bold')
    ax1.grid(True, alpha=0.3)
    ax1.legend(loc='upper left')

    # 标注峰值点
    max_idx = df['ports_per_sec'].idxmax()
    max_threads = df['threads'].iloc[max_idx]
    max_rate = df['ports_per_sec'].iloc[max_idx]
    ax1.annotate(f'峰值: {max_rate:.0f} 端口/秒\n@ {max_threads} 线程',
                 xy=(max_threads, max_rate),
                 xytext=(max_threads - 200, max_rate * 0.75),
                 fontsize=10,
                 arrowprops=dict(arrowstyle='->', color='gray'))

    # 图2: 扫描耗时 vs 线程数
    ax2.plot(df['threads'], df['duration_sec'], 's-',
             color=color2, linewidth=2.5, markersize=10)

    ax2.set_xlabel('线程数', fontsize=12)
    ax2.set_ylabel('扫描耗时 (秒)', fontsize=12)
    ax2.set_title('扫描耗时曲线', fontsize=14, fontweight='bold')
    ax2.grid(True, alpha=0.3)

    # 标注最快点
    min_idx = df['duration_sec'].idxmin()
    min_threads = df['threads'].iloc[min_idx]
    min_duration = df['duration_sec'].iloc[min_idx]
    ax2.annotate(f'最快: {min_duration:.2f} 秒\n@ {min_threads} 线程',
                 xy=(min_threads, min_duration),
                 xytext=(min_threads - 200, min_duration * 1.5),
                 fontsize=10,
                 arrowprops=dict(arrowstyle='->', color='gray'))

    plt.tight_layout()
    plt.savefig(args.output, dpi=150, bbox_inches='tight')
    print(f"图表已保存: {args.output}")

    # 打印分析结论
    print("\n=== 分析结论 ===")
    print(f"最优线程数: {max_threads} (峰值吞吐量: {max_rate:.0f} ports/sec)")
    print(f"最短耗时: {min_duration:.2f}s @ {min_threads} 线程")

    # 计算扩展效率
    if len(df) >= 2:
        efficiency = (df['ports_per_sec'].iloc[-1] / df['ports_per_sec'].iloc[0]) / \
                     (df['threads'].iloc[-1] / df['threads'].iloc[0]) * 100
        print(f"扩展效率: {efficiency:.1f}% (相对于线性扩展)")

        if efficiency < 50:
            print("⚠️  扩展效率较低，可能存在锁竞争或资源瓶颈")
        elif efficiency > 80:
            print("✅ 扩展效率良好")


if __name__ == '__main__':
    main()
