package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

var glbConfig *GlobalConfig

type GlobalConfig struct {
	App   App   `yaml:"app"`
	Mysql Mysql `yaml:"mysql"`
	Log   Log   `yaml:"log"`
}

type App struct {
	Name string `yaml:"name"`
}

type Mysql struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	User   string `yaml:"user"`
	Passwd string `yaml:"passwd"`
	DBName string `yaml:"DBName"`
}

type Log struct {
	Level       string `yaml:"level"`
	Env         string `yaml:"env"`
	MaxAge      int    `yaml:"maxAge"`
	MaxBackup   int    `yaml:"maxBackup"`
	MaxFileSize int    `yaml:"maxFileSize"`
	LogFileDir  string `yaml:"logFileDir"`
}

var dftConfig = GlobalConfig{
	App: App{
		Name: "mybutler",
	},
	Mysql: Mysql{
		Host:   "127.0.0.1",
		Port:   3306,
		User:   "root",
		Passwd: "1ASDinnocent",
		DBName: "mybutler",
	},
	Log: Log{
		Level:       "debug",
		Env:         "dev",
		MaxAge:      10,
		MaxBackup:   10,
		MaxFileSize: 128,
		LogFileDir:  "./logs",
	},
}

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./backend/pkg/configs")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("=== warning: config file not found === %v\n", err)
			// use default config
			glbConfig = &dftConfig
			return nil
		} else {
			return err
		}
	}
	glbConfig = &GlobalConfig{}
	if err := viper.Unmarshal(glbConfig); err != nil {
		return err
	}
	return nil
}

func Default() *GlobalConfig {
	return glbConfig
}
