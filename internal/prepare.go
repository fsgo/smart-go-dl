// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/1

package internal

import (
	"log"
	"os"
)

// Prepare 在其他正式命令之前的预处理逻辑
func Prepare() error {
	tmpDir, err := GetTmpDir()
	if err != nil {
		return err
	}

	if err = ParserGOBIN(); err != nil {
		return err
	}

	log.Println("Use TmpDir:", tmpDir)

	if err = os.MkdirAll(tmpDir, 0777); err != nil && !os.IsExist(err) {
		return err
	}

	loadConfig()
	printProxy()

	if err = chdir(tmpDir); err != nil {
		return err
	}
	return Download()
}
