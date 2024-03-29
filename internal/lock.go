// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/1

package internal

import (
	"os"
	"path/filepath"
)

const lockedName = "smart-go-dl.locked"

// Lock 给指定版本添加 lock 标记文件
func Lock(version string, action string) error {
	sdk, err := goroot(version)
	if err != nil {
		return err
	}
	name := filepath.Join(sdk, lockedName)
	if action == "remove" {
		if err = os.Remove(name); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	return os.WriteFile(name, []byte("clean locked"), 0655)
}

func isLocked(version string) bool {
	sdk, err := goroot(version)
	if err != nil {
		return false
	}
	name := filepath.Join(sdk, lockedName)
	_, err = os.Stat(name)
	return err == nil
}
