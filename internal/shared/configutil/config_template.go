package configutil

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

// 将配置加载到 viper 实例中
func LoadViperConfig(cfgType, relPath string) (*viper.Viper, error) {
	// 为了方便调试，加入绝对文件路径日志
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return nil, fmt.Errorf("配置文件路径解析错误：%s\n", err.Error())
	}

	v := viper.New()
	v.SetConfigType(cfgType)
	v.SetConfigFile(absPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("配置文件读入 viper 失败：%s", err.Error())
	}

	return v, nil
}

// 将配置加载到结构体中，实现类型安全操作
// 暂时不支持参数设置
func LoadConfig[T any](cfg *T, cfgPath string, cfgType string) error {
	v, err := LoadViperConfig(cfgType, cfgPath)
	if err != nil {
		return err
	}

	return v.Unmarshal(cfg)
}
