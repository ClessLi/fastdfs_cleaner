package fastdfs_cleaner

import (
	"gopkg.in/yaml.v2"
	"sync"
)

var (
	onceForSingletonConfig  sync.Once
	singletonConfigInstance *Config
	configFilepath          string
)

type Config struct {
	CleanerGoroutineCap int    `yaml:"CleanerGoroutineCap"`
	FastDfsStoragePath  string `yaml:"FastDfsStoragePath"`
	DBType              string `yaml:"DBType"`
	DatabaseName        string `yaml:"DatabaseName"`
	TableName           string `yaml:"TableName"`
	//Fields []string `yaml:"Fields,flow"`
	Field      string `yaml:"Field"`
	Username   string `yaml:"Username"`
	Password   string `yaml:"Password"`
	IPAddr     string `yaml:"IPAddr"`
	ListenPort uint   `yaml:"ListenPort"`
	Protocol   string `yaml:"Protocol"`
}

func GetSingletonConfigInstance() *Config {
	onceForSingletonConfig.Do(func() {
		if singletonConfigInstance == nil {
			configData, err := readFile(configFilepath)
			if err != nil {
				panic("read config failed, cased by: " + err.Error())
			}
			var config Config
			err = yaml.Unmarshal(configData, &config)
			if err != nil {
				panic("unmarshal config failed, cased by: " + err.Error())
			}
			// TODO: 判断Config各必选项都完整
			singletonConfigInstance = &config
		}
	})

	return singletonConfigInstance
}
