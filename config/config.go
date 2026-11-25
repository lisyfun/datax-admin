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
	Security SecurityConfig `mapstructure:"security"`
	SSH      SSHConfig      `mapstructure:"ssh"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
}

type ServerConfig struct {
	Port           string   `mapstructure:"port"`
	Mode           string   `mapstructure:"mode"`
	BasePath       string   `mapstructure:"base_path"`
	MaxFileSize    int64    `mapstructure:"max_file_size"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	// TrustedProxies 可信反向代理的 IP/CIDR 列表。
	// 为空时不信任任何代理的 X-Forwarded-For，避免客户端伪造来源 IP。
	TrustedProxies []string `mapstructure:"trusted_proxies"`
}

type DatabaseConfig struct {
	Type            string `mapstructure:"type"` // 数据库类型: mysql, postgres
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	LogMode         string `mapstructure:"log_mode"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	TimeZone        string `mapstructure:"time_zone"`
	Charset         string `mapstructure:"charset"`
	SSLMode         string `mapstructure:"ssl_mode"` // PostgreSQL SSL模式
	Schema          string `mapstructure:"schema"`   // PostgreSQL Schema
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
	Secret        string `mapstructure:"secret"`         // Session 密钥（建议用 AUTH_SECRET 环境变量注入）
	Expiration    int    `mapstructure:"expiration"`     // Session 过期时间(秒)
	AllowRegister bool   `mapstructure:"allow_register"` // 是否允许公开注册（默认 false）
	CookieSecure  bool   `mapstructure:"cookie_secure"`  // 是否只通过 HTTPS 发送会话 Cookie
}

// SecurityConfig 敏感数据加密配置
type SecurityConfig struct {
	// EncryptionKey 用于加密数据库中的 SSH/Kafka/Redis 凭据。
	// 为空时使用 SECURITY_ENCRYPTION_KEY 环境变量，或自动生成 data/encryption.key。
	EncryptionKey string `mapstructure:"encryption_key"`
}

// KafkaConfig Kafka 客户端安全配置
type KafkaConfig struct {
	// InsecureSkipVerify 跳过 TLS 证书与主机名校验。
	// 危险：仅建议在测试环境使用，生产环境应配置 ca_file 或使用受信任证书。
	InsecureSkipVerify bool `mapstructure:"insecure_skip_verify"`
	// CAFile 自签 CA 证书路径（PEM），用于校验 Kafka broker 证书。
	CAFile string `mapstructure:"ca_file"`
}

// SSHConfig SSH 连接安全配置
type SSHConfig struct {
	KnownHostsPath        string `mapstructure:"known_hosts_path"`         // known_hosts 路径，为空使用 ~/.ssh/known_hosts
	InsecureIgnoreHostKey bool   `mapstructure:"insecure_ignore_host_key"` // 危险：跳过主机密钥校验，仅建议测试环境使用
}

// JobLogConfig 任务日志配置
type JobLogConfig struct {
	MaxAge      int  `mapstructure:"max_age"`      // 日志保留天数
	AutoCleanup bool `mapstructure:"auto_cleanup"` // 是否启用自动清理
}

type InitDataConfig struct {
	ForceReset      bool   `mapstructure:"force_reset"`      // 是否强制重新初始化数据
	DefaultPassword string `mapstructure:"default_password"` // 默认管理员密码
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

	// 配置中的敏感项允许通过环境变量注入。
	// 注意：viper 的 AutomaticEnv 不会作用于 Unmarshal，这里显式回读。
	if GlobalConfig.Auth.Secret == "" {
		GlobalConfig.Auth.Secret = viper.GetString("auth.secret")
	}
	if GlobalConfig.Database.Password == "" {
		GlobalConfig.Database.Password = viper.GetString("db.password")
	}
	if GlobalConfig.Security.EncryptionKey == "" {
		GlobalConfig.Security.EncryptionKey = viper.GetString("security.encryption_key")
	}
	if GlobalConfig.Init.DefaultPassword == "" {
		GlobalConfig.Init.DefaultPassword = viper.GetString("init.default_password")
	}
	if GlobalConfig.SSH.KnownHostsPath == "" {
		GlobalConfig.SSH.KnownHostsPath = viper.GetString("ssh.known_hosts_path")
	}

	// 布尔项：环境变量优先级高于配置文件
	if v := viper.GetString("auth.allow_register"); v != "" {
		GlobalConfig.Auth.AllowRegister = parseBool(v, GlobalConfig.Auth.AllowRegister)
	}
	if v := viper.GetString("auth.cookie_secure"); v != "" {
		GlobalConfig.Auth.CookieSecure = parseBool(v, GlobalConfig.Auth.CookieSecure)
	}
	if v := viper.GetString("ssh.insecure_ignore_host_key"); v != "" {
		GlobalConfig.SSH.InsecureIgnoreHostKey = parseBool(v, GlobalConfig.SSH.InsecureIgnoreHostKey)
	}
}

func parseBool(v string, fallback bool) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}
