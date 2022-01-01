// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/1

package internal

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const lockedName = "smart-go-dl.locked"

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
	return ioutil.WriteFile(name, []byte("clean locked"), 0655)
}
