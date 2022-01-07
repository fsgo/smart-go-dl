// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/1

package internal

import (
	"fmt"
	"strings"
)

// List 列出已安装和可安装的 go 版本
func List() error {
	versions, err := LastVersions()
	if err != nil {
		return err
	}

	format := "%-20s %-20s %-20s\n"
	formatColor := "%-31s %-20s %-20s\n"
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf(format, "version", "latest", "installed")
	fmt.Println(strings.Repeat("-", 80))

	for _, mv := range versions {
		latest := mv.PatchVersions[0]
		cell1 := mv.NormalizedVersion
		localFormat := format
		installed := strings.Join(installedVersions(mv.PatchVersions), " ")
		if !isWindows() {
			if latest.Installed() {
				cell1 = green(cell1)
				localFormat = formatColor
			} else if len(installed) > 0 {
				cell1 = yellow(cell1)
				localFormat = formatColor
			}
		}
		fmt.Printf(localFormat, cell1, latest.Raw, installed)
	}
	return nil
}

func installedVersions(vs []*Version) []string {
	var result []string
	for _, v := range vs {
		if v.Installed() {
			name := v.RawFormatted()
			if isLocked(v.Raw) {
				name += "(L)"
			}
			result = append(result, fmt.Sprintf("%-12s", name))
		}
	}
	return result
}
