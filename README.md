# smart-go-dl
Go 多版本管理辅助工具, 可以快速安装 Go ( 次要版本 ) 的最新版本，并对过期版本进行清理。  

底层使用 https://github.com/golang/dl 来进行多 Go 版本的安装。

依赖：
 1. git 工具
 2. 需要设置环境变量 `$GOBIN`，可参考如下：
```bash
export GOBIN=$HOME/go/bin

export PATH=$PATH:$GOBIN
```
注意：
若之前按照 go1.x 不在上述 $GOBIN 路径里，请删除掉，以避免使用`smart-go-dl` install 或者 clean 后不生效。

## 安装/更新
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
会自动找到`go1.18` 最新的版本进行安装，go 命令安装为 $GOBIN/`go1.18`。  
如当前 go1.18 的最新版本是 `go1.18beta1`，则上述 $GOBIN/`go1.18` 是 $GOBIN/`go1.18beta1` 的软连接。

在使用的时候，可以直接使用 `go1.18` 即可：
```bash
go1.18 version

# go version go1.18beta1 darwin/amd64
```
以后有新的版本了，重新使用 `smart-go-dl install/update go1.18` 即可安装最新版本，
go 的命令依旧保持为 `go1.18` 不变。

### 安装指定的 3 位版本：
```
smart-go-dl install go1.17.3
```
### 安装首个正式版本
Go 的每个正式版本是如 `go1.17` 这种，3 位版本号 0 是缺省的，若要安装，可以这样：
```
smart-go-dl install go1.17.0
```
之后这样使用，如 `go1.17.0 version` 。



## 清理过期的 Go SDK
将 `go1.17` 除了最新版本的老版本清理掉：
```bash
smart-go-dl clean go1.17
```

若期望指定版本不被清理，可以使用子命令 `lock`，如下为让 `go1.17.3`这个版本不被清理：
```
smart-go-dl lock go1.17.3
```
于此对应的有 `unlock` 命令，用于解除 lock 状态。

## 更新 Go SDK
```bash
smart-go-dl update go1.17
```
等价于先执行 clean，再执行 install。  
还可以使用`smart-go-dl update all` 来更新所有已安装版本( gotip 除外 )。


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
go1.17               go1.17.5             go1.17.5     go1.17.3
go1.16               go1.16.12            go1.16.12
go1.15               go1.15.15            go1.15.15
go1.14               go1.14.15            go1.14
go1.13               go1.13.15            go1.13
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


## 自动版本选择
在不同目录，执行 go 命令，使用不同的 go 版本：  
https://github.com/fsgo/bin-auto-switcher