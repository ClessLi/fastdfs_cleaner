package fastdfs_cleaner

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

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
	// TODO: Remove Method
	panic("implement me")
}

func (m mysqlStorage) GetAllGarbageInfo() []GarbageInfo {
	// TODO: Get All Method
	panic("implement me")
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
