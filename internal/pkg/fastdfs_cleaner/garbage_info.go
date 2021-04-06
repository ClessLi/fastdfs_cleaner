package fastdfs_cleaner

import (
	"path/filepath"
)

type GarbageInfo interface {
	GetFilePath() string
	GetIndex() string
}

type relativePathGarbageInfo struct {
	workspace    string
	relativePath string
	index        string
}

func newRelativePathGarbageInfo(workspace, relativePath string, index string) GarbageInfo {
	return relativePathGarbageInfo{
		workspace:    workspace,
		relativePath: relativePath,
		index:        index,
	}
}

func (r relativePathGarbageInfo) GetFilePath() string {
	return filepath.Join(r.workspace, r.relativePath)
}

func (r relativePathGarbageInfo) GetIndex() string {
	return r.index
}
