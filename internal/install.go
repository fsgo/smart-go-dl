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
		// 用于支持安装 3 位版本，如  go1.16.0、go1.16.3
		err = installVV(version, versions)
		if err == nil {
			return nil
		}
		return fmt.Errorf("install %q failed: %w", version, err)
	}
	last := vinfos[0]

	log.Println("[install]", "found last", version, "version is", last.Raw)

	goBinTo := last.RawGoBinPath()

	if err = installWithVersion(last); err != nil {
		return err
	}

	goBinLink := last.NormalizedGoBinPath()
	if goBinLink == goBinTo {
		return nil
	}

	// create link for go bin
	// go1.16.6 -> go1.16
	{
		if err = os.Remove(goBinLink); err != nil && !os.IsNotExist(err) {
			return err
		}
		if isWindows() {
			if err = copyFile(goBinTo, goBinLink); err != nil {
				return err
			}
		} else {
			if err = os.Symlink(goBinTo, goBinLink); err != nil {
				return err
			}
		}
		log.Println("[link]", goBinTo, "->", goBinLink, "success")
	}

	// create sdk dir link
	if !isWindows() {
		sdkDir := last.Raw
		sdkDirLink := mustGoRoot(last.Normalized) + ".latest"

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

func installWithVersion(ver *Version) error {
	defer os.Chdir(TmpDir())
	log.Println("[chdir]", ver.DlDir())
	if err := os.Chdir(ver.DlDir()); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	goBinTo := ver.RawGoBinPath()
	cmd := exec.CommandContext(ctx, "go", "build", "-o", goBinTo)
	log.Println("[exec]", cmd.String())
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}

	downloadCmd := exec.Command(goBinTo, "download")
	log.Println("[exec]", downloadCmd.String())
	downloadCmd.Stderr = os.Stderr
	downloadCmd.Stdout = os.Stdout
	return downloadCmd.Run()
}

// installVV 安装指定的小版本
func installVV(version string, vvs map[string][]*Version) error {
	if len(vvs[version]) > 0 {
		// 不应该执行到这个逻辑
		return fmt.Errorf("now allow, bug here")
	}
	vu, err := parserVersion(version)
	if err != nil {
		return err
	}
	mvs := vvs[vu.Normalized]
	if len(mvs) == 0 {
		return fmt.Errorf("minor version not found")
	}
	var installVersion *Version
	for _, mv := range mvs {
		if mv.Raw == version || mv.Raw+".0" == version {
			installVersion = mv
			break
		}
	}
	if installVersion == nil {
		return fmt.Errorf("version not found")
	}
	return installWithVersion(installVersion)
}
