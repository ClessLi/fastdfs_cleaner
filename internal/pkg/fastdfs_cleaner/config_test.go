package fastdfs_cleaner

import (
	"gopkg.in/yaml.v2"
	"reflect"
	"testing"
)

func TestGetSingletonConfigInstance(t *testing.T) {
	configFilepath = "F:\\GO_Project\\src\\fastdfs_cleaner\\test\\config\\cleaner_config.yml"
	var testConfig Config
	data, err := readFile(configFilepath)
	if err != nil {
		t.Fatal(err)
	}

	err = yaml.Unmarshal(data, &testConfig)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "test get singleton config instance",
			want: &testConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSingletonConfigInstance(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSingletonConfigInstance() = %v, want %v", got, tt.want)
			} else {
				t.Log("got:", got)
				t.Log("want:", tt.want)
			}
		})
	}
}
