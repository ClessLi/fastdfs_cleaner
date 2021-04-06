package fastdfs_cleaner

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type mockRemover struct {
}

func (m mockRemover) Remove(filepath string) error {
	time.Sleep(time.Millisecond)
	fmt.Println(filepath, "is removed.")
	return nil
}

type mockStorage struct {
	once *sync.Once
}

func (m mockStorage) RemoveGarbageInfo(info GarbageInfo) {
	fmt.Println("remove garbage file: " + "'" + info.GetFilePath() + "' in storage.")
}

func (m mockStorage) GetAllGarbageInfo() []GarbageInfo {
	infos := make([]GarbageInfo, 0)
	ws := "/home/fastdfs/storage/"
	m.once.Do(func() {
		for i := 0; i < 100; i++ {
			//relativePath := fmt.Sprintf("group1/M00/%d/B9/CgFMY2BcAi-AHYzqAAbghrRaoKw831.pdf", i)
			infos = append(infos, newRelativePathGarbageInfo(ws, fmt.Sprintf("group1/M00/%d/B9/CgFMY2BcAi-AHYzqAAbghrRaoKw831.pdf", i), i))
		}
	})
	return infos
}

func TestCleaner_Clean(t *testing.T) {
	type fields struct {
		poolCap     int
		fileRemover remover
		locker      sync.Locker
		storage     Storage
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "test 10 capacity goroutine pool",
			fields: fields{
				poolCap:     10,
				fileRemover: new(mockRemover),
				locker:      new(sync.Mutex),
				storage:     &mockStorage{once: new(sync.Once)},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cleaner{
				poolCap:     tt.fields.poolCap,
				fileRemover: tt.fields.fileRemover,
				locker:      tt.fields.locker,
				storage:     tt.fields.storage,
			}
			if err := c.Clean(); (err != nil) != tt.wantErr {
				t.Errorf("Clean() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
