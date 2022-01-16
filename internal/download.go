// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/12/31

package internal

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const dlStatsFile = "download.status"
const golangDLDir = "golang_dl"

// Download 下载 golang/dl.git
func Download() error {
	dlStatsPath := filepath.Join(TmpDir(), dlStatsFile)
	writeStats := func() {
		_ = ioutil.WriteFile(dlStatsPath, []byte(time.Now().String()), 0655)
	}
	info, _ := os.Stat(dlStatsPath)

	dlDir := filepath.Join(TmpDir(), golangDLDir)
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
		log.Println("[exec]", cmdPull.String())
		cmdPull.Stderr = os.Stderr
		cmdPull.Stdout = os.Stdout
		if err = cmdPull.Run(); err != nil {
			log.Println("skipped error: git pull failed, ", err)
			// 忽略异常,可能由于 Q 的存在，更新最新版本不是很稳定
		}
		writeStats()
		return nil
	}

	repo := getGolangDlRepo()
	args := []string{"clone", repo, golangDLDir}
	cmdClone := exec.Command("git", args...)
	log.Println("[exec]", cmdClone.String())
	cmdClone.Stdin = os.Stdin
	cmdClone.Stderr = os.Stderr
	cmdClone.Stdout = os.Stdout
	if err = cmdClone.Run(); err != nil {
		// 若直接下载失败了，则使用内置的，将其解压到对应目录下去
		log.Println("[failback] extract", defaultRepo, "by embed datas")
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

func getGolangDlRepo() string {
	fp := filepath.Join(TmpDir(), "golang_dl.txt")
	txt, err := ioutil.ReadFile(fp)
	if err != nil {
		return defaultRepo
	}
	txt = bytes.TrimSpace(txt)
	if len(txt) == 0 {
		return defaultRepo
	}
	return string(txt)
}
