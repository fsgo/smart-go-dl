// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/12/31

package internal

import (
	"bytes"
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

	"github.com/fsgo/cmdutils"
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

	logPrint("install", fmt.Sprintf("found %s's latest version is %s", version, last.Raw))

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
			if err = os.Symlink(filepath.Base(goBinTo), goBinLink); err != nil {
				return err
			}
		}
		logPrint("link", goBinTo, "->", goBinLink, "success")
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
		logPrint("link", mustGoRoot(last.Raw), "->", sdkDirLink, "success")
	}
	log.Printf("Success. You may now run '%s'\n", version)
	printPATHMessage(goBinTo)
	return nil
}

func printPATHMessage(goBinTo string) {
	name := filepath.Base(goBinTo)
	_, err := exec.LookPath(name)
	if err == nil {
		return
	}
	dir := filepath.Dir(goBinTo)
	log.Printf("%q not in $PATH", dir)
}

func installWithVersion(ver *Version) error {
	gb, err := findGoBin()
	if err != nil {
		// 当没有找到 go 的时候，尝试直接使用下载编译好的 go
		err = installByArchive(ver.Raw)
		if err != nil {
			return err
		}
	}

	gb, err = findGoBin()
	if err != nil {
		return err
	}

	goBinTo := ver.RawGoBinPath()

	// if _, err = exec.LookPath(goBinTo); err != nil {
	// 	err = installByArchive(ver.Raw)
	// }

	if err = chdir(ver.DlDir()); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, gb, "build", "-o", goBinTo)
	setGoEnv(cmd, gb)
	logPrint("exec", cmd.String())
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		printGoEnv(gb)
		return err
	}

	removeGoTmpTar(ver.Raw)
	log.Printf("Success. You may now run '%s'\n", filepath.Base(goBinTo))
	return err
}

func printGoEnv(gb string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, gb, "env")
	setGoEnv(cmd, gb)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

func setGoEnv(cmd *exec.Cmd, gb string) {
	goROOT := filepath.Dir(filepath.Dir(gb))
	fp := filepath.Join(goROOT, "api", "go1.1.txt")
	_, err := os.Stat(fp)
	if err != nil {
		return
	}
	cmd.Env = append(cmd.Env,
		"GOROOT="+goROOT,
		"GOCACHE="+filepath.Join(os.TempDir(), "go_build_cache"),
		"GOBIN="+GOBIN(),
	)
}

func removeGoTmpTar(version string) {
	sdk, err := sdkRoot()
	if err != nil {
		return
	}
	name := versionArchiveName(version)
	tmpTar := filepath.Join(sdk, version, name)
	_, err = os.Stat(tmpTar)
	if err == nil {
		logPrint("remove", tmpTar)
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
	if p, err := lookGoBinPath(gb); err == nil {
		return p, nil
	}
	if ep := findGoInSdkDir(); len(ep) > 0 {
		return ep, nil
	}

	return gb, fmt.Errorf("cannot find %q in $PATH", gb)
}

// lookGoBinPath 查找判断是否一个有效的 go bin
func lookGoBinPath(goFile string) (string, error) {
	gb, err := exec.LookPath(goFile)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, gb, "version").Output()
	if err != nil {
		return "", err
	}

	// 完整的 out : go version go1.16.15 darwin/amd64
	if bytes.HasPrefix(out, []byte("go version go")) {
		return gb, nil
	}
	return "", fmt.Errorf("%s is not valid go bin", goFile)
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
		if ep, err := lookGoBinPath(filepath.Join(m, "bin", "go"+exe())); err == nil {
			return ep
		}
	}
	return ""
}

// installByArchive 安装指定的 3 位版本
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
	wget := newWget()
	logPrint("download", "from", u, "to", out)

	if err = wget.Download(u, out); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	if err = unpackArchive(out); err != nil {
		return err
	}
	return nil
}

const unpackedOkay = ".unpacked-success"

func unpackArchive(f string) (err error) {
	logPrint("unpack", f)
	defer func() {
		if err != nil {
			return
		}
		_ = os.WriteFile(unpackedOkay, nil, 0644)
	}()

	if strings.HasSuffix(f, ".zip") {
		z := &cmdutils.Zip{
			StripComponents: 1,
		}
		return z.Unpack(f, "./")
	}
	tr := &cmdutils.Tar{
		StripComponents: 1,
	}
	return tr.Unpack(f, "./")
}

func versionArchiveName(version string) string {
	goos := getOS()

	ext := ".tar.gz"
	if goos == "windows" {
		ext = ".zip"
	}
	arch := runtime.GOARCH
	if goos == "linux" && runtime.GOARCH == "arm" {
		arch = "armv6l"
	}
	name := version + "." + goos + "-" + arch + ext
	return name
}

func versionArchiveURL(version string) string {
	name := versionArchiveName(version)
	return defaultConfig.getTarUrL(name)
}
