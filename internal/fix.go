// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2022/11/29

package internal

import (
	"context"
	"os"
	"path/filepath"
)

func Fix(ctx context.Context) error {
	return installGoLatestBin(ctx)
}

func installGoLatestBin(ctx context.Context) error {
	versions, err := LastVersions(ctx)
	if err != nil {
		return err
	}
	var latest *Version
	var def *Version
	for _, mv := range versions {
		for _, z := range mv.PatchVersions {
			if z.Raw == "gotip" {
				continue
			}
			if z.Installed() {
				if z.IsNormal() && (latest == nil || z.Num > latest.Num) {
					latest = z
				}
				if def == nil {
					def = z
				}
			}
		}
	}

	if latest == nil {
		latest = def
	}

	if latest == nil {
		return nil
	}
	latest.NormalizedGoBinPath()

	latestBinPath := filepath.Join(GOBIN(), "go.latest"+exe())

	if err1 := createLink(latest.NormalizedGoBinPath(), latestBinPath); err1 != nil {
		return err1
	}

	// 若是 $GOBIN/go 不存在，则创建一个软连接
	goPath := filepath.Join(GOBIN(), "go"+exe())
	if _, err2 := os.Stat(goPath); os.IsNotExist(err2) {
		_ = createLink(latestBinPath, goPath)
	}

	return nil
}
