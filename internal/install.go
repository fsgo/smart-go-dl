// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/12/31

package internal

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Install 安装 go1.x 的最新版本
func Install(version string) error {
	versions, err := LastVersions()
	if err != nil {
		return err
	}
	vinfos := versions[version]
	if len(vinfos) == 0 {
		return fmt.Errorf("version %q not found", version)
	}
	last := vinfos[0]

	log.Println("[install]", "found last", version, "version is", last.Raw)

	defer os.Chdir(TmpDir())

	if err = os.Chdir(last.DlDir()); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	goBinTo := filepath.Join(GOBIN(), last.Raw)
	cmd := exec.CommandContext(ctx, "go", "build", "-o", goBinTo)
	log.Println("[exec]", cmd.String())
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		return err
	}

	downloadCmd := exec.Command(goBinTo, "download")
	log.Println("[exec]", downloadCmd.String())
	downloadCmd.Stderr = os.Stderr
	downloadCmd.Stdout = os.Stdout
	if err = downloadCmd.Run(); err != nil {
		return err
	}

	goBinLink := filepath.Join(GOBIN(), last.Normalized)
	if goBinLink == goBinTo {
		return nil
	}

	// create link for go bin
	// go1.16.6 -> go1.16
	{
		if err = os.Remove(goBinLink); err != nil && !os.IsNotExist(err) {
			return err
		}

		if err = os.Symlink(goBinTo, goBinLink); err != nil {
			return err
		}
		log.Println("[link]", goBinTo, "->", goBinLink, "success")
	}

	// create sdk dir link
	{
		sdkDir := last.Raw
		sdkDirLink := mustGoRoot(last.Normalized)

		if err = os.Remove(sdkDirLink); err != nil && !os.IsNotExist(err) {
			return err
		}

		if err = os.Symlink(sdkDir, sdkDirLink); err != nil {
			return err
		}
		log.Println("[link]", mustGoRoot(last.Raw), "->", sdkDirLink, "success")
	}
	log.Printf("Success. You may now run '%s'\n", version)
	return nil
}
