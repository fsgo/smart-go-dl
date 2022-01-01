// Copyright(C) 2022 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2022/1/1

package internal

func Update(version string) error {
	if err := Clean(version); err != nil {
		return err
	}
	return Install(version)
}
