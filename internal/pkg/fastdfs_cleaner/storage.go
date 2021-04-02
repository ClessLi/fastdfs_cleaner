package fastdfs_cleaner

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

type TabFastDfs struct {
	Col string
}

type Storage interface {
	RemoveGarbageInfo(info GarbageInfo)
	GetAllGarbageInfo() []GarbageInfo
}

type mysqlStorage struct {
	db         *gorm.DB
	deleteBuff GarbageInfosQueue
	rwLocker   *sync.RWMutex
}

func newMySQLStorage(db *gorm.DB, queue GarbageInfosQueue) Storage {
	return &mysqlStorage{
		db:         db,
		deleteBuff: queue,
		rwLocker:   new(sync.RWMutex),
	}
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
	return newMySQLStorage(mysqlDB, NewGarbageInfosQueue())
}

func (m *mysqlStorage) RemoveGarbageInfo(info GarbageInfo) {
	// DONE: Remove Method
	m.rwLocker.Lock()
	defer m.rwLocker.Unlock()
	config := GetSingletonConfigInstance()
	colValue := info.GetRelativePath()
	m.db.Exec(fmt.Sprintf("delete from %s where `%s` = '%s'", config.TableName, config.Field, colValue))
	var count int64
	m.db.Table(config.TableName).Select(config.Field).Where("? = ?", config.Field, colValue).Count(&count)
	if count > 0 {
		fmt.Printf("garbage data '%s' is not removed in database", colValue)
	}
	fmt.Printf("garbage data '%s' is removed in database", colValue)
}

func (m mysqlStorage) GetAllGarbageInfo() []GarbageInfo {
	// DONE: Get All Method
	var (
		config = GetSingletonConfigInstance()
		infos  = make([]GarbageInfo, 0)
	)
	m.rwLocker.RLock()
	defer m.rwLocker.RUnlock()

	rows, err := m.db.Table(config.TableName).Select(config.Field).Limit(1000).Rows()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for rows.Next() {
		var value string
		err = rows.Scan(&value)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		infos = append(infos, newRelativePathGarbageInfo(config.FastDfsStoragePath, value))
	}
	return infos
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
