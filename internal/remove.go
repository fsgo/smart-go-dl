// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/3

package internal

import (
	"context"
	"fmt"
	"os"
	"strings"
)

// Remove 删除指定的版本
func Remove(ctx context.Context, version string) error {
	return remove(ctx, version)
}

func remove(ctx context.Context, version string) error {
	defer installGoLatestBin(ctx)

	v, err := parserVersion(version)
	if err != nil {
		return err
	}

	version = strings.TrimSuffix(version, ".0")

	sdkDir, err := goroot(version)
	if err != nil {
		return err
	}

	if _, err = os.Stat(sdkDir); err != nil && os.IsNotExist(err) {
		return fmt.Errorf("version %q not installed", version)
	}

	goBin := v.RawGoBinPath()

	logPrint("remove", goBin)
	if err = os.Remove(goBin); err != nil && !os.IsNotExist(err) {
		return err
	}

	logPrint("remove", sdkDir)
	if err = os.RemoveAll(sdkDir); err != nil && !os.IsNotExist(err) {
		return err
	}

	vs, err := LastVersions(ctx)
	if err != nil {
		return err
	}
	mv := vs.Get(v.Normalized)
	if mv != nil {
		last := mv.Latest()
		if last.Raw == v.Raw {
			link := v.NormalizedGoBinPath()
			if err = os.Remove(link); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}

	return nil
}
