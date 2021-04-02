package main

import (
	"flag"
	"github.com/ClessLi/fastdfs_cleaner/internal/pkg/fastdfs_cleaner"
	"os"
)

var (
	configPath = flag.String("f", "", "the cleaner `config`uration file path.")
	help       = flag.Bool("h", false, "this `help`")
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
}
