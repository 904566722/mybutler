package configs

import "github.com/spf13/viper"

var glbConfig *GlobalConfig

type GlobalConfig struct {
	Mysql Mysql `yaml:"mysql"`
}

type Mysql struct {
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	User   string `yaml:"user"`
	Passwd string `yaml:"passwd"`
	DBName string `yaml:"DBName"`
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

func GetGlbConfig() *GlobalConfig {
	return glbConfig
}
