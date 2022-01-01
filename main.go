// Copyright(C) 2021 github.com/fsgo  All Rights Reserved.
// Author: fsgo
// Date: 2021/12/31

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fsgo/smart-go-dl/internal"
)

var helpMessage = `
install {go1.x} :
    install the latest go1.x, 'x' must be a number, x >= 5
    eg: "install go1.18", then you can run "go1.18"

clean {go1.x} :
    clean up expired go versions.
    lower than the latest version will be removed.
    it will remove $GOBIN/{go1.x.y} and $HOME/sdk/{go1.x.y}
    eg: "clean go1.15"

    add $HOME/sdk/{go1.x.y}/smart-go-dl.ignore_clean to ignore clean

update {go1.x} :
    alias of "install {go1.x}" && "clean {go1.x}"

list :
    list all go versions that can be installed.

update self :
    go install github.com/fsgo/smart-go-dl@main

Site    : https://github.com/fsgo/smart-go-dl
Version : 0.1.0
Date    : 2022-01-01
`

func init() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(out, strings.TrimSpace(helpMessage))
	}
}

func main() {
	flag.Parse()

	args := stringSlice(os.Args)
	if len(args) < 2 || args.get(1) == "help" {
		flag.Usage()
		return
	}

	if err := internal.Prepare(); err != nil {
		log.Fatalln(err)
	}

	var err error
	switch args[1] {
	case "install":
		err = internal.Install(args.get(2))
	case "clean":
		err = internal.Clean(args.get(2))
	case "update":
		err = internal.Update(args.get(2))
	case "list":
		err = internal.List()
	default:
		err = fmt.Errorf("not support")
	}

	if err != nil {
		log.Fatalf("%s error: %v\n", args[1], err)
	} else {
		log.Printf("%s success", args[1])
	}
}

func init() {
	log.SetFlags(log.Lmsgprefix)
	log.SetPrefix("[smart-go-dl] ")
}

type stringSlice []string

func (s stringSlice) get(index int) string {
	if index >= len(s) {
		return ""
	}
	return s[index]
}
