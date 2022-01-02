// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/1

package internal

import (
	"log"
	"os"
)

func Prepare() error {
	tmpDir, err := InitTmpDir()
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

	if err = os.Chdir(tmpDir); err != nil {
		return err
	}
	return Download()
}
