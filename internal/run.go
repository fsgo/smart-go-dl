package internal

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fsgo/cmdutil/gosdk"
)

var goCMDReg = regexp.MustCompile(`^go1\.\d+`)

// TryRunGo 尝试运行 go 命令，如 go env
func TryRunGo(name string) {
	name = strings.TrimRight(filepath.Base(name), exe())
	if name == "go" || name == "go.latest" {
		runLatest()
		return
	}

	if goCMDReg.MatchString(name) {
		run(name)
	}
}

func runLatest() {
	sd := &gosdk.SDK{
		ExtDirs: []string{SDKRootDir()},
	}
	goBin := sd.Latest()
	if goBin == "" {
		log.Fatalln("not found go")
	}
	root := filepath.Dir(filepath.Dir(goBin))
	gosdk.RunGo(root)
}

func run(version string) {
	log.SetFlags(0)

	sd := &gosdk.SDK{
		ExtDirs: []string{SDKRootDir()},
	}

	goBin := sd.Find(version)
	if goBin == "" {
		log.Fatalln("not found", version)
	}

	if len(os.Args) == 2 && os.Args[1] == "download" {
		if err := installByArchive(version); err != nil {
			log.Fatalf("%s: install failed: %v", version, err)
		}
		os.Exit(0)
	}

	root := filepath.Dir(filepath.Dir(goBin))

	if _, err := os.Stat(filepath.Join(root, unpackedOkay)); err != nil {
		log.Fatalf("%s: not downloaded. Run '%s download' to install to %v", version, version, root)
	}

	gosdk.RunGo(root)
}
