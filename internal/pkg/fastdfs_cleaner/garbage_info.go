package fastdfs_cleaner

import (
	"path/filepath"
)

type GarbageInfo interface {
	GetFilePath() string
	GetIndex() interface{}
}

type relativePathGarbageInfo struct {
	workspace    string
	relativePath string
	index        int
}

func newRelativePathGarbageInfo(workspace, relativePath string, index int) GarbageInfo {
	return relativePathGarbageInfo{
		workspace:    workspace,
		relativePath: relativePath,
		index:        index,
	}
}

func (r relativePathGarbageInfo) GetFilePath() string {
	return filepath.Join(r.workspace, r.relativePath)
}

func (r relativePathGarbageInfo) GetIndex() interface{} {
	return r.index
}
