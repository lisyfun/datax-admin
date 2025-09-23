package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"db"`
	DataX    DataXConfig    `mapstructure:"datax"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Auth     AuthConfig     `mapstructure:"auth"`
	JobLog   JobLogConfig   `mapstructure:"job_log"`
	Init     InitDataConfig `mapstructure:"init"`
}

type ServerConfig struct {
	Port        string `mapstructure:"port"`
	Mode        string `mapstructure:"mode"`
	BasePath    string `mapstructure:"base_path"`
	MaxFileSize int64  `mapstructure:"max_file_size"`
}

type DatabaseConfig struct {
	Type         string `mapstructure:"type"` // 数据库类型: mysql, postgres
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	LogMode      string `mapstructure:"log_mode"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxLifetime  int    `mapstructure:"max_lifetime"`
	TimeZone     string `mapstructure:"time_zone"`
	Charset      string `mapstructure:"charset"`
	SSLMode      string `mapstructure:"ssl_mode"` // PostgreSQL SSL模式
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

// AuthConfig Auth配置
type AuthConfig struct {
	Secret     string `mapstructure:"secret"`     // Session 密钥
	Expiration int    `mapstructure:"expiration"` // Session 过期时间(秒)
}

// JobLogConfig 任务日志配置
type JobLogConfig struct {
	MaxAge      int  `mapstructure:"max_age"`      // 日志保留天数
	AutoCleanup bool `mapstructure:"auto_cleanup"` // 是否启用自动清理
}

type InitDataConfig struct {
	ForceReset      bool   `mapstructure:"force_reset"`       // 是否强制重新初始化数据
	DefaultPassword string `mapstructure:"default_password"`  // 默认管理员密码
}

var GlobalConfig Config

func InitConfig() {
	InitConfigWithFile("")
}

func InitConfigWithFile(configFile string) {
	// 如果指定了配置文件，使用指定的文件
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// 默认配置
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./")
	}

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
