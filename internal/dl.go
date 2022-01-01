// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/12/31

package internal

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
)

func Prepare() error {
	tmpDir, err := TmpDir()
	if err != nil {
		return err
	}

	if err = ParserGOBIN(); err != nil {
		return err
	}

	log.Println("TmpDir:", tmpDir)

	if err = os.MkdirAll(tmpDir, 0777); err != nil && !os.IsExist(err) {
		return err
	}

	if err = os.Chdir(tmpDir); err != nil {
		return err
	}
	return Download()
}

func TmpDir() (string, error) {
	sdk, err := sdkRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(sdk, "smart-go-dl"), nil
}

func Download() error {
	_, err := os.Stat("dl")
	if err == nil {
		pwd, err1 := os.Getwd()
		if err1 != nil {
			return err
		}
		if err = os.Chdir("dl"); err != nil {
			return err
		}
		cmdPull := exec.Command("git", "pull")
		log.Println("[exec]", cmdPull.String())
		cmdPull.Stderr = os.Stderr
		cmdPull.Stdout = os.Stdout
		if err = cmdPull.Run(); err != nil {
			return err
		}
		return os.Chdir(pwd)
	}
	args := []string{
		"clone",
		"git@github.com:golang/dl.git",
	}
	cmdClone := exec.Command("git", args...)
	log.Println("[exec]", cmdClone.String())
	cmdClone.Stderr = os.Stderr
	cmdClone.Stdout = os.Stdout
	return cmdClone.Run()
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
