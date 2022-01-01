// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/12/31

package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	// "golang.org/x/mod/semver"
)

type Version struct {
	// 原始的版本号，如 go1.10，go1.9rc2，go1.18beta1
	Raw string

	// 归一化的值，
	Num int

	// 归一化的，如 go1.17
	Normalized string
}

func (v *Version) String() string {
	bf, _ := json.Marshal(v)
	return string(bf)
}

func (v *Version) Installed() bool {
	sdk, err := goroot(v.Raw)
	if err != nil {
		return false
	}
	info, err := os.Stat(sdk)
	return err == nil && info.IsDir()
}

var versionReg = regexp.MustCompile(`^(go1\.\d+)((\.\d+)|(rc\d+)|(beta\d+))?$`)

func parserVersion(version string) (*Version, error) {
	matches := versionReg.FindStringSubmatch(version)
	if len(matches) == 0 {
		return nil, fmt.Errorf("not goBinPath version")
	}
	// "go1.1rc1"   -> ["go1.1rc1" "go1.1" "rc1" "" "rc1" ""]
	// "go1.12.12"  -> ["go1.12.12" "go1.12" ".12" ".12" "" ""]
	// "go1.12"     -> ["go1.12" "go1.12" "" "" "" ""]
	// "go1.1beta1" -> ["go1.1beta1" "go1.1" "beta1" "" "" "beta1"]

	vv := &Version{
		Raw:        version,
		Normalized: matches[1],
	}
	num, _ := strconv.Atoi(matches[1][4:])
	num = num * 10000

	if strings.HasPrefix(matches[2], "rc") {
		m, _ := strconv.Atoi(matches[2][2:])
		num += m
	}

	if strings.HasPrefix(matches[2], "beta") {
		m, _ := strconv.Atoi(matches[2][4:])
		num += m * 100
	}

	if strings.HasPrefix(matches[2], ".") {
		m, _ := strconv.Atoi(matches[2][1:])
		num += m * 1000
	}
	vv.Num = num
	return vv, nil
}

func LastVersions() (map[string][]*Version, error) {
	os.Chdir("/Users/baidu/tmp/dl")

	matches, err := filepath.Glob("go1.*")
	if err != nil {
		return nil, err
	}
	versions := make(map[string][]*Version)
	for _, name := range matches {
		vv, err := parserVersion(name)
		if err != nil {
			continue
		}
		versions[vv.Normalized] = append(versions[vv.Normalized], vv)
		sort.Slice(versions[vv.Normalized], func(i, j int) bool {
			a := versions[vv.Normalized][i]
			b := versions[vv.Normalized][j]
			return a.Num >= b.Num
		})
	}
	return versions, nil
}
