// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/1

package internal

// Prepare 在其他正式命令之前的预处理逻辑
func Prepare() error {
	if err := ParserGOBIN(); err != nil {
		return err
	}

	loadConfig()
	printProxy()
	dataDir := DataDir()
	logPrint("data dir", dataDir)

	if err := chdir(dataDir); err != nil {
		return err
	}
	return Download()
}
