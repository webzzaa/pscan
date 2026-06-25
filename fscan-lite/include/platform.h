#ifndef PLATFORM_H
#define PLATFORM_H

/* 
 * platform.h - 极致兼容性的平台抽象层
 * 
 * 支持范围：
 * Windows: 98/ME/NT4/2000/XP/Vista/7/8/10/11
 * Linux:   glibc 2.3+ (2003年后的所有发行版)
 * 编译器:  MSVC 6.0+, GCC 3.0+, Clang 3.0+
 */

/* C89兼容性 - 最古老但最可靠的标准 */
#ifndef __STDC__
#define __STDC__ 1
#endif

/* 平台检测 */
#ifdef _WIN32
    #define PLATFORM_WINDOWS
    #ifdef _WIN64
        #define PLATFORM_WIN64
    #else
        #define PLATFORM_WIN32
    #endif
#else
    #define PLATFORM_UNIX
    #ifdef __linux__
        #define PLATFORM_LINUX
    #elif defined(__APPLE__)
        #define PLATFORM_MACOS
    #endif
#endif

/* Windows头文件包含 - 兼容Win98 */
#ifdef PLATFORM_WINDOWS
    /* 定义最低Windows版本 - Win98 */
    #ifndef _WIN32_WINNT
        #define _WIN32_WINNT 0x0410  /* Windows 98 */
    #endif
    #ifndef WINVER
        #define WINVER 0x0410
    #endif

    /* 必须先包含winsock2.h，否则windows.h会包含旧版winsock.h导致冲突 */
    #include <winsock2.h>
    #include <windows.h>

    /* 老版本Windows兼容性 */
    #ifdef _MSC_VER
        #if _MSC_VER < 1300  /* MSVC 6.0 */
            #pragma comment(lib, "wsock32.lib")
        #else
            #pragma comment(lib, "ws2_32.lib")
        #endif
    #endif

    /* Windows类型定义 */
    typedef SOCKET socket_t;
    typedef int socklen_t;
    #define INVALID_SOCKET_VALUE INVALID_SOCKET
    #define close_socket closesocket
    #define socket_errno WSAGetLastError()

    /* Windows错误码转换 - 使用ifndef避免与errno.h冲突 */
    #ifndef EWOULDBLOCK
        #define EWOULDBLOCK WSAEWOULDBLOCK
    #endif
    #ifndef EINPROGRESS
        #define EINPROGRESS WSAEINPROGRESS
    #endif
    #ifndef ECONNREFUSED
        #define ECONNREFUSED WSAECONNREFUSED
    #endif
    
#else
    /* Unix/Linux头文件 */
    #include <sys/types.h>
    #include <sys/socket.h>
    #include <sys/time.h>
    #include <netinet/in.h>
    #include <arpa/inet.h>
    #include <unistd.h>
    #include <errno.h>
    #include <netdb.h>
    
    /* Unix类型定义 */
    typedef int socket_t;
    #define INVALID_SOCKET_VALUE (-1)
    #define close_socket close
    #define socket_errno errno
    
#endif

/* 标准C头文件 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

/* 线程抽象 - 最简单的实现 */
#ifdef PLATFORM_WINDOWS
    typedef HANDLE thread_t;
    typedef DWORD thread_id_t;
    typedef unsigned (__stdcall *thread_func_t)(void *);
    
    #define CREATE_THREAD(func, arg) \
        (HANDLE)_beginthreadex(NULL, 0, (thread_func_t)(func), (arg), 0, NULL)
    #define WAIT_THREAD(handle) WaitForSingleObject((handle), INFINITE)
    #define CLOSE_THREAD(handle) CloseHandle(handle)
    
#else
    #include <pthread.h>
    #include <fcntl.h>
    typedef pthread_t thread_t;
    typedef pthread_t thread_id_t;
    typedef void* (*thread_func_t)(void *);
    
    #define CREATE_THREAD(func, arg) ({ \
        pthread_t t; \
        pthread_create(&t, NULL, (thread_func_t)(func), (arg)) == 0 ? t : 0; \
    })
    #define WAIT_THREAD(handle) pthread_join((handle), NULL)
    #define CLOSE_THREAD(handle) /* pthread handles are auto-cleaned */
    
#endif

/* 时间函数抽象 */
#ifdef PLATFORM_WINDOWS
    #define sleep_ms(ms) Sleep(ms)
#else
    #define sleep_ms(ms) usleep((ms) * 1000)
#endif

/* 编译器特定定义 */
#ifdef _MSC_VER
    /* MSVC特定 */
    #define snprintf _snprintf
    #define vsnprintf _vsnprintf
    #define strcasecmp _stricmp
    #define strncasecmp _strnicmp
#endif

/* 常用常量 */
#define MAX_HOST_LEN 256
#define MAX_PORT_COUNT 65536
#define DEFAULT_TIMEOUT 3
#define DEFAULT_THREAD_COUNT 100

/* 函数声明 */
int platform_init(void);
void platform_cleanup(void);
int set_socket_timeout(socket_t sock, int timeout_seconds);
int make_socket_nonblocking(socket_t sock);

#endif /* PLATFORM_H */