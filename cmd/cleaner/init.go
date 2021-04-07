package main

import (
	"flag"
	"fmt"
	"github.com/ClessLi/fastdfs_cleaner/internal/pkg/fastdfs_cleaner"
	"os"
	"path/filepath"
)

var (
	configPath = flag.String("f", "", "the cleaner `config`uration file path.")
	help       = flag.Bool("h", false, "this `help`")
	logF       *os.File
)

func init() {
	flag.Parse()
	if *configPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	fastdfs_cleaner.SetConfigFilePath(*configPath)
	initLogout()
}

func initLogout() {
	config := fastdfs_cleaner.GetSingletonConfigInstance()
	logDir := config.LogDir
	if !filepath.IsAbs(logDir) {
		var err error
		logDir = filepath.Join(filepath.Dir(*configPath), logDir)
		logDir, err = filepath.Abs(logDir)
		if err != nil {
			panic(err)
		}
	}

	logDirStat, err := os.Stat(logDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(logDir, 644)
			if err != nil {
				panic(err)
			}
		}
		panic(err)
	} else {
		if !logDirStat.IsDir() {
			panic(fmt.Sprintf("%s is not a directry", logDir))
		}
	}

	logPath := filepath.Join(logDir, "cleaner.log")
	logF, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	stdoutPath := filepath.Join(logDir, "cleaner.out")
	outF, err := os.OpenFile(stdoutPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	os.Stdout = outF
	os.Stderr = outF
}
