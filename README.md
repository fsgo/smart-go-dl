# smart-go-dl
go 多版本管理辅助工具, 可以快速安装 go1.x 版本的最新版本，并对过期版本进行清理。  

依赖 https://github.com/golang/dl

请提前设置好环境变量 `$GOBIN`,若没有设置不能正常工作。

## 安装/更新
```bash
go install github.com/fsgo/smart-go-dl@main
```

## 查看使用帮助
```
smart-go-dl -help
```

## 安装 go sdk
如下为安装最新的 `go1.18`：
```bash
smart-go-dl install go1.18
```
会自动找到最新的 `go1.18` 并进行安装，安装到 $GOBIN/`go1.18`。  
如当前 go1.18 的最新版本是 `go1.18beta1`，则上述 $GOBIN/`go1.18` 是 $GOBIN/`go1.18beta1` 的软连接。

在使用的时候，可以直接使用 go1.18 即可：
```bash
# go1.18 version
```
go version go1.18beta1 darwin/amd64


## 清理过期的 go sdk
将 `go1.17` 除了最新版本的老版本清理掉：
```bash
smart-go-dl clean go1.17
```

## 更新 go sdk
```bash
smart-go-dl clean go1.17
```
等价于先执行 install，再执行 clean。

## 列出已安装/可按照的 go sdk
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