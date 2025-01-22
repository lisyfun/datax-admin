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
	Logger   LoggerConfig   `mapstructure:"logger"`
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
	Home string `mapstructure:"home"`
	Bin  string `mapstructure:"bin"`
}

type LoggerConfig struct {
	LogPath    string `mapstructure:"log_path"`    // 日志文件路径
	MaxSize    int    `mapstructure:"max_size"`    // 单个日志文件最大大小（MB）
	MaxBackups int    `mapstructure:"max_backups"` // 最大保留的旧日志文件数
	MaxAge     int    `mapstructure:"max_age"`     // 日志文件保留的最大天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩旧日志文件
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
