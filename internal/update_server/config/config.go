package config

// Config 是更新服务器的配置结构体
type Config struct {
	// Server 配置
	Server ServerConfig `mapstructure:"server"`
}

// ServerConfig 是更新服务器的服务器配置
type ServerConfig struct {
	// Mode 运行模式
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
