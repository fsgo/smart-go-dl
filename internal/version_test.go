// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/12/31

package internal

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
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
				Num:        11000,
				Normalized: "go1.1",
			},
			wantErr: false,
		},
		{
			version: "go1.10beta1",
			want: &Version{
				Raw:        "go1.10beta1",
				Num:        100001,
				Normalized: "go1.10",
			},
			wantErr: false,
		},
		{
			version: "go1.10rc1",
			want: &Version{
				Raw:        "go1.10rc1",
				Num:        100010,
				Normalized: "go1.10",
			},
			wantErr: false,
		},
		{
			version: "go1.10",
			want: &Version{
				Raw:        "go1.10",
				Num:        101000,
				Normalized: "go1.10",
			},
			wantErr: false,
		},
		{
			version: "go1.10.1",
			want: &Version{
				Raw:        "go1.10.1",
				Num:        102000,
				Normalized: "go1.10",
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
				t.Errorf("parserVersion()\n  got = %#v,\n want = %#v", got, tt.want)
			}
		})
	}
}

func Test_parserVersions(t *testing.T) {
	versions := []string{
		"go1.1",
		"go1.10", "go1.10.1", "go1.10.11",
		"go1.7beta2", "go1.7rc1",
		"go1.9beta2", "go1.9rc1", "go1.9rc2", "go1.9",
		"go1.8beta1",
		"go1.18beta2", "go1.18rc1", "go1.18", "go1.18.1",
		"gotip",
	}
	vs, err := parserVersions(versions)
	require.NoError(t, err)
	var got []string
	for _, item := range vs {
		got = append(got, item.Latest().Raw)
	}
	want := []string{"gotip", "go1.18.1", "go1.10.11", "go1.9", "go1.8beta1", "go1.7rc1", "go1.1"}
	require.Equal(t, want, got)
}
