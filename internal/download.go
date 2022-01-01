// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/12/31

package internal

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const dlStatsFile = "dl.status"

func Download() error {
	dlStatsPath := filepath.Join(TmpDir(), dlStatsFile)
	info, _ := os.Stat(dlStatsPath)

	dlDir := filepath.Join(TmpDir(), "dl")
	_, err := os.Stat(dlDir)
	if err == nil {
		// 5 分钟内更新过
		if info != nil && time.Since(info.ModTime()) < 5*time.Minute {
			return nil
		}

		if err = os.Chdir(dlDir); err != nil {
			return err
		}

		cmdPull := exec.Command("git", "pull")
		log.Println("[exec]", cmdPull.String())
		cmdPull.Stderr = os.Stderr
		cmdPull.Stdout = os.Stdout
		if err = cmdPull.Run(); err != nil {
			return err
		}
		_ = ioutil.WriteFile(dlStatsPath, []byte(time.Now().String()), 0655)
		return nil
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
