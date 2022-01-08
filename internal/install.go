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
	"runtime"
	"sort"
	"strings"
	"time"
)

// Install 安装 go1.x 的最新版本
func Install(version string) error {
	versions, err := LastVersions()
	if err != nil {
		return err
	}
	mv := versions.Get(version)
	if mv == nil {
		// 用于支持安装 3 位版本，如  go1.16.0、go1.16.3
		err = installVV(version, versions)
		if err == nil {
			return nil
		}
		return fmt.Errorf("install %q failed: %w", version, err)
	}
	last := mv.Latest()

	log.Printf("[install] found %s's latest version is %s\n", version, last.Raw)

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
	gb, err := findGoBin()
	if err != nil {
		if !isWindows() {
			// 当没有找到 go 的时候，尝试直接使用下载编译好的 go
			err = installByArchive(ver.Raw)
		}
		if err != nil {
			return err
		}
	}

	gb, err = findGoBin()
	if err != nil {
		return err
	}

	if err = chdir(ver.DlDir()); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	goBinTo := ver.RawGoBinPath()

	cmd := exec.CommandContext(ctx, gb, "build", "-o", goBinTo)
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
	err = downloadCmd.Run()
	if err == nil {
		removeGoTmpTar(ver.Raw)
		log.Printf("Success. You may now run '%s'\n", filepath.Base(goBinTo))
	}
	return err
}

func removeGoTmpTar(version string) {
	sdk, err := sdkRoot()
	if err != nil {
		return
	}
	tmpTar := filepath.Join(sdk, version, version+"."+runtime.GOOS+"-"+runtime.GOARCH+".tar.gz")
	_, err = os.Stat(tmpTar)
	if err == nil {
		log.Println("[remove]", tmpTar)
		_ = os.Remove(tmpTar)
	}

}

// installVV 安装指定的小版本
func installVV(version string, vvs Versions) error {
	if vvs.Get(version) != nil {
		// 不应该执行到这个逻辑
		return fmt.Errorf("now allow, bug here")
	}
	vu, err := parserVersion(version)
	if err != nil {
		return err
	}
	mv := vvs.Get(vu.Normalized)
	if mv == nil {
		return fmt.Errorf("minor version not found")
	}
	var installVersion *Version
	for _, pv := range mv.PatchVersions {
		if pv.Raw == version || pv.Raw+".0" == version {
			installVersion = pv
			break
		}
	}
	if installVersion == nil {
		return fmt.Errorf("version not found")
	}
	return installWithVersion(installVersion)
}

func findGoBin() (string, error) {
	gb := "go" + exe()
	if p, err := exec.LookPath(gb); err == nil {
		return p, nil
	}
	if ep := findGoInSdkDir(); len(ep) > 0 {
		return ep, nil
	}

	return gb, fmt.Errorf("cannot find %q in $PATH", gb)
}

// 查找 ~/sdk/ 目录下已经安装的 go 版本
func findGoInSdkDir() string {
	sr, err := sdkRoot()
	if err != nil {
		return ""
	}
	ms, err := filepath.Glob(filepath.Join(sr, "go*"))
	if err != nil {
		return ""
	}
	sort.Slice(ms, func(i, j int) bool {
		return strings.Compare(ms[i], ms[j]) >= 0
	})
	for _, m := range ms {
		if ep, err := exec.LookPath(filepath.Join(m, "bin", "go"+exe())); err == nil {
			return ep
		}
	}
	return ""
}

func installByArchive(version string) error {
	gr, err := goroot(version)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(gr, 0755); err != nil && !os.IsExist(err) {
		return err
	}
	if err = chdir(gr); err != nil {
		return err
	}
	u := versionArchiveURL(version)
	out := u[strings.LastIndex(u, "/")+1:]
	cmd := exec.Command("wget", u, "-O", out)
	log.Println("[exec]", cmd.String())
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		return err
	}
	if err = unpackArchive(out); err != nil {
		return err
	}
	return nil
}

func unpackArchive(f string) error {
	cmd := exec.Command("tar", "-xzf", f, "--strip-components", "1")
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	log.Println("[exec]", cmd.String())
	return cmd.Run()
}

func versionArchiveURL(version string) string {
	goos := getOS()

	ext := ".tar.gz"
	if goos == "windows" {
		ext = ".zip"
	}
	arch := runtime.GOARCH
	if goos == "linux" && runtime.GOARCH == "arm" {
		arch = "armv6l"
	}
	return "https://dl.google.com/go/" + version + "." + goos + "-" + arch + ext
}
