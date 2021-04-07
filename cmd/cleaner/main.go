package main

import (
	"fmt"
	"github.com/ClessLi/fastdfs_cleaner/internal/pkg/fastdfs_cleaner"
	"github.com/apsdehal/go-logger"
	"os"
)

func main() {
	//cleaner := fastdfs_cleaner.NewCleanerFromConfig()
	cleaner := fastdfs_cleaner.NewCleanerFromConfigWithLogger(logF, logger.InfoLevel)
	err := cleaner.Clean()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
