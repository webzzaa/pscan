/*
 * main.c - fscan-lite 主程序
 * 
 * 极简的TCP内网端口扫描器
 * 目标：最大兼容性，最小复杂度
 */

#include "../include/platform.h"

/* 函数声明 */
int tcp_connect_test(const char* host, int port, int timeout);
int scan_host_ports(const char* host, const int* ports, int port_count, int timeout);
int parse_ports(const char* port_str, int* ports, int max_ports);
int parse_hosts(const char* host_str, char hosts[][MAX_HOST_LEN], int max_hosts);

/* 显示版本信息 */
void show_version(void) {
    printf("fscan-lite v1.0 - Lightweight TCP Port Scanner\n");
    printf("Built for maximum compatibility (Windows 98 - Windows 11, Linux glibc 2.3+)\n");
    printf("Copyright (c) 2024\n");
}

/* 显示使用帮助 */
void show_usage(const char* program_name) {
    printf("Usage: %s [OPTIONS]\n", program_name);
    printf("\n");
    printf("Options:\n");
    printf("  -h HOST     Target host (IP address)\n");
    printf("  -p PORTS    Ports to scan (e.g. 80,443 or 1-1000)\n");
    printf("  -t TIMEOUT  Connection timeout in seconds (default: 3)\n");
    printf("  --help      Show this help message\n");
    printf("  --version   Show version information\n");
    printf("\n");
    printf("Examples:\n");
    printf("  %s -h 192.168.1.1 -p 22,80,443\n", program_name);
    printf("  %s -h 10.0.0.1 -p 1-1000 -t 2\n", program_name);
    printf("\n");
}

/* 主函数 */
int main(int argc, char* argv[]) {
    char* target_host = NULL;
    char* port_string = NULL;
    int timeout = DEFAULT_TIMEOUT;
    int ports[1000];  /* 支持最多1000个端口 */
    int port_count = 0;
    int i;
    int result;
    
    /* 参数解析 */
    for (i = 1; i < argc; i++) {
        if (strcmp(argv[i], "-h") == 0 && i + 1 < argc) {
            target_host = argv[++i];
        }
        else if (strcmp(argv[i], "-p") == 0 && i + 1 < argc) {
            port_string = argv[++i];
        }
        else if (strcmp(argv[i], "-t") == 0 && i + 1 < argc) {
            timeout = atoi(argv[++i]);
            if (timeout <= 0) timeout = DEFAULT_TIMEOUT;
        }
        else if (strcmp(argv[i], "--help") == 0) {
            show_usage(argv[0]);
            return 0;
        }
        else if (strcmp(argv[i], "--version") == 0) {
            show_version();
            return 0;
        }
        else {
            printf("Unknown option: %s\n", argv[i]);
            show_usage(argv[0]);
            return 1;
        }
    }
    
    /* 验证必需参数 */
    if (!target_host) {
        printf("Error: Target host (-h) is required\n");
        show_usage(argv[0]);
        return 1;
    }
    
    if (!port_string) {
        printf("Error: Ports (-p) are required\n");
        show_usage(argv[0]);
        return 1;
    }
    
    /* 初始化平台 */
    if (platform_init() != 0) {
        printf("Error: Failed to initialize platform\n");
        return 1;
    }
    
    /* 解析端口 */
    port_count = parse_ports(port_string, ports, sizeof(ports) / sizeof(ports[0]));
    if (port_count == 0) {
        printf("Error: No valid ports specified\n");
        platform_cleanup();
        return 1;
    }
    
    printf("fscan-lite - Starting scan\n");
    printf("Target: %s\n", target_host);
    printf("Ports: %d ports to scan\n", port_count);
    printf("Timeout: %d seconds\n", timeout);
    printf("=================================\n");
    
    /* 执行扫描 */
    result = scan_host_ports(target_host, ports, port_count, timeout);
    
    printf("=================================\n");
    printf("Scan completed: %d open ports found\n", result);
    
    /* 清理资源 */
    platform_cleanup();
    
    return 0;
}