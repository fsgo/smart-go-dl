// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/12/31

package internal

import (
	"bytes"
	"context"
	"errors"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
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
		if info != nil && time.Since(info.ModTime()) < time.Minute {
			return nil
		}
		if err = chdir(dlDir); err != nil {
			return err
		}

		err = gitPull()
		if err == nil {
			writeStats()
		}
		return nil
	}

	repo := defaultRepo
	args := []string{"clone", repo, golangDLDir}
	cmdClone := exec.Command("git", args...)
	logPrint("exec", cmdClone.String())
	setGitCmdEnv(cmdClone)
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

func setGitCmdEnv(cmd *exec.Cmd) {
	if defaultConfig.InsecureSkipVerify {
		cmd.Env = append(os.Environ(), "GIT_SSL_NO_VERIFY=false")
	}
}

var useGoGit = len(os.Getenv("Smart_Go_Dl_GoGit")) != 0

func gitPull() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if !useGoGit {
		cmdPull := exec.CommandContext(ctx, "git", "pull", "-v")
		logPrint("exec", cmdPull.String())
		setGitCmdEnv(cmdPull)
		cmdPull.Stderr = os.Stderr
		cmdPull.Stdout = os.Stdout
		cmdPull.Stdin = os.Stdin
		err := cmdPull.Run()
		if err == nil {
			return nil
		}

		logPrint("git pull failed, ", err)
	}

	gr, err := git.PlainOpen("./")
	if err != nil {
		logPrint("try open with pure Go git failed,", err)
		return err
	}
	w, err := gr.Worktree()
	if err != nil {
		logPrint("pure Go git Worktree:", err)
		return err
	}
	err = w.PullContext(ctx, &git.PullOptions{})
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			logPrint("pure GoGit:", "git pull ", err)
			return nil
		}
		logPrint("pure GoGit:", "git pull ", err)
	}
	return err
}

const defaultRepo = "https://github.com/golang/dl.git"

func wget(url string, to string) error {
	logPrint("download", "from", url, "to", to)
	w1 := newWget()
	err1 := w1.Download(url, to)
	if err1 == nil {
		return nil
	}

	logPrint("go-wget", "failed:", err1, ", will retry")

	var args []string
	if defaultConfig.InsecureSkipVerify {
		args = append(args, "--no-check-certificate")
	}
	args = append(args, "--connect-timeout=5", "--tries=1", "-O", to)
	args = append(args, url)
	cmd1 := exec.Command("wget", args...)
	logPrint("exec", cmd1.String())
	cmd1.Stderr = os.Stderr
	cmd1.Stdin = os.Stdin
	cmd1.Stdout = os.Stdout
	return cmd1.Run()
}

const enableSort = false

// sortURLs 对 url 地址排序，速度快的排前面
func sortURLs(urls []string) []string {
	if !enableSort {
		return urls
	}

	type info struct {
		err  error
		url  string
		cost time.Duration
	}

	w1 := newWget()
	w1.LogWriter = nil
	w1.ConnectTimeout = 3 * time.Second
	w1.Timeout = 5 * time.Second

	var items []info
	for _, api := range urls {
		func() {
			start := time.Now()
			u, err := url.Parse(api)
			if err != nil {
				items = append(items, info{url: api, err: err, cost: time.Hour})
				return
			}
			u.Path = "/404_page"
			bf := &bytes.Buffer{}
			err2 := w1.DownloadToWriter(u.String(), bf)
			item := info{
				url:  api,
				err:  err2,
				cost: time.Since(start),
			}

			if err2 != nil && strings.Contains(err2.Error(), "i/o timeout") {
				item.cost = time.Hour
			}

			// 这种认为是成功了
			if strings.Contains(bf.String(), "text/html") {
				item.err = nil
			}
			items = append(items, item)
		}()
	}

	sort.Slice(items, func(i, j int) bool {
		// 失败的排后面
		if items[i].err != nil && items[j].err == nil {
			return false
		}
		// 耗时少的排前面
		return items[i].cost < items[j].cost
	})
	result := make([]string, 0, len(urls))
	for _, item := range items {
		result = append(result, item.url)
	}
	return result
}
