// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/3

package internal

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func Remove(version string) error {
	return remove(version)
}

func remove(version string) error {
	v, err := parserVersion(version)
	if err != nil {
		return err
	}

	if strings.HasSuffix(version, ".0") {
		version = version[0 : len(version)-2]
	}

	sdkDir, err := goroot(version)
	if err != nil {
		return err
	}

	if _, err = os.Stat(sdkDir); err != nil && os.IsNotExist(err) {
		return fmt.Errorf("version %q not installed", version)
	}

	goBin := v.RawGoBinPath()
	log.Println("[clean] remove ", goBin)
	if err = os.Remove(goBin); err != nil && !os.IsNotExist(err) {
		return err
	}

	log.Println("[clean] remove ", sdkDir)
	if err = os.RemoveAll(sdkDir); err != nil && !os.IsNotExist(err) {
		return err
	}

	vs, err := LastVersions()
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
