/*
 * scanner.c - 核心TCP端口扫描实现
 * 
 * 极简但可靠的端口扫描逻辑，专注于内网环境
 */

#include "../include/platform.h"

/*
 * 基础TCP连接测试
 * 返回: 1=端口开放, 0=端口关闭, -1=错误
 */
int tcp_connect_test(const char* host, int port, int timeout) {
    socket_t sock;
    struct sockaddr_in addr;
    int result;
    
    /* 参数验证 */
    if (!host || port <= 0 || port > 65535) {
        return -1;
    }
    
    /* 创建socket */
    sock = socket(AF_INET, SOCK_STREAM, 0);
    if (sock == INVALID_SOCKET_VALUE) {
        return -1;
    }
    
    /* 设置超时 */
    if (set_socket_timeout(sock, timeout) != 0) {
        close_socket(sock);
        return -1;
    }
    
    /* 设置目标地址 */
    memset(&addr, 0, sizeof(addr));
    addr.sin_family = AF_INET;
    addr.sin_port = htons((unsigned short)port);
    
    /* 转换IP地址 */
    addr.sin_addr.s_addr = inet_addr(host);
    if (addr.sin_addr.s_addr == INADDR_NONE) {
        /* 如果不是有效IP，当作域名处理 */
        struct hostent* he;
        he = gethostbyname(host);
        if (!he) {
            close_socket(sock);
            return -1;
        }
        memcpy(&addr.sin_addr, he->h_addr_list[0], he->h_length);
    }
    
    /* 执行连接测试 */
    result = connect(sock, (struct sockaddr*)&addr, sizeof(addr));
    
    /* 关闭socket */
    close_socket(sock);
    
    /* 返回结果 */
    return (result == 0) ? 1 : 0;
}

/*
 * 扫描单个主机的多个端口
 */
int scan_host_ports(const char* host, const int* ports, int port_count, int timeout) {
    int i;
    int open_count = 0;
    
    if (!host || !ports || port_count <= 0) {
        return 0;
    }
    
    printf("Scanning %s...\n", host);
    
    for (i = 0; i < port_count; i++) {
        int result = tcp_connect_test(host, ports[i], timeout);
        
        if (result == 1) {
            printf("%s:%d open\n", host, ports[i]);
            open_count++;
        } else if (result == -1) {
            /* 静默处理错误，继续扫描 */
        }
        
        /* 简单的进度指示 */
        if ((i + 1) % 100 == 0 || i == port_count - 1) {
            printf("Progress: %d/%d ports scanned\n", i + 1, port_count);
        }
    }
    
    return open_count;
}

/*
 * 解析端口范围字符串
 * 支持: "80", "80,443", "1-1000", "80,443,8000-8080"
 */
int parse_ports(const char* port_str, int* ports, int max_ports) {
    char* str_copy;
    char* token;
    int count = 0;
    
    if (!port_str || !ports || max_ports <= 0) {
        return 0;
    }
    
    /* 复制字符串以便修改 */
    str_copy = malloc(strlen(port_str) + 1);
    if (!str_copy) {
        return 0;
    }
    strcpy(str_copy, port_str);
    
    /* 使用strtok（兼容性更好） */
    token = strtok(str_copy, ",");
    
    while (token && count < max_ports) {
        char* dash = strchr(token, '-');
        
        if (dash) {
            /* 处理范围 "start-end" */
            int start, end, i;
            *dash = '\0';
            start = atoi(token);
            end = atoi(dash + 1);
            
            if (start > 0 && end > 0 && start <= end && end <= 65535) {
                for (i = start; i <= end && count < max_ports; i++) {
                    ports[count++] = i;
                }
            }
        } else {
            /* 处理单个端口 */
            int port = atoi(token);
            if (port > 0 && port <= 65535) {
                ports[count++] = port;
            }
        }
        
        token = strtok(NULL, ",");
    }
    
    free(str_copy);
    return count;
}

/*
 * 简单的IP范围解析
 * 目前只支持单个IP，后续可扩展
 */
int parse_hosts(const char* host_str, char hosts[][MAX_HOST_LEN], int max_hosts) {
    if (!host_str || !hosts || max_hosts <= 0) {
        return 0;
    }
    
    /* 目前简化实现：只处理单个主机 */
    if (strlen(host_str) < MAX_HOST_LEN) {
        strcpy(hosts[0], host_str);
        return 1;
    }
    
    return 0;
}