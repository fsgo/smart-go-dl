// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/11/29

package internal

import (
	"os"
	"path/filepath"
)

func Fix() error {
	return installGoLatestBin()
}

func installGoLatestBin() error {
	versions, err := LastVersions()
	if err != nil {
		return err
	}
	var latest *Version
	for _, mv := range versions {
		for _, z := range mv.PatchVersions {
			if z.Raw == "gotip" {
				continue
			}
			if z.Installed() && (latest == nil || z.Num > latest.Num) {
				latest = z
			}
		}
	}
	if latest == nil {
		return nil
	}
	latest.NormalizedGoBinPath()
	latestBinPath := filepath.Join(GOBIN(), "go.latest")

	if err1 := createLink(latest.NormalizedGoBinPath(), latestBinPath); err1 != nil {
		return err1
	}

	// 若是 $GOBIN/go 不存在，则创建一个软连接
	goPath := filepath.Join(GOBIN(), "go")
	if _, err2 := os.Stat(goPath); os.IsNotExist(err2) {
		_ = createLink(latestBinPath, goPath)
	}

	// 给最后的
	latestGoRoot, _ := goroot("go1.latest")
	if len(latestGoRoot) > 0 {
		_ = createLink(latest.GOROOT(), latestGoRoot)
	}

	return nil
}
