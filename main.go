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
smart-go-dl subCommand [options]

SubCommands:
    install {go1.x} :
        install the latest go1.x, 'x' must be a number, x >= 5
          eg: "install go1.18", then you can run "go1.18"
        install the specified version:
          eg: install go1.17.0 | go1.17.3 | gotip
    
    clean {go1.x} :
        clean up expired go versions.
        lower than the latest version will be removed.
        it will remove $GOBIN/{go1.x.y} and $HOME/sdk/{go1.x.y}
        eg: "clean go1.15"
    
    lock {go1.x.y} :
        add lock file. eg: "lock go1.16.1"
    
    unlock {go1.x.y} :
        remove lock file. eg: "unlock go1.18beta1"
    
    update {go1.x} / all :
        alias of  "clean {go1.x}" && "install {go1.x}"
        "all": update all installed go versions, eg: "update all" or "update"

    remove {go1.x.y} :
        remove patch version like 'go1.12.1'
    
    list :
        list all go versions that can be installed.

Self-Update :
          go install github.com/fsgo/smart-go-dl@main

Site    : https://github.com/fsgo/smart-go-dl
Version : 0.1.5
Date    : 2022-03-16
`

func init() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(out, strings.TrimSpace(helpMessage)+"\n")
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
	case "lock":
		err = internal.Lock(args.get(2), "add")
	case "unlock":
		err = internal.Lock(args.get(2), "remove")
	case "list":
		err = internal.List()
	case "remove", "uninstall":
		err = internal.Remove(args.get(2))
	default:
		err = fmt.Errorf("not support")
	}

	if err != nil {
		log.Fatalf("error: %s failed, %v\n", args[1], err)
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
