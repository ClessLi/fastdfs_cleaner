package fastdfs_cleaner

import (
	"github.com/apsdehal/go-logger"
	"os"
)

type baseLogger struct {
	logger *logger.Logger
}

func newBaseLogger(logFile *os.File, level logger.LogLevel) *baseLogger {
	l, err := logger.New("Fastdfs Cleaner", level, logFile)
	if err != nil {
		panic(err)
	}
	l.SetFormat("[%{module}] %{time:2006-01-02 15:04:05.000} [%{level}] %{message}\n")
	return &baseLogger{
		logger: l,
	}
}

type cleanerLogger struct {
	*baseLogger
	cleaner Cleaner
}

func newCleanerLogger(baseLogger *baseLogger, cleaner Cleaner) Cleaner {
	return &cleanerLogger{
		baseLogger: baseLogger,
		cleaner:    cleaner,
	}
}

func CleanerLoggerProxy(logFile *os.File, level logger.LogLevel) func(cleaner Cleaner) Cleaner {
	return func(cleaner Cleaner) Cleaner {
		bl := newBaseLogger(logFile, level)
		return newCleanerLogger(bl, cleaner)
	}
}

func (c *cleanerLogger) Clean() error {
	err := c.cleaner.Clean()
	if err != nil {
		c.baseLogger.logger.WarningF("something wrong with clean fastdfs garbage data, cased by: %s.", err)
		return err
	}
	c.baseLogger.logger.Info("clean fastdfs garbage data success.")
	return nil
}

type storageLogger struct {
	*baseLogger
	storage Storage
}

func newStorageLogger(baseLogger *baseLogger, storage Storage) Storage {
	return &storageLogger{
		baseLogger: baseLogger,
		storage:    storage,
	}
}

func StorageLoggerProxy(logFile *os.File, level logger.LogLevel) func(storage Storage) Storage {
	return func(storage Storage) Storage {
		bl := newBaseLogger(logFile, level)
		return newStorageLogger(bl, storage)
	}
}

func (s *storageLogger) RemoveGarbageInfo(info GarbageInfo) {
	s.storage.RemoveGarbageInfo(info)
	s.baseLogger.logger.InfoF("removed garbage info(%s) from Storage object.", info.GetFilePath())
}

func (s *storageLogger) GetAllGarbageInfo() []GarbageInfo {
	infos := s.storage.GetAllGarbageInfo()
	s.baseLogger.logger.InfoF("get garbage infos(%s), count %d", infos, len(infos))
	return infos
}
