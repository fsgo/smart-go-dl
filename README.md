# smart-go-dl
Go 多版本管理辅助工具, 可以快速安装 Go ( 次要版本 ) 的最新版本，并对过期版本进行清理。  

使用 https://github.com/golang/dl 以获取 Go 版本列表。

依赖：
 1. 需要设置环境变量 `$GOBIN`，可参考如下进行配置：
```bash
export GOBIN=$HOME/go/bin   # go install 安装的二进制文件所在目录，go1.x命令也将安装到此目录

export PATH=$GOBIN:$PATH    # 可以直接在任意位置使用 GOBIN 目录下的所有命令
```

## 安装/更新
未安装过 Go 的，请先在 https://go.dev/dl/ 下载安装 Go，
若非 windows 系统(如 Linux & mac )，也可以直接 [下载编译好的二进制文件](https://github.com/fsgo/smart-go-dl/releases) 。

已安装过 Go ，安装和更新：
```bash
go install github.com/fsgo/smart-go-dl@latest
```


## 查看使用帮助
```bash
smart-go-dl -help
```

## 安装 Go SDK
### 安装 `go1.25` 的最新版本：
```bash
smart-go-dl install go1.25
```
会自动找到`go1.25` 最新的版本进行安装，并安装为 `$GOBIN/go1.25.0` 和 `$GOBIN/go1.25`。  
如当前 go1.25 的最新版本是 `go1.25.5`，则上述 `$GOBIN/go1.22` 是 `$GOBIN/go1.22.5` 的软连接。  

安装或更新后，会创建软连 `$GOBIN/go.latest`，其为当前安装的最新版本。
若 `$GOBIN/go` 不存在，则也会创建这个软连接，相当于 `ln -s go.latest go`,即这个 `$GOBIN/go`总是最新版本的 go。

在使用的时候，可以直接使用 `go`、`go.latest`、`go1.25`、`go1.25.0` 之一：
```bash
# go1.25 version
或者
# go1.25.0 version
或者
# go version
```
输出：
```
go version go1.25.0 darwin/amd64
```

以后有新的版本了，重新使用 `smart-go-dl install/update go1.25` 即可安装最新版本。

使用其他版本示例：
```
go1.22.0 version       # 使用首个正式版本，对应版本号为 go1.22.0
go1.22.1 version       # 使用第 1 个正式修正版本,对应版本号为 go1.22.1
go1.22.2 version       # 使用第 2 个正式修正版本,对应版本号为 go1.22.2
```
### 安装指定的 3 位版本：
```bash
smart-go-dl install go1.22.5
```
### 安装首个正式版本
Go 的每个正式版本是如 `go1.22` 这种，3 位版本号 0 是缺省的，若要安装，可以这样：
```bash
smart-go-dl install go1.22.0
```
之后这样使用，如 `go1.22.0 version` 。


## 清理过期的 Go SDK
将 `go1.21` 除了最新版本的老版本清理掉：
```bash
smart-go-dl clean go1.21
```

若期望指定版本不被清理，可以使用子命令 `lock`，如下为让 `go1.22.3`这个版本不被清理：
```bash
smart-go-dl lock go1.22.3
```
于此对应的有 `unlock` 命令，用于解除 lock 状态。

## 更新 Go SDK
```bash
smart-go-dl update go1.22
```
等价于先执行 clean，再执行 install。  
还可以使用`smart-go-dl update` 来更新所有已安装版本( gotip 除外 )。


## 列出已安装/可安装的 Go SDK
```bash
smart-go-dl list
```

输出：
```
--------------------------------------------------------------------------------
version              latest               installed
--------------------------------------------------------------------------------
gotip                gotip
go1.23               go1.23rc2
go1.22               go1.22.5             go1.22.5
go1.21               go1.21.12
go1.20               go1.20.14
go1.19               go1.19.13
go1.18               go1.18.10
go1.17               go1.17.13
go1.16               go1.16.15
go1.15               go1.15.15
go1.14               go1.14.15
go1.13               go1.13.15
go1.12               go1.12.17
go1.11               go1.11.13
go1.10               go1.10.8
go1.9                go1.9.7
go1.8                go1.8.7
go1.7                go1.7.6
go1.6                go1.6.4
go1.5                go1.5.4
[smart-go-dl] list success
```

第一列，若是绿色，说明当前已按照最新版本，若是黄色，安装的不是最新版本。    
windows 环境下目前未做终端颜色的适配。  

## 删除指定版本的 Go SDK
```bash
smart-go-dl remove go1.19.1
```

## 配置文件
可选的配置文件为 `~/.config/smart-go-dl/app.toml`:
```toml
# 下载时使用的 Proxy，可选
# 不配置或者为空时，会使用环境变量的代理配置
# Proxy="http://127.0.0.1:8128"

# 下载文件时，是否跳过证书校验，可选，默认 false
# InsecureSkipVerify = true

# 下载 Go tar 文件的地址前缀，可选
# 会一次使用每个地址进行尝试
# 默认值是 "https://dl.google.com/go/,https://dl-ssl.google.com/go/"
#TarURLPrefix="https://dl.google.com/go/"

# 安装目录，可选，默认为 ~/sdk
# 不同的 Go 版本在 SDKDir 中以子目录方式存在，如 ~/sdk/go1.22.0/
# SDKDir = ""
```
该文件在不存在的时候，会尝试自动创建

## 数据/缓存目录
该程序使用 `${SDKDir}/smart-go-dl/` 目录缓存数据，依赖的 https://github.com/golang/dl 
也会自动下载到此目录下的 `golang_dl` 子目录中。  
首次使用时会使用 `git clone` 命令下载 `golang_dl`，之后会使用 `git pull` 命令检查更新。  
因 golang_dl 更新频率很低，也为了使用 `smart-go-dl` 时更流畅，更新时间间隔在 1 分钟内，
再次使用时不会使用 `git pull` 检查更新。  
若因为某些原因，git 命令下载和更新不能正常工作，也可以手工创建和更新该目录。


## 自动版本选择
在不同目录，执行 go 命令，使用不同的 go 版本：  
https://github.com/fsgo/bin-auto-switcher