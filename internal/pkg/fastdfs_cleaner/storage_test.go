package fastdfs_cleaner

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"testing"
)

func Test_mysqlStorage_GetAllGarbageInfo(t *testing.T) {
	configFilepath = "F:\\GO_Project\\src\\fastdfs_cleaner\\test\\config\\cleaner_config.yml"
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
		t.Fatal(err)
	}

	type fields struct {
		db         *gorm.DB
		deleteBuff GarbageInfosQueue
		rwLocker   *sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		//want   []GarbageInfo
	}{
		{
			name: "test get garbage infos",
			fields: fields{
				db:         mysqlDB,
				deleteBuff: NewGarbageInfosQueue(),
				rwLocker:   new(sync.RWMutex),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mysqlStorage{
				db:         tt.fields.db,
				deleteBuff: tt.fields.deleteBuff,
				rwLocker:   tt.fields.rwLocker,
			}
			//if got := m.GetAllGarbageInfo(); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("GetAllGarbageInfo() = %v, want %v", got, tt.want)
			//}
			got := m.GetAllGarbageInfo()
			t.Log(got)
		})
	}
}

func Test_mysqlStorage_RemoveGarbageInfo(t *testing.T) {
	configFilepath = "F:\\GO_Project\\src\\fastdfs_cleaner\\test\\config\\cleaner_config.yml"
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
		t.Fatal(err)
	}
	type fields struct {
		db         *gorm.DB
		deleteBuff GarbageInfosQueue
		rwLocker   *sync.RWMutex
	}
	type args struct {
		info GarbageInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test remove garbage data",
			fields: fields{
				db:         mysqlDB,
				deleteBuff: NewGarbageInfosQueue(),
				rwLocker:   new(sync.RWMutex),
			},
			args: args{info: newRelativePathGarbageInfo(config.FastDfsStoragePath, "group1/1/2/kljxklf.pdf")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &mysqlStorage{
				db:         tt.fields.db,
				deleteBuff: tt.fields.deleteBuff,
				rwLocker:   tt.fields.rwLocker,
			}
			m.RemoveGarbageInfo(tt.args.info)
			t.Log(m.GetAllGarbageInfo())
		})
	}
}
