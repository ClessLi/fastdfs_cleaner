package fastdfs_cleaner

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"sync"
	"time"
)

type Storage interface {
	RemoveGarbageInfo(info GarbageInfo)
	GetAllGarbageInfo() []GarbageInfo
}

type mysqlStorage struct {
	db         *gorm.DB
	rowLimit   uint
	deleteBuff GarbageInfosQueue
	rwLocker   *sync.RWMutex
}

func newMySQLStorage(db *gorm.DB, rowLimit uint, queue GarbageInfosQueue) Storage {
	storage := &mysqlStorage{
		db:         db,
		rowLimit:   rowLimit,
		deleteBuff: queue,
		rwLocker:   new(sync.RWMutex),
	}
	go storage.removeGarbageInfos()
	return storage
}

func NewMySQLStorageFromConfig() Storage {
	config := GetSingletonConfigInstance()
	mysqlDSN := fmt.Sprintf(
		"%s:%s@%s(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Protocol,
		config.IPAddr,
		config.ListenPort,
		config.DatabaseName,
	)
	mysqlDB, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return newMySQLStorage(mysqlDB, 1000, NewGarbageInfosQueue())
}

func (m *mysqlStorage) RemoveGarbageInfo(info GarbageInfo) {
	// DONE: Remove Method
	m.rwLocker.Lock()
	defer m.rwLocker.Unlock()
	m.deleteBuff.Append(info)
}

func (m *mysqlStorage) GetAllGarbageInfo() []GarbageInfo {
	// DONE: Get All Method
	var (
		config = GetSingletonConfigInstance()
		infos  = make([]GarbageInfo, 0)
	)
	m.rwLocker.RLock()
	defer m.rwLocker.RUnlock()

	rows, err := m.db.Table(config.TableName).Select(config.IndexField, config.Field).Limit(int(m.rowLimit)).Rows()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for rows.Next() {
		var value string
		var idx string
		err = rows.Scan(&idx, &value)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		infos = append(infos, newRelativePathGarbageInfo(config.FastDfsStoragePath, value, idx))
	}
	return infos
}

func (m *mysqlStorage) removeGarbageInfos() {
	timeTicker := time.NewTicker(time.Second)
	buffTicker := time.NewTicker(time.Millisecond)
	for true {
		select {
		case <-timeTicker.C:
			if m.deleteBuff.IsEmpty() {
				continue
			}
			break
		case <-buffTicker.C:
			if m.deleteBuff.Size() < int(m.rowLimit) {
				continue
			}
			break
		}
		m.rwLocker.Lock()
		indexesSql := ""
		var records int64 = 0
		for !m.deleteBuff.IsEmpty() {
			indexesSql += fmt.Sprintf("'%s', ", m.deleteBuff.Pop().GetIndex())
			records++
		}
		indexesSql = strings.TrimRight(indexesSql, ", ")
		config := GetSingletonConfigInstance()
		sql := fmt.Sprintf("delete from %s where `%s` in (%s)", config.TableName, config.IndexField, indexesSql)
		m.db.Exec(sql)
		var count int64
		m.db.Table(config.TableName).Select(config.Field).Where("? in (?)", config.IndexField, indexesSql).Count(&count)
		if count > 0 {
			fmt.Printf("%d garbage record(s) removed failed in database.\n", count)
			records -= count
		}
		fmt.Printf("%d garbage record(s) removed in database.\n", records)
		m.rwLocker.Unlock()
	}
}

type GarbageInfosQueue interface {
	Size() int
	Append(info GarbageInfo)
	Pop() GarbageInfo
	IsEmpty() bool
}

type garbageInfosQueue []GarbageInfo

func (g garbageInfosQueue) Size() int {
	return len(g)
}

func (g *garbageInfosQueue) Append(info GarbageInfo) {
	*g = append(*g, info)
}

func (g *garbageInfosQueue) Pop() GarbageInfo {
	p := (*g)[0]
	*g = (*g)[1:]
	return p
}

func (g garbageInfosQueue) IsEmpty() bool {
	return g.Size() == 0
}

func NewGarbageInfosQueue() GarbageInfosQueue {
	return new(garbageInfosQueue)
}
