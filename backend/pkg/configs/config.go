package configs

import (
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

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./backend/pkg/configs")
	if err := viper.ReadInConfig(); err != nil {
		return err
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
