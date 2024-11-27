package config

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

type DBConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	Username     string        `yaml:"username"`
	Password     string        `yaml:"password"`
	Database     string        `yaml:"database"`
	Charset      string        `yaml:"charset"`
	MaxIdleConns int           `yaml:"max_idle_conns"` // 最大空闲连接数
	MaxOpenConns int           `yaml:"max_open_conns"` // 最大打开连接数
	MaxLifetime  time.Duration `yaml:"max_lifetime"`   // 连接最大生命周期
}

type Config struct {
	Host string   `yaml:"host"`
	Port int      `yaml:"port"`
	DB   DBConfig `yaml:"db"`
}

func LoadConfig(path string) (*Config, error) {
	conf := &Config{}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
