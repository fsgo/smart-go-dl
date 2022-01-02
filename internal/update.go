// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/1

package internal

import (
	"fmt"
	"log"
)

func Update(version string) error {
	if version == "all" {
		return updateAll()
	}
	return update(version)
}

func update(version string) error {
	if err := Clean(version); err != nil {
		return err
	}
	return Install(version)
}

func updateAll() error {
	versions, err := LastVersions()
	if err != nil {
		return err
	}
	var failed []string
	for _, mv := range versions {
		if mv.NormalizedVersion == "gotip" {
			log.Println("[update] skip gotip")
			continue
		}
		if mv.Installed() {
			if err = update(mv.NormalizedVersion); err != nil {
				log.Println("[update]", mv.NormalizedVersion, "failed:", err)
				failed = append(failed, mv.NormalizedVersion)
			} else {
				log.Println("[update]", mv.NormalizedVersion, "success")
			}
			log.Println()
		}
	}
	if len(failed) > 0 {
		return fmt.Errorf("update %q failed", failed)
	}
	return nil
}
