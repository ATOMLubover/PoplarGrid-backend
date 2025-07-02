package config

// Config 目前适配 viper 的配置文件格式（mapstructure）
// 直接支持 YAML、JSON 等格式，
// 未来可以扩展为支持更多配置文件格式（但是似乎没必要）

// 网关服务器的配置结构体
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
}

// 总体服务器配置结构体
type ServerConfig struct {
	// 启动模式
	Mode string `mapstructure:"mode"`

	// 控制服务器的地址和端口
	Port int `mapstructure:"port"`
}

// 数据库配置结构体
type DatabaseConfig struct {
	// 数据库类型，目前只支持 PostgreSQL
	Type string `mapstructure:"type"`

	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`

	User string `mapstructure:"user"`
	// 为了提高安全性，密码可以通过环境变量传递
	PwdEnvVar string `mapstructure:"pwd_env_var"`

	// 数据库名称
	DbName string `mapstructure:"db_name"`

	// 是否启用 SSL
	SslEnabled bool `mapstructure:"ssl_enabled"`

	// 连接超时时间
	ConnectTimeout int `mapstructure:"connect_timeout"`
}
