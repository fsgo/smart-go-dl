# smart-go-dl
Go 多版本管理辅助工具, 可以快速安装 Go ( 次要版本 ) 的最新版本，并对过期版本进行清理。  

底层使用 https://github.com/golang/dl 来进行多 Go 版本的安装。

依赖：
 1. 安装过 git，内部使用了 `git clone` 和 `git pull` 命令
 2. 需要设置环境变量 `$GOBIN`，可参考如下进行配置：
```bash
export GOBIN=$HOME/go/bin   # go install 安装的二进制文件所在目录，go1.x命令也将安装到此目录

export PATH=$PATH:$GOBIN    # 可以直接在任意位置使用 GOBIN 目录下的所有命令
```
若之前安装的 `go1.x` 命令(`go`命令不受影响)不在上述 `$GOBIN` 路径里，请删除掉或者移动到 `$GOBIN` 里，
以避免使用`smart-go-dl` install 或者 clean 后，使用命令 `go1.x`(如 go.16) 使用的是旧版本的。

## 安装/更新
未安装过 Go 的，请先在 https://go.dev/dl/ 下载安装 Go，
若非 windows 系统(如 Linux & mac )，也可以直接 [下载编译好的二进制文件](https://github.com/fsgo/smart-go-dl/releases) 。

已安装过 Go ，安装和更新：
```bash
go install github.com/fsgo/smart-go-dl@main
```


## 查看使用帮助
```bash
smart-go-dl -help
```

## 安装 Go SDK
### 安装最新的 `go1.18`：
```bash
smart-go-dl install go1.18
```
会自动找到`go1.18` 最新的版本进行安装，并安装为 `$GOBIN/go1.18beta1` 和 `$GOBIN/go1.18`。  
如当前 go1.18 的最新版本是 `go1.18beta1`，则上述 `$GOBIN/go1.18` 是 `$GOBIN/go1.18beta1` 的软连接。

在使用的时候，可以直接使用 `go1.18` 或者 `go1.18beta1`：
```bash
go1.18 version           # 使用的总是 go1.18 系列的最新版本
# go1.18beta1 version    # 使用的是 指定的小版本
```
输出：
```
# go version go1.18beta1 darwin/amd64
```

以后有新的版本了，重新使用 `smart-go-dl install/update go1.18` 即可安装最新版本。

使用其他版本示例：
```
go1.18.0 version       # 使用首个正式版本，对应版本号为 go1.18
go1.18.1 version       # 使用第 1 个正式修正版本,对应版本号为 go1.18.1
go1.18.2 version       # 使用第 2 个正式修正版本,对应版本号为 go1.18.2
```
### 安装指定的 3 位版本：
```bash
smart-go-dl install go1.17.3
```
### 安装首个正式版本
Go 的每个正式版本是如 `go1.17` 这种，3 位版本号 0 是缺省的，若要安装，可以这样：
```bash
smart-go-dl install go1.17.0
```
之后这样使用，如 `go1.17.0 version` 。


## 清理过期的 Go SDK
将 `go1.17` 除了最新版本的老版本清理掉：
```bash
smart-go-dl clean go1.17
```

若期望指定版本不被清理，可以使用子命令 `lock`，如下为让 `go1.17.3`这个版本不被清理：
```bash
smart-go-dl lock go1.17.3
```
于此对应的有 `unlock` 命令，用于解除 lock 状态。

## 更新 Go SDK
```bash
smart-go-dl update go1.17
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
go1.18               go1.18beta1          go1.18beta1
go1.17               go1.17.5             go1.17.5.0     go1.17.3
go1.16               go1.16.12            go1.16.12
go1.15               go1.15.15            go1.15.15
go1.14               go1.14.15            go1.14.0
go1.13               go1.13.15            go1.13.0
go1.12               go1.12.17            go1.12.17
go1.11               go1.11.13            go1.11.13
go1.10               go1.10.8             go1.10.8
go1.9                go1.9.7              go1.9.7
go1.8                go1.8.7
go1.7                go1.7.6
go1.6                go1.6.4
go1.5                go1.5.4
```

第一列，若是绿色，说明当前已按照最新版本，若是黄色，安装的不是最新版本。    
windows 环境下目前未做终端颜色的适配。  

## 删除指定版本的 Go SDK
```bash
smart-go-dl remove go1.17.3
```

## 数据/缓存目录
该程序使用 `$HOME/sdk/smart-go-dl/` 目录缓存数据，依赖的 https://github.com/golang/dl 
也会自动下载到此目录下的 `golang_dl` 子目录中。  
首次使用时会使用 `git clone` 命令下载 `golang_dl`，之后会使用 `git pull` 命令检查更新。  
因 golang_dl 更新频率很低，也为了使用 `smart-go-dl` 时更流畅，更新时间间隔在 1 小时内，
再次使用时不会使用 `git pull` 检查更新。  
若因为某些原因，git 命令下载和更新不能正常工作，也可以手工创建和更新该目录。


## 自动版本选择
在不同目录，执行 go 命令，使用不同的 go 版本：  
https://github.com/fsgo/bin-auto-switcher