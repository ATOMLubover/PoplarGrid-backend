package config

import (
	"fmt"
	"path/filepath"
	"poplargrid/internal/shared/configutil"
	"sync"

	"github.com/spf13/viper"
)

// 网关服务器的配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Api      ApiConfig      `mapstructure:"api"`
}

// 总体服务器配置结构体
type ServerConfig struct {
	// 启动模式
	Mode string `mapstructure:"mode"`
}

// 数据库配置结构体
type DatabaseConfig struct {
	// 数据库类型，目前只支持 PostgreSQL
	Type string `mapstructure:"type"`

	// 数据库主机地址和端口
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`

	// 登录用户名
	User string `mapstructure:"user"`
	// 为了提高安全性，密码可以通过环境变量传递
	PwdEnvVar string `mapstructure:"pwd_env_var"`

	// 数据库名称
	DbName string `mapstructure:"db_name"`

	// 是否启用 SSL
	SslEnabled string `mapstructure:"ssl_enabled"`

	// 连接超时时间
	ConnectTimeout int `mapstructure:"connect_timeout"`
}

// API 配置结构体
type ApiConfig struct {
	// 尨译 API 基础 URL
	BaseUrl string `mapstructure:"base_url"`
	// 尨译 API 的授权 Token（初始化时）
	AuthToken string `mapstructure:"auth_token"`
}

var (
	// 全局配置变量
	config *Config = &Config{}
	// 保护 config 的互斥锁
	mtx sync.RWMutex

	// 控制用于热更新的配置加载函数
	configUpdater func(*viper.Viper) = func(v *viper.Viper) {
		mtx.Lock()
		defer mtx.Unlock()

		if config == nil {
			config = &Config{}
		}

		if err := v.Unmarshal(config); err != nil {
			fmt.Printf("配置文件解析失败: %v", err)
		}

		// 在 debug 模式下打印内存中的结构体检查
		if config.Server.Mode == "debug" {
			fmt.Printf("重新加载配置结构体：%+v\n", config)
		}
	}
)

// LoadConfig 使用 configutil 包加载配置文件
func LoadConfig(relPath string, cfgType string) error {
	// 首先检查路径的绝对路径
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return err
	}

	// 使用 configutil 包加载配置文件，并且注入热更新的实现函数
	if err := configutil.LoadConfig(config, cfgType, absPath, configUpdater); err != nil {
		return err
	}

	// 在 debug 模式下打印内存中的结构体检查
	if config.Server.Mode == "debug" {
		fmt.Printf("已加载配置结构体：%+v\n", config)
	}

	return nil
}

// GetConfig 返回全局配置
func GetConfig() *Config {
	mtx.RLock()
	defer mtx.RUnlock()

	if config == nil {
		return nil
	}

	return config
}
