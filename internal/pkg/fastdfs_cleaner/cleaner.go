package fastdfs_cleaner

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type remover interface {
	Remove(filepath string) error
}

type osRemover struct {
}

func (o osRemover) Remove(path string) error {
	err := os.Remove(path)
	if err != nil {
		fmt.Printf("remove file %s failed, cased by %s\n", path, err)
		dir := filepath.Dir(path)
		_, err = os.Stat(dir)
		if err == nil || os.IsExist(err) {
			return nil
		}
		return fmt.Errorf("file directory '%s' check error: %s", dir, err)
	}
	fmt.Println(path, "is removed in file system.")
	return nil
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

	// build pool
	pool, _ := ants.NewPool(c.poolCap)
	defer pool.Release()

	// loop get and clean garbage infos, until empty.
	for garbageInfos := c.storage.GetAllGarbageInfo(); len(garbageInfos) > 0; garbageInfos = c.storage.GetAllGarbageInfo() {
		c.backgroundClean(pool, garbageInfos)
		c.waitGoroutines(pool)
	}

	return nil
}

func (c Cleaner) backgroundClean(pool *ants.Pool, garbageInfos []GarbageInfo) {
	for i := range garbageInfos {
		idx := i
		filePath := garbageInfos[idx].GetFilePath()

		// wait per 10ms, if pool is full
		for pool.Free() <= 0 {
			time.Sleep(time.Millisecond * 10)
		}

		// submit remove task into pool
		err := pool.Submit(func() {
			//err := os.Remove(filePath)
			err := c.fileRemover.Remove(filePath)
			if err != nil {
				fmt.Printf("%s removed failed in file system, cased by: %s", filePath, err)
				return
			}
			c.storage.RemoveGarbageInfo(garbageInfos[idx])
		})
		if err != nil {
			fmt.Println(garbageInfos[idx], err)
		}

	}
}

func (c *Cleaner) waitGoroutines(pool *ants.Pool) {
	// wait per 10ms, if pool is not empty
	for pool.Running() > 0 {
		time.Sleep(time.Millisecond * 10)
	}
	// wait goroutines, if task num less than pool cap
	//if taskNum < c.poolCap {
	//}
}
