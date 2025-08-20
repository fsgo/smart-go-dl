//  Copyright(C) 2025 github.com/hidu  All Rights Reserved.
//  Author: hidu <duv123+git@gmail.com>
//  Date: 2025-08-20

package internal

import (
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsgo/cmdutil/gosdk"
)

func logDir() string {
	return filepath.Join(DataDir(), "log")
}

func TrySetLogFile(name string) func() {
	logFile, err1 := newLogFile(name)
	if err1 != nil {
		log.Println("create log file:", err1)
		return func() {}
	}
	var lf io.Writer
	if name == "go" {
		lf = logFile
		log.SetFlags(log.LstdFlags)
	} else {
		lf = io.MultiWriter(os.Stderr, logFile)
	}

	log.SetOutput(lf)
	gosdk.SetLogger(log.Default())
	return func() {
		_ = logFile.Close()
	}
}

func newLogFile(name string) (*os.File, error) {
	dir := logDir()
	_ = os.MkdirAll(dir, 0777)
	fileName := filepath.Join(dir, name+".log."+time.Now().Format("20060102"))
	log.Println("logFile    : ", fileName)
	go cleanLogFiles()
	return os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
}

func cleanLogFiles() {
	expireTime := time.Now().Add(-72 * time.Hour)
	_ = filepath.WalkDir(logDir(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := os.Stat(path)
		if err != nil {
			return nil // ignore error
		}
		if info.ModTime().Before(expireTime) {
			_ = os.Remove(path)
		}
		return nil
	})
}
