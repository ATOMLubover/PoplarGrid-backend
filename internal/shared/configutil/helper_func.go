package configutil

import "os"

// LoadEnvVariable 辅助函数获取环境变量
func LoadEnvVariable(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
