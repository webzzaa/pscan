/*
 * platform.c - 平台抽象层实现
 * 
 * 实现最基础但最可靠的平台相关功能
 */

#include "../include/platform.h"

/* 全局初始化标志 */
static int platform_initialized = 0;

/*
 * 平台初始化
 * Windows: 初始化Winsock
 * Unix: 无需特殊初始化
 */
int platform_init(void) {
    if (platform_initialized) {
        return 0;
    }
    
#ifdef PLATFORM_WINDOWS
    WSADATA wsaData;
    int result;
    
    /* 初始化Winsock - 请求版本2.0，兼容Win98 */
    result = WSAStartup(MAKEWORD(2, 0), &wsaData);
    if (result != 0) {
        /* 如果2.0失败，尝试1.1（Win95/NT兼容） */
        result = WSAStartup(MAKEWORD(1, 1), &wsaData);
        if (result != 0) {
            return -1;
        }
    }
#endif
    
    platform_initialized = 1;
    return 0;
}

/*
 * 平台清理
 */
void platform_cleanup(void) {
    if (!platform_initialized) {
        return;
    }
    
#ifdef PLATFORM_WINDOWS
    WSACleanup();
#endif
    
    platform_initialized = 0;
}

/*
 * 设置socket超时
 * 兼容所有平台的最可靠方法
 */
int set_socket_timeout(socket_t sock, int timeout_seconds) {
#ifdef PLATFORM_WINDOWS
    DWORD timeout_ms = timeout_seconds * 1000;
    
    if (setsockopt(sock, SOL_SOCKET, SO_RCVTIMEO, 
                   (char*)&timeout_ms, sizeof(timeout_ms)) != 0) {
        return -1;
    }
    
    if (setsockopt(sock, SOL_SOCKET, SO_SNDTIMEO, 
                   (char*)&timeout_ms, sizeof(timeout_ms)) != 0) {
        return -1;
    }
#else
    struct timeval tv;
    tv.tv_sec = timeout_seconds;
    tv.tv_usec = 0;
    
    if (setsockopt(sock, SOL_SOCKET, SO_RCVTIMEO, 
                   (void*)&tv, sizeof(tv)) != 0) {
        return -1;
    }
    
    if (setsockopt(sock, SOL_SOCKET, SO_SNDTIMEO, 
                   (void*)&tv, sizeof(tv)) != 0) {
        return -1;
    }
#endif
    
    return 0;
}

/*
 * 设置socket为非阻塞模式
 * 跨平台兼容实现
 */
int make_socket_nonblocking(socket_t sock) {
#ifdef PLATFORM_WINDOWS
    u_long mode = 1;
    return ioctlsocket(sock, FIONBIO, &mode);
#else
    int flags;
    
    flags = fcntl(sock, F_GETFL, 0);
    if (flags == -1) {
        return -1;
    }
    
    flags |= O_NONBLOCK;
    return fcntl(sock, F_SETFL, flags);
#endif
}