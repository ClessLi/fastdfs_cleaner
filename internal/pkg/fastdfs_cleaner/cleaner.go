package fastdfs_cleaner

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"os"
	"sync"
	"time"
)

type remover interface {
	Remove(filepath string) error
}

type osRemover struct {
}

func (o osRemover) Remove(filepath string) error {
	return os.Remove(filepath)
}

type Cleaner struct {
	poolCap     int
	fileRemover remover
	locker      sync.Locker
	storage     Storage
}

func NewCleaner(poolCap int, rmr remover, storage Storage) Cleaner {
	return Cleaner{
		poolCap:     poolCap,
		fileRemover: rmr,
		locker:      new(sync.Mutex),
		storage:     storage,
	}
}

func NewCleanerFromConfig() Cleaner {
	return NewCleaner(GetSingletonConfigInstance().CleanerGoroutineCap, new(osRemover), NewMySQLStorageFromConfig())
}

func (c *Cleaner) Clean() error {
	c.locker.Lock()
	defer c.locker.Unlock()
	pool, _ := ants.NewPool(c.poolCap)
	defer pool.Release()
	garbageInfos := c.storage.GetAllGarbageInfo()
	c.backgroundClean(pool, garbageInfos)
	return nil
}

func (c Cleaner) backgroundClean(pool *ants.Pool, garbageInfos []GarbageInfo) {
	for i := range garbageInfos {
		idx := i
		filePath := garbageInfos[idx].GetFilePath()

		for pool.Free() <= 0 {
			time.Sleep(time.Millisecond * 10)
		}
		err := pool.Submit(func() {
			//err := os.Remove(filePath)
			err := c.fileRemover.Remove(filePath)
			if err != nil {
				fmt.Println(garbageInfos[idx], err)
				return
			}
			fmt.Println(garbageInfos[idx].GetFilePath(), "is removed in file system.")
			c.storage.RemoveGarbageInfo(garbageInfos[idx])
		})
		if err != nil {
			fmt.Println(garbageInfos[idx], err)
		}
	}
}
