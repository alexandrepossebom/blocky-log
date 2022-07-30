package config

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	confOnce sync.Once
	conf     *ConfYaml
)

type ConfYaml struct {
	Hours            int          `yaml:"hours"`
	DefaultEventType string       `yaml:"defaulteventtype"`
	Blockys          []BlockyItem `yaml:"blockys"`
}

type SectionDatabase struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	Type     string `yaml:"type"`
	Database string `yaml:"database"`
}

type BlockyItem struct {
	Name     string          `yaml:"name"`
	Database SectionDatabase `yaml:"database"`
}

func Get() *ConfYaml {
	confOnce.Do(func() {
		loadConf()
	})
	return conf
}

func GetDatabase(host string) *SectionDatabase {
	for _, blocky := range conf.Blockys {
		if blocky.Name == host {
			return &blocky.Database
		}
	}
	return nil
}

func save(filename string) error {
	v := viper.New()
	v.SetConfigType("yaml")

	b, err := yaml.Marshal(&conf)
	if err != nil {
		return err
	}

	defaultConfig := bytes.NewReader(b)
	if err := v.MergeConfig(defaultConfig); err != nil {
		return err
	}

	return v.WriteConfigAs(filename)
}

func loadConf() {
	filename := "blocky-log.yaml"
	conf = &ConfYaml{}
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc")
	viper.AddConfigPath("$HOME/etc")
	viper.SetConfigName(filename)
	viper.SetDefault("hours", "1")
	viper.SetDefault("defaulteventtype", "all")
	viper.SetDefault("blockys",
		map[string]any{
			"name": "home",
			"database": map[string]any{
				"type":     "mysql",
				"host":     "localhost",
				"port":     3306,
				"username": "blocky",
				"password": "secret",
				"database": "blocky",
			},
		})

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.Unmarshal(&conf); err != nil {
				panic(err)
			}
			if err := save(filename); err != nil {
				panic(err)
			}
			fmt.Printf("Config file not found, new config file is generated [%s] configure it.\n", filename)
			os.Exit(1)
		}
		fmt.Println(err)
		os.Exit(1)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		panic(err)
	}
}
