// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/1

package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Clean 将go1.x的老版本删除掉
func Clean(version string) error {
	versions, err := LastVersions()
	if err != nil {
		return err
	}

	vinfos := versions[version]
	if len(vinfos) == 0 {
		return fmt.Errorf("version %q not found", version)
	}

	log.Printf("%s has total %d versions, latest is %q\n", version, len(vinfos), vinfos[0].Raw)

	if len(vinfos) < 2 {
		log.Println("no old versions need to be clean")
		return nil
	}

	for i := 1; i < len(vinfos); i++ {
		cur := vinfos[i]
		if err = cleanVersion(cur); err != nil {
			log.Println("clean ", cur.Raw, "failed:", err)
		}
	}

	return nil
}

func cleanVersion(v *Version) error {
	sdkDir, err := goroot(v.Raw)
	if err != nil {
		return err
	}

	if _, err = os.Stat(sdkDir); err != nil && os.IsNotExist(err) {
		return nil
	}

	ignoreFile := filepath.Join(sdkDir, "smart-go-dl.ignore_clean")
	if _, err = os.Stat(ignoreFile); err == nil {
		log.Println("clean ignored")
		return nil
	}

	bin := filepath.Join(GOBIN(), v.Raw)
	log.Println("remove ", bin)
	if err = os.Remove(bin); err != nil && !os.IsNotExist(err) {
		return err
	}
	log.Println("remove ", sdkDir)
	if err = os.RemoveAll(sdkDir); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func goroot(version string) (string, error) {
	dir, err := sdkRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, version), nil
}
