// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/1

package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

var gTmpDir string

func InitTmpDir() (string, error) {
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

func TmpDir() string {
	if len(gTmpDir) != 0 {
		return gTmpDir
	}
	t, err := InitTmpDir()
	if err != nil {
		panic(err)
	}
	return t
}

var goBinPath string

func ParserGOBIN() error {
	paths := getEnvSlice("GOBIN")
	if len(paths) == 0 || len(paths[0]) == 0 {
		return fmt.Errorf("GOBIN has not setted")
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
	home, err := homedir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}
	return filepath.Join(home, "sdk"), nil
}

func getOS() string {
	return runtime.GOOS
}

func homedir() (string, error) {
	// This could be replaced with os.UserHomeDir, but it was introduced too
	// recently, and we want this to work with go as packaged by Linux
	// distributions. Note that user.Current is not enough as it does not
	// prioritize $HOME. See also Issue 26463.
	switch getOS() {
	case "plan9":
		return "", fmt.Errorf("%q not yet supported", runtime.GOOS)
	case "windows":
		if dir := os.Getenv("USERPROFILE"); dir != "" {
			return dir, nil
		}
		return "", errors.New("can't find user home directory; %USERPROFILE% is empty")
	default:
		if dir := os.Getenv("HOME"); dir != "" {
			return dir, nil
		}
		if u, err := user.Current(); err == nil && u.HomeDir != "" {
			return u.HomeDir, nil
		}
		return "", errors.New("can't find user home directory; $HOME is empty")
	}
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
