// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/16

package internal

import (
	_ "embed" // embed file for go version list
	"os"
	"path/filepath"

	"github.com/fsgo/cmdutils"
)

//go:embed files/golang_dl.tar.gz
var golangDlTar []byte

func extractGolangDLTar(dstDir string) error {
	tarPath := filepath.Join(DataDir(), "golang_dl.tar.gz")
	defer os.Remove(tarPath)

	if err := os.WriteFile(tarPath, golangDlTar, 0644); err != nil {
		return err
	}
	tr := &cmdutils.Tar{
		StripComponents: 1,
	}
	return tr.Unpack(tarPath, dstDir)
}
