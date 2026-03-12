package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// GlobalConfig 全局配置实例
var GlobalConfig *Config

// LoadConfig 加载配置并设置全局实例
func LoadConfig() *Config {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}
	GlobalConfig = cfg
	return cfg
}

// Config 应用配置结构
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	CORS     CORSConfig     `mapstructure:"cors"`
	Upload   UploadConfig   `mapstructure:"upload"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string `mapstructure:"name"`
	Env     string `mapstructure:"env"`
	Version string `mapstructure:"version"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Mode string `mapstructure:"mode"`
	Port string `mapstructure:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	SSLMode         string `mapstructure:"sslmode"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime"`
}

// DSN 返回数据库连接字符串
func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// Addr 返回 Redis 地址
func (r *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret         string        `mapstructure:"secret"`
	ExpireHours    int           `mapstructure:"expire_hours"`
	RefreshHours   int           `mapstructure:"refresh_hours"`
	AccessTokenTTL time.Duration `mapstructure:"-"`
	RefreshTokenTTL time.Duration `mapstructure:"-"`
}

// AfterLoad 配置加载后初始化派生字段
func (j *JWTConfig) AfterLoad() {
	j.AccessTokenTTL = time.Duration(j.ExpireHours) * time.Hour
	j.RefreshTokenTTL = time.Duration(j.RefreshHours) * time.Hour
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowOrigins  string `mapstructure:"allow_origins"`
	AllowMethods  string `mapstructure:"allow_methods"`
	AllowHeaders  string `mapstructure:"allow_headers"`
}

// GetAllowOrigins 返回允许的源列表
func (c *CORSConfig) GetAllowOrigins() []string {
	return strings.Split(c.AllowOrigins, ",")
}

// GetAllowMethods 返回允许的方法列表
func (c *CORSConfig) GetAllowMethods() []string {
	return strings.Split(c.AllowMethods, ",")
}

// GetAllowHeaders 返回允许的头部列表
func (c *CORSConfig) GetAllowHeaders() []string {
	return strings.Split(c.AllowHeaders, ",")
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	MaxSize      int64  `mapstructure:"max_size"`
	AllowedTypes string `mapstructure:"allowed_types"`
	StoragePath  string `mapstructure:"storage_path"`
}

// GetAllowedTypes 返回允许的文件类型列表
func (u *UploadConfig) GetAllowedTypes() []string {
	return strings.Split(u.AllowedTypes, ",")
}

// Load 加载配置
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/app/config")
	viper.AddConfigPath(".")

	// 设置默认值
	setDefaults()

	// 从环境变量读取
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 初始化派生字段
	cfg.JWT.AfterLoad()

	return &cfg, nil
}

// setDefaults 设置默认配置值
func setDefaults() {
	// 应用默认配置
	viper.SetDefault("app.name", "ConsultingSystem")
	viper.SetDefault("app.env", "development")
	viper.SetDefault("app.version", "1.0.0")

	// 服务器默认配置
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("server.port", "8080")

	// 数据库默认配置
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "consulting")
	viper.SetDefault("database.password", "consulting123")
	viper.SetDefault("database.dbname", "consulting_db")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "5m")

	// Redis 默认配置
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)

	// JWT 默认配置
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expire_hours", 24)
	viper.SetDefault("jwt.refresh_hours", 168)

	// 日志默认配置
	viper.SetDefault("log.level", "debug")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")

	// CORS 默认配置
	viper.SetDefault("cors.allow_origins", "http://localhost:5173,http://localhost:3000")
	viper.SetDefault("cors.allow_methods", "GET,POST,PUT,DELETE,OPTIONS")
	viper.SetDefault("cors.allow_headers", "Authorization,Content-Type,X-Request-ID")

	// 上传默认配置
	viper.SetDefault("upload.max_size", 10485760)
	viper.SetDefault("upload.allowed_types", "image/jpeg,image/png,image/gif,application/pdf")
	viper.SetDefault("upload.storage_path", "./uploads")
}
