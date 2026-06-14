# FScan 插件开发规范

## 概述

FScan 采用简化的单文件插件架构，每个插件一个 `.go` 文件，消除了过度设计的多文件结构。


1. **简洁至上**：消除所有不必要的抽象层
2. **直击本质**：专注于解决实际问题，不为架构而架构  
3. **向后兼容**：不破坏用户接口和现有功能
4. **消除特殊情况**：统一处理逻辑，减少 if/else 分支

## 插件架构

### 核心接口

```go
// Plugin 插件接口 - 只保留必要的方法
type Plugin interface {
    GetName() string                                        // 插件名称
    GetPorts() []int                                       // 支持的端口
    Scan(ctx context.Context, info *common.HostInfo) *ScanResult  // 扫描功能
}

// 可选接口：如果插件支持利用功能
type Exploiter interface {
    Exploit(ctx context.Context, info *common.HostInfo, creds Credential, config *common.Config) *ExploitResult
}
```

### 数据结构

```go
// ScanResult 扫描结果 - 删除所有冗余字段
type ScanResult struct {
    Success  bool   // 扫描是否成功
    Service  string // 服务类型
    Username string // 发现的用户名（弱密码）
    Password string // 发现的密码（弱密码）  
    Banner   string // 服务版本信息
    Error    error  // 错误信息（如果失败）
}

// ExploitResult 利用结果（仅有利用功能的插件需要）
type ExploitResult struct {
    Success bool   // 利用是否成功
    Output  string // 命令执行输出
    Error   error  // 错误信息
}

// Credential 凭据结构
type Credential struct {
    Username string
    Password string
    KeyData  []byte // SSH私钥等
}
```

## 插件开发模板

### 1. 纯扫描插件（如MySQL）

```go
package plugins

import (
    "context"
    "fmt"
    // 其他必要导入
)

// PluginName服务扫描插件
type PluginNamePlugin struct {
    name  string
    ports []int
}

// 构造函数
func NewPluginNamePlugin() *PluginNamePlugin {
    return &PluginNamePlugin{
        name:  "plugin_name",
        ports: []int{default_port},
    }
}

// 实现Plugin接口
func (p *PluginNamePlugin) GetName() string { return p.name }
func (p *PluginNamePlugin) GetPorts() []int { return p.ports }

func (p *PluginNamePlugin) Scan(ctx context.Context, info *common.HostInfo) *ScanResult {
    // 如果禁用暴力破解，只做服务识别
    if common.DisableBrute {
        return p.identifyService(info)
    }

    // 生成测试凭据
    credentials := GenerateCredentials("plugin_name")
    
    // 逐个测试凭据
    for _, cred := range credentials {
        select {
        case <-ctx.Done():
            return &ScanResult{Success: false, Error: ctx.Err()}
        default:
        }

        if p.testCredential(ctx, info, cred) {
            return &ScanResult{
                Success:  true,
                Service:  "plugin_name",
                Username: cred.Username,
                Password: cred.Password,
            }
        }
    }

    return &ScanResult{Success: false, Service: "plugin_name"}
}

// 核心认证逻辑
func (p *PluginNamePlugin) testCredential(ctx context.Context, info *common.HostInfo, cred Credential) bool {
    // 实现具体的认证测试逻辑
    return false
}

// 服务识别（-nobr模式）
func (p *PluginNamePlugin) identifyService(info *common.HostInfo) *ScanResult {
    // 实现服务识别逻辑
    return &ScanResult{Success: false, Service: "plugin_name"}
}

// 自动注册
func init() {
    RegisterPlugin("plugin_name", func() Plugin {
        return NewPluginNamePlugin()
    })
}
```

### 2. 带利用功能的插件（如SSH）

```go
package plugins

// SSH插件结构
type SSHPlugin struct {
    name  string
    ports []int
}

// 同时实现Plugin和Exploiter接口
func (p *SSHPlugin) Scan(ctx context.Context, info *common.HostInfo) *ScanResult {
    // 扫描逻辑（同上）
}

func (p *SSHPlugin) Exploit(ctx context.Context, info *common.HostInfo, creds Credential, config *common.Config) *ExploitResult {
    // 建立SSH连接
    client, err := p.connectSSH(info, creds)
    if err != nil {
        return &ExploitResult{Success: false, Error: err}
    }
    defer client.Close()

    // 执行命令或其他利用操作
    output, err := p.executeCommand(client, "whoami")
    return &ExploitResult{
        Success: err == nil,
        Output:  output,
        Error:   err,
    }
}

// 辅助方法
func (p *SSHPlugin) connectSSH(info *common.HostInfo, creds Credential) (*ssh.Client, error) {
    // SSH连接实现
}

func (p *SSHPlugin) executeCommand(client *ssh.Client, cmd string) (string, error) {
    // 命令执行实现
}
```

## 开发规范

### 文件组织

```
plugins/
├── base.go           # 核心接口和注册系统
├── mysql.go          # MySQL插件
├── ssh.go            # SSH插件  
├── redis.go          # Redis插件
└── README.md         # 开发文档（本文件）
```

### 命名规范

- **插件文件**：`{service_name}.go`
- **插件结构体**：`{ServiceName}Plugin`
- **构造函数**：`New{ServiceName}Plugin()`
- **插件名称**：小写，与文件名一致

### 代码规范

1. **错误处理**：始终使用Context进行超时控制
2. **日志输出**：成功时使用 `common.LogSuccess`，调试用 `common.LogDebug`
3. **凭据生成**：使用 `GenerateCredentials(service_name)` 生成测试凭据
4. **资源管理**：及时关闭连接，使用 defer 确保清理

### 测试要求

每个插件必须支持：

1. **暴力破解模式**：`common.DisableBrute = false`
2. **服务识别模式**：`common.DisableBrute = true`
3. **Context超时处理**：正确响应 `ctx.Done()`
4. **代理支持**：如果 `common.Socks5Proxy` 不为空

## 迁移指南

### 从三文件架构迁移

1. **提取核心逻辑**：从 connector.go 提取认证逻辑
2. **合并实现**：将 plugin.go 中的组装逻辑内联
3. **删除垃圾**：删除空的 exploiter.go
4. **简化数据结构**：只保留必要的字段

### 从Legacy插件迁移

1. **保留核心逻辑**：复制扫描和认证的核心算法
2. **标准化接口**：实现统一的Plugin接口  
3. **移除全局依赖**：通过返回值而不是全局变量传递结果
4. **统一日志**：使用统一的日志接口

## 性能优化

1. **连接复用**：在同一次扫描中复用连接
2. **内存管理**：及时释放不需要的资源
3. **并发控制**：通过Context控制并发度
4. **超时设置**：合理设置各阶段超时时间

## 示例

参考 `mysql.go` 作为标准的纯扫描插件实现
参考 `ssh.go` 作为带利用功能的插件实现

---

**记住：好的代码不是写出来的，是重构出来的。消除所有不必要的复杂性，直击问题本质。**