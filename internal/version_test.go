// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/12/31

package internal

import (
	"reflect"
	"testing"
)

func Test_versionReg(t *testing.T) {
	testsOk := []string{
		"go1.1", "go1.10", "go1.10.1", "go1.10.11", "go1.10.11",
		"go1.9rc1", "go1.9rc2",
		"go1.8beta1",
		"go1.18beta2",
	}
	for _, tt := range testsOk {
		t.Run(tt, func(t *testing.T) {
			if !versionReg.MatchString(tt) {
				t.Fatalf("not match")
			}
		})
	}

	testsNot := []string{
		"ggo1.1", "1.10", "go1.10.1v2", "go1.10.11x",
	}
	for _, tt := range testsNot {
		t.Run(tt, func(t *testing.T) {
			if versionReg.MatchString(tt) {
				t.Fatalf("expect not match")
			}
		})
	}
}

func Test_parserVersion(t *testing.T) {
	tests := []struct {
		version string
		want    *Version
		wantErr bool
	}{
		{
			version: "go1.1",
			want: &Version{
				Raw:        "go1.1",
				Num:        10000,
				Normalized: "go1.1.x",
			},
			wantErr: false,
		},
		{
			version: "go1.10",
			want: &Version{
				Raw:        "go1.10",
				Num:        100000,
				Normalized: "go1.10.x",
			},
			wantErr: false,
		},
		{
			version: "go1.10.1",
			want: &Version{
				Raw:        "go1.10.1",
				Num:        101000,
				Normalized: "go1.10.x",
			},
			wantErr: false,
		},
		{
			version: "go1.10rc1",
			want: &Version{
				Raw:        "go1.10rc1",
				Num:        100001,
				Normalized: "go1.10.x",
			},
			wantErr: false,
		},
		{
			version: "go1.10beta1",
			want: &Version{
				Raw:        "go1.10beta1",
				Num:        100100,
				Normalized: "go1.10.x",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			got, err := parserVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("parserVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parserVersion() got = %v, want %v", got, tt.want)
			}
		})
	}
}
