// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/1

package internal

import (
	"fmt"
	"sort"
	"strings"
)

func List() error {
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
