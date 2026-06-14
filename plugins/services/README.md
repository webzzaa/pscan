# 服务扫描插件目录

本目录包含所有服务扫描插件，采用简化的单文件插件架构。

## 已实现插件

### 数据库服务
- `mysql.go` - MySQL数据库扫描
- `postgresql.go` - PostgreSQL数据库扫描  
- `redis.go` - Redis内存数据库扫描
- `mongodb.go` - MongoDB文档数据库扫描
- `mssql.go` - Microsoft SQL Server扫描
- `oracle.go` - Oracle数据库扫描
- `memcached.go` - Memcached缓存扫描
- `neo4j.go` - Neo4j图数据库扫描

### 消息队列服务
- `rabbitmq.go` - RabbitMQ消息队列扫描
- `activemq.go` - ActiveMQ消息队列扫描
- `kafka.go` - Apache Kafka扫描

### 网络服务
- `ssh.go` - SSH远程登录服务扫描
- `ftp.go` - FTP文件传输服务扫描
- `telnet.go` - Telnet远程终端服务扫描
- `smtp.go` - SMTP邮件服务扫描
- `snmp.go` - SNMP网络管理协议扫描
- `ldap.go` - LDAP目录服务扫描
- `rsync.go` - Rsync文件同步服务扫描

### Windows服务
- `findnet.go` - Windows网络发现插件 (RPC端点映射)
- `smbinfo.go` - SMB协议信息收集插件

### 其他服务
- `vnc.go` - VNC远程桌面服务扫描
- `cassandra.go` - Apache Cassandra数据库扫描

## 插件特性

每个插件都包含：
- ✅ 服务识别功能
- ✅ 弱密码检测功能
- ✅ 完整的利用功能
- ✅ 错误处理和超时控制
- ✅ 统一的结果输出格式

## 开发规范

所有插件都遵循 `../README.md` 中定义的开发规范。