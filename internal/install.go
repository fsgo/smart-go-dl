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
	vinfos := versions[version]
	if len(vinfos) == 0 {
		return fmt.Errorf("version %q not found", version)
	}
	last := vinfos[0]

	log.Println("[install]", "found last", version, "version is", last.Raw)

	if err := os.Chdir(last.Raw); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	to := filepath.Join(GOBIN(), last.Raw)
	cmd := exec.CommandContext(ctx, "go", "build", "-o", to)
	log.Println("[exec]", cmd.String())
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		return err
	}

	downloadCmd := exec.Command(to, "download")
	log.Println("[exec]", downloadCmd.String())
	downloadCmd.Stderr = os.Stderr
	downloadCmd.Stdout = os.Stdout
	if err = downloadCmd.Run(); err != nil {
		return err
	}

	link := filepath.Join(GOBIN(), last.Normalized)
	if link == to {
		return nil
	}
	if err = os.Remove(link); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err = os.Symlink(to, link); err != nil {
		return err
	}
	log.Println("[link]", to, "->", link, "success")
	log.Printf("Success. You may now run '%s'\n", version)
	return nil
}

var goBinPath string

func ParserGOBIN() error {
	paths := strings.Split(os.Getenv("GOBIN"), ":")
	if len(paths) == 0 {
		return fmt.Errorf("GOBIN has not setted")
	}
	goBinPath = paths[len(paths)-1]
	return nil
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
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

func Update(version string) error {
	if err := Install(version); err != nil {
		return err
	}
	return Clean(version)
}

func List() error {
	log.Println("list call")
	versions, err := LastVersions()
	if err != nil {
		return err
	}
	vlist := make([]string, 0, len(versions))
	for v := range versions {
		vlist = append(vlist, v)
	}
	sort.Slice(vlist, func(i, j int) bool {
		a := versions[vlist[i]][0]
		b := versions[vlist[j]][0]
		return a.Num > b.Num
	})

	format := "%-20s %-20s %-20s\n"
	formatColor := "%-31s %-20s %-20s\n"
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf(format, "version", "latest", "installed")
	fmt.Println(strings.Repeat("-", 80))

	for _, v := range vlist {
		infos := versions[v]
		latest := infos[0]
		cell1 := v
		localFormat := format
		installed := strings.Join(installedVersions(infos), " ")
		if latest.Installed() {
			cell1 = green(v)
			localFormat = formatColor
		} else if len(installed) > 0 {
			cell1 = yellow(v)
			localFormat = formatColor
		}
		fmt.Printf(localFormat, cell1, latest.Raw, installed)
	}
	return nil
}

func installedVersions(vs []*Version) []string {
	var result []string
	for _, v := range vs {
		if v.Installed() {
			result = append(result, fmt.Sprintf("%-12s", v.Raw))
		}
	}
	return result
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
