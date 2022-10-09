// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/12/31

package internal

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const dlStatsFile = "download.status"
const golangDLDir = "golang_dl"

// Download 下载 golang/dl.git
func Download() error {
	dlStatsPath := filepath.Join(DataDir(), dlStatsFile)
	writeStats := func() {
		_ = os.WriteFile(dlStatsPath, []byte(time.Now().String()), 0655)
	}
	info, _ := os.Stat(dlStatsPath)

	dlDir := filepath.Join(DataDir(), golangDLDir)
	_, err := os.Stat(dlDir)
	if err == nil {
		// 短期内更新过的
		// 这个库本来更新也非常少
		if info != nil && time.Since(info.ModTime()) < 1*time.Hour {
			return nil
		}

		if err = chdir(dlDir); err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		cmdPull := exec.CommandContext(ctx, "git", "pull", "-v")
		logPrint("exec", cmdPull.String())
		cmdPull.Stderr = os.Stderr
		cmdPull.Stdout = os.Stdout
		if err = cmdPull.Run(); err != nil {
			logPrint("skipped", "git pull failed, ", err)
			// 忽略异常,可能由于 Q 的存在，更新最新版本不是很稳定
		}
		writeStats()
		return nil
	}

	repo := defaultRepo
	args := []string{"clone", repo, golangDLDir}
	cmdClone := exec.Command("git", args...)
	logPrint("exec", cmdClone.String())
	cmdClone.Stdin = os.Stdin
	cmdClone.Stderr = os.Stderr
	cmdClone.Stdout = os.Stdout
	if err = cmdClone.Run(); err != nil {
		// 若直接下载失败了，则使用内置的，将其解压到对应目录下去
		logPrint("fallback", "extract", defaultRepo, "by embed datas")
		err2 := extractGolangDLTar(dlDir)
		if err2 == nil {
			return nil
		}
		return err
	}
	writeStats()
	return nil
}

const defaultRepo = "https://github.com/golang/dl.git"
