package helperUtils

import (
	"github.com/spf13/viper"
	"log"
)

type AppConfig struct {
	Postgres      Postgres    `yaml:"postgres"`
	ContainerPort string      `yaml:"containerPort"`
	DataSources   DataSources `yaml:"dataSources"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

type DataSources struct {
	Buckets Buckets `yaml:"buckets"`
}

type Buckets struct {
	Shoes string `yaml:"shoes"`
}

var config *AppConfig

func InitConfig() *AppConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")      // look for config locally for development
	viper.AddConfigPath("/config") // look for config when running in the cluster (volume mount)

	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Panic(err)
	}

	return config

}
