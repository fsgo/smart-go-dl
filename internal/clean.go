// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/1

package internal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// Clean 将go1.x的老版本删除掉
func Clean(ctx context.Context, version string) error {
	versions, err := LastVersions(ctx)
	if err != nil {
		return err
	}

	mv := versions.Get(version)
	if mv == nil {
		return fmt.Errorf("version %q not found", version)
	}

	if len(mv.PatchVersions) < 2 {
		logPrint("clean", "no old versions need to be clean")
		return nil
	}

	for i := 1; i < len(mv.PatchVersions); i++ {
		cur := mv.PatchVersions[i]
		if err = cleanVersion(cur); err != nil {
			logPrint("clean", cur.Raw, "failed:", err)
		}
	}
	return nil
}

func cleanVersion(v *Version) error {
	sdkDir, err := goroot(v.Raw)
	if err != nil {
		return err
	}

	if _, err = os.Stat(sdkDir); err != nil && os.IsNotExist(err) {
		return nil
	}

	ignoreFile := filepath.Join(sdkDir, lockedName)
	if _, err = os.Stat(ignoreFile); err == nil {
		logPrint("clean", v.Raw, "locked")
		return nil
	}

	goBin := v.RawGoBinPath()
	logPrint("clean", "remove ", goBin)
	if err = os.Remove(goBin); err != nil && !os.IsNotExist(err) {
		return err
	}
	logPrint("clean", "remove ", sdkDir)
	if err = os.RemoveAll(sdkDir); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
