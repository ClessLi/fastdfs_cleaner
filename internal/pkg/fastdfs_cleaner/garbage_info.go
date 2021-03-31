package fastdfs_cleaner

import (
	"path/filepath"
)

type GarbageInfo interface {
	GetFilePath() string
}

type relativePathGarbageInfo struct {
	workspace    string
	relativePath string
}

func newRelativePathGarbageInfo(workspace, relativePath string) GarbageInfo {
	return relativePathGarbageInfo{
		workspace:    workspace,
		relativePath: relativePath,
	}
}

func (r relativePathGarbageInfo) GetFilePath() string {
	return filepath.Join(r.workspace, r.relativePath)
}
