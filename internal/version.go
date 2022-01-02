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

var regNumSuf = regexp.MustCompile(`\.\d+$`)

func (v *Version) RawGoBinPath() string {
	return filepath.Join(GOBIN(), v.RawFormatted()) + exe()
}

func (v *Version) RawFormatted() string {
	name := v.Raw
	// 如 原始版本的 go1.16，实际应该是 go1.16.0
	// 用正则过滤掉非数字版本，如 gotip
	if regNumSuf.MatchString(v.Raw) && v.Raw == v.Normalized {
		name += ".0"
	}
	return name
}

func (v *Version) NormalizedGoBinPath() string {
	return filepath.Join(GOBIN(), v.Normalized) + exe()
}

func (v *Version) Installed() bool {
	sdk, err := goroot(v.Raw)
	if err != nil {
		return false
	}
	info, err := os.Stat(sdk)
	if err != nil || (err == nil && !info.IsDir()) {
		return false
	}
	_, err = os.Readlink(sdk)
	if err != nil {
		return true
	}
	return false
}

func (v *Version) DlDir() string {
	return filepath.Join(TmpDir(), "dl", v.Raw)
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

type MinorVersion struct {
	NormalizedVersion string
	PatchVersions     []*Version
}

func (mv *MinorVersion) Latest() *Version {
	return mv.PatchVersions[0]
}

func (mv *MinorVersion) Installed() bool {
	for _, pv := range mv.PatchVersions {
		if pv.Installed() {
			return true
		}
	}
	return false
}

type Versions []*MinorVersion

func (vs Versions) Get(version string) *MinorVersion {
	for _, mv := range vs {
		if mv.NormalizedVersion == version {
			return mv
		}
	}
	return nil
}

func LastVersions() (Versions, error) {
	pt := filepath.Join(TmpDir(), "dl", "go1.*")
	matches, err := filepath.Glob(pt)
	if err != nil {
		return nil, err
	}
	versions := make(map[string][]*Version)
	for _, name := range matches {
		vv, err := parserVersion(filepath.Base(name))
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

	versions["gotip"] = []*Version{
		{
			Raw:        "gotip",
			Normalized: "gotip",
			Num:        2000000,
		},
	}

	var result Versions
	for v, list := range versions {
		sort.Slice(list, func(i, j int) bool {
			a := list[i]
			b := list[j]
			return a.Num > b.Num
		})
		mv := &MinorVersion{
			NormalizedVersion: v,
			PatchVersions:     list,
		}
		result = append(result, mv)
	}

	sort.Slice(result, func(i, j int) bool {
		a := result[i]
		b := result[j]
		return a.Latest().Num > b.Latest().Num
	})

	return result, nil
}
