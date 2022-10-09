// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/1

package internal

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fsgo/cmdutils"
)

var gTmpDir string

// GetDataDir 获取当前应用的缓存目录, 路径为 ~/sdk/smart-go-dl
func GetDataDir() (string, error) {
	if len(gTmpDir) != 0 {
		return gTmpDir, nil
	}
	sdk, err := sdkRoot()
	if err != nil {
		return "", err
	}
	gTmpDir = filepath.Join(sdk, "smart-go-dl")
	return gTmpDir, nil
}

func chdir(dir string) error {
	err := os.Chdir(dir)
	if err == nil {
		logPrint("chdir", dir)
	} else {
		logPrint("chdir", dir, "failed:", err)
	}
	return err
}

// DataDir 获取临时目录，路径为 ~/sdk/smart-go-dl
func DataDir() string {
	if len(gTmpDir) != 0 {
		return gTmpDir
	}
	t, err := GetDataDir()
	if err != nil {
		panic(err)
	}
	return t
}

var goBinPath string

// ParserGOBIN 解析 GOBIN 环境变量
func ParserGOBIN() error {
	paths := getEnvSlice("GOBIN")
	if len(paths) == 0 || len(paths[0]) == 0 {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		goBinPath = filepath.Join(home, "go", "bin")
		return nil
	}
	goBinPath = paths[0]
	return nil
}

func getEnvSlice(key string) []string {
	sep := ":"
	if isWindows() {
		sep = ";"
	}
	return strings.Split(os.Getenv(key), sep)
}

// GOBIN 获取 GOBIN 环境变量
func GOBIN() string {
	if len(goBinPath) == 0 {
		if err := ParserGOBIN(); err != nil {
			panic(err)
		}
	}
	return goBinPath
}

func exe() string {
	if isWindows() {
		return ".exe"
	}
	return ""
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func goroot(version string) (string, error) {
	dir, err := sdkRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, version), nil
}

func mustGoRoot(version string) string {
	dir, err := goroot(version)
	if err != nil {
		panic(err)
	}
	return dir
}

func green(txt string) string {
	return colorText(txt, 32)
}

func yellow(txt string) string {
	return colorText(txt, 33)
}

func colorText(txt string, color int) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, txt)
}

func sdkRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, "sdk"), nil
}

func getOS() string {
	return runtime.GOOS
}

func copyFile(src, dst string) error {
	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()
	si, err := sf.Stat()
	if err != nil {
		return err
	}
	df, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_RDWR, si.Mode())
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	return err
}

func newWget() *cmdutils.Wget {
	gt := &cmdutils.Wget{
		LogWriter:          os.Stderr,
		Proxy:              defaultConfig.getProxy(),
		ConnectTimeout:     30 * time.Second,
		InsecureSkipVerify: defaultConfig.InsecureSkipVerify,
	}
	// 只有使用默认的下载地址的时候，才需要代理
	if !defaultConfig.isDefaultRarURLPrefix() {
		gt.Proxy = nil
	}
	return gt
}

func logPrint(key string, msgs ...any) {
	ks := fmt.Sprintf("%-10s : ", key)
	var bs strings.Builder
	bs.WriteString(ks)
	bs.WriteString(" ")
	for _, m := range msgs {
		bs.WriteString(fmt.Sprint(m))
		bs.WriteString(" ")
	}
	_ = log.Output(1, bs.String())
}
