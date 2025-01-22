package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"db"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	DataX    DataXConfig    `mapstructure:"datax"`
}

type ServerConfig struct {
	Port     string `mapstructure:"port"`
	Mode     string `mapstructure:"mode"`
	BasePath string `mapstructure:"base_path"`
}

type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	LogMode      string `mapstructure:"log_mode"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxLifetime  int    `mapstructure:"max_lifetime"`
	SSLMode      string `mapstructure:"ssl_mode"`
	TimeZone     string `mapstructure:"time_zone"`
	Charset      string `mapstructure:"charset"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"` // token过期时间（小时）
}

type DataXConfig struct {
	Bin  string `yaml:"bin"`  // DataX可执行文件路径
	Home string `yaml:"home"` // DataX安装目录
}

var GlobalConfig Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.AutomaticEnv()

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	if err := viper.Unmarshal(&GlobalConfig); err != nil {
		log.Fatalf("Error unmarshaling config: %s", err)
	}
}
