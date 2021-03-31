package main

import (
	"fmt"
	"github.com/ClessLi/fastdfs_cleaner/internal/pkg/fastdfs_cleaner"
	"os"
)

func main() {
	cleaner := fastdfs_cleaner.NewCleanerFromConfig()
	err := cleaner.Clean()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
