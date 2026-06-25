# fscan-lite

极简但极致兼容的TCP内网端口扫描器

## 设计理念

**兼容性第一，简洁至上**

- 支持从 Windows 98 到 Windows 11
- 支持从 Ubuntu 8.04 到最新版本  
- 使用 C89 标准，最大兼容性
- 静态编译，零依赖运行
- 单个可执行文件 < 1MB

## 功能特性

- ✅ TCP端口连接扫描
- ✅ 支持端口范围 (1-65535, 80,443)
- ✅ 可配置超时时间
- ✅ 静态编译，零依赖
- ✅ 跨平台兼容

## 编译

### Linux/Unix

```bash
# 动态编译
make

# 静态编译（推荐）
make static

# 最小化编译
make small
```

### Windows

```bash
# MinGW 编译
mingw32-make -f Makefile

# 或使用MSVC
cl /TC src/*.c /Febin/fscan-lite.exe ws2_32.lib
```

## 使用方法

```bash
# 基本用法
./bin/fscan-lite -h 192.168.1.1 -p 22,80,443

# 扫描端口范围
./bin/fscan-lite -h 10.0.0.1 -p 1-1000

# 自定义超时
./bin/fscan-lite -h 192.168.1.100 -p 80,443 -t 2
```

## 参数说明

| 参数 | 说明 | 示例 |
|------|------|------|
| -h HOST | 目标主机IP | -h 192.168.1.1 |
| -p PORTS | 端口列表 | -p 80,443,8000-8080 |
| -t TIMEOUT | 超时时间（秒） | -t 3 |
| --help | 显示帮助 | --help |
| --version | 显示版本 | --version |

## 二进制大小对比

| 版本 | 大小 | 说明 |
|------|------|------|
| fscan (Go) | ~30MB | 包含运行时 |
| fscan-lite | ~900KB | 静态编译 |
| fscan-lite (strip) | ~700KB | 去除调试信息 |

## 兼容性测试

### Linux 发行版
- ✅ Ubuntu 8.04 - 24.04
- ✅ CentOS 5 - 9  
- ✅ Debian 5 - 12
- ✅ RHEL 5 - 9

### Windows 版本
- ✅ Windows 98 SE
- ✅ Windows XP
- ✅ Windows 7/8/10/11
- ✅ Windows Server 2003-2022

## 技术实现

- **语言**: C89 (最大兼容性)
- **网络**: 原生socket API
- **编译**: GCC/MSVC/Clang
- **链接**: 静态链接，零依赖
- **大小**: < 1MB 单文件

## 性能对比

| 指标 | fscan | fscan-lite |
|------|-------|------------|
| 启动时间 | ~50ms | ~5ms |
| 内存占用 | ~20MB | ~2MB |
| 扫描速度 | 1000 ports/s | 1000 ports/s |
| 兼容性 | 现代系统 | 25年跨度 |

## 构建配置

```bash
# 查看构建信息
make info

# 所有构建选项
make help
```

## 许可证

与 fscan 主项目保持一致

---

**理念**: 一个工具应该在它设计的任何系统上都能运行，而不需要用户去寻找依赖项。