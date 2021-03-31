package fastdfs_cleaner

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"os"
	"sync"
	"time"
)

type Cleaner struct {
	poolCap int
	locker  sync.Locker
	storage Storage
}

func NewCleaner(poolCap int, storage Storage) Cleaner {
	return Cleaner{
		poolCap: poolCap,
		locker:  new(sync.Mutex),
		storage: storage,
	}
}

func NewCleanerFromConfig() Cleaner {
	return NewCleaner(GetSingletonConfigInstance().CleanerGoroutineCap, NewMySQLStorageFromConfig())
}

func (c *Cleaner) Clean() error {
	c.locker.Lock()
	defer c.locker.Unlock()
	pool, _ := ants.NewPool(c.poolCap)
	defer pool.Release()
	pool.Running()
	garbageInfos := c.storage.GetAllGarbageInfo()
	for _, garbageInfo := range garbageInfos {
		filePath := garbageInfo.GetFilePath()
		for pool.Cap() >= c.poolCap {
			time.Sleep(time.Millisecond * 100)
		}
		err := pool.Submit(func() {
			err := os.Remove(filePath)
			if err != nil {
				fmt.Println(garbageInfo, err)
				return
			}
			c.storage.RemoveGarbageInfo(garbageInfo)
		})
		if err != nil {
			fmt.Println(garbageInfo, err)
		}
	}
	return nil
}
