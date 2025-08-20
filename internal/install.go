// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/12/31

package internal

import (
	"context"
	"debug/buildinfo"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/fsgo/cmdutil"
)

// Install 安装 go1.x 的最新版本
//
// version: 版本号，如 1.21
func Install(ctx context.Context, version string) error {
	versions, err := LastVersions(ctx)
	if err != nil {
		return err
	}
	defer installGoLatestBin(ctx)

	mv := versions.Get(version)
	if mv == nil {
		logPrint("installVV", version)

		// 用于支持安装 3 位版本，如  go1.16.0、go1.16.3
		err = installVV(version, versions)
		if err != nil {
			return fmt.Errorf("install %q failed: %w", version, err)
		}
		return nil
	}

	last := mv.Latest()

	logPrint("install", fmt.Sprintf("found %s's latest version is %s", version, last.Raw))

	goBinTo := last.RawGoBinPath()

	if err = installWithVersion(last); err != nil {
		return err
	}

	goBinLink := last.NormalizedGoBinPath()
	logPrint("trace", "goBinLink=", goBinLink, "goBinTo=", goBinTo)
	if goBinLink == goBinTo {
		return nil
	}

	// create link for go bin
	// go1.16.6 -> go1.16
	if err = createLink(goBinTo, goBinLink); err != nil {
		return err
	}

	log.Printf("Success. You may now run '%s'\n", version)
	printPATHMessage(goBinTo)
	return nil
}

func createLink(from string, to string) error {
	from = filepath.Clean(from)
	to = filepath.Clean(to)
	if from == to {
		return nil
	}
	if err := os.Remove(to); err != nil && !os.IsNotExist(err) {
		return err
	}
	if isWindows() {
		if err := copyFile(from, to); err != nil {
			return err
		}
	} else {
		if err := os.Symlink(filepath.Base(from), to); err != nil {
			return err
		}
	}
	logPrint("link", from, "->", to, "success")
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
	logPrint("trace", "installWithVersion", ver.String())
	_, err := findGoBin()
	if err != nil {
		// 当没有找到 go 的时候，尝试直接使用下载编译好的 go
		err = installByArchive(ver.Raw)
		if err != nil {
			return err
		}
	}

	goBinTo := ver.RawGoBinPath()

	// smart-go-dl 可以将自己重命名为 go，并支持运行的时候使用 go download 下载 sdk 文件
	selfPath := os.Getenv("_")
	if selfPath == "" {
		selfPath = os.Args[0]
	}

	if err = createLink(selfPath, goBinTo); err != nil {
		logPrint("createLink", selfPath, "->", goBinTo, ", err=", err)
	}

	// if err = copyFile(selfPath, goBinTo); err != nil {
	//	logPrint("copyFile", selfPath, "->", goBinTo, "err=", err)
	//	return err
	// }

	out, err1 := lookGoBinPath(goBinTo)
	logPrint("trace", "check", goBinTo, out, err1)
	if err1 != nil || strings.Contains(out, "not downloaded") {
		if err2 := installByArchive(ver.Raw); err2 != nil {
			logPrint("download", err2.Error())
			return err2
		}
	}

	removeGoTmpTar(ver.Raw)
	log.Printf("Success. You may now run '%s'\n", filepath.Base(goBinTo))
	return err
}

func printGoEnv(gb string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, gb, "env")
	setGoEnv(cmd, gb)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	logPrint("trace", cmd.String(), err)
}

func setGoEnv(cmd *exec.Cmd, gb string) {
	goROOT := filepath.Dir(filepath.Dir(gb))
	fp := filepath.Join(goROOT, "api", "go1.1.txt")
	_, err := os.Stat(fp)
	if err != nil {
		logPrint("trace", "setGoEnv", err)
	}
	cmd.Env = append(os.Environ(),
		"GOROOT="+goROOT,
		"GOCACHE="+filepath.Join(os.TempDir(), "go_build_cache"),
		"GOPATH="+filepath.Dir(GOBIN()),
		"GOBIN="+GOBIN(),
	)
}

func removeGoTmpTar(version string) {
	sdkDir := defaultConfig.getSDKDir()
	name := versionArchiveName(version)
	tmpTar := filepath.Join(sdkDir, version, name)
	_, err := os.Stat(tmpTar)
	if err == nil {
		logPrint("remove", tmpTar)
		_ = os.Remove(tmpTar)
	}
}

// installVV 安装指定的小版本
func installVV(version string, vvs Versions) error {
	if vvs.Get(version) != nil {
		// 不应该执行到这个逻辑
		return errors.New("now allow, bug here")
	}
	vu, err := parserVersion(version)
	if err != nil {
		return err
	}
	mv := vvs.Get(vu.Normalized)
	if mv == nil {
		return errors.New("minor version not found")
	}
	var installVersion *Version
	for _, pv := range mv.PatchVersions {
		if pv.Raw == version || pv.Raw+".0" == version {
			installVersion = pv
			break
		}
	}
	if installVersion == nil {
		return errors.New("version not found")
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

	bi, err := buildinfo.ReadFile(gb)
	if err != nil {
		return "", err
	}

	if bi.Path != "cmd/go" {
		return "", fmt.Errorf("not cmd/go, got %q", bi.Path)
	}

	// go1.16.15
	if strings.HasPrefix(bi.GoVersion, "go") {
		return gb, nil
	}
	return bi.GoVersion, fmt.Errorf("%s is not valid go bin", goFile)
}

// 查找 ~/sdk/ 目录下已经安装的 go 版本
func findGoInSdkDir() string {
	sdkDir := defaultConfig.getSDKDir()
	ms, err := filepath.Glob(filepath.Join(sdkDir, "go*"))
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
	urls := versionArchiveURLs(version)
	logPrint("trace", "urls", urls)

	for _, u := range urls {
		out := u[strings.LastIndex(u, "/")+1:]
		if err = wget(u, out); err != nil {
			continue
		}
		if err = unpackArchive(out); err == nil {
			break
		}
	}
	return err
}

const unpackedOkay = ".unpacked-success"

func unpackArchive(f string) (err error) {
	info, err := os.Stat(f)
	if err != nil {
		logPrint("unpack", "error,", err)
		return err
	}
	logPrint("unpack", f, "size=", info.Size())
	defer func() {
		logPrint("unpack", "done,", err)
		if err != nil {
			return
		}
		_ = os.WriteFile(unpackedOkay, nil, 0644)
	}()

	if strings.HasSuffix(f, ".zip") {
		z := &cmdutil.Zip{
			StripComponents: 1,
		}
		return z.Unpack(f, "./")
	}
	tr := &cmdutil.Tar{
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

func versionArchiveURLs(version string) []string {
	name := versionArchiveName(version)
	urls := defaultConfig.getTarURLs(name)
	return sortURLs(urls)
}
