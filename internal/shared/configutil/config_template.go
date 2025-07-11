package configutil

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 将配置加载到 viper 实例中
func LoadViperConfig(cfgType, relPath string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType(cfgType)
	v.SetConfigFile(relPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("配置文件读入 viper 失败：%s", err.Error())
	}

	return v, nil
}

// 将配置加载到结构体中，实现类型安全操作
// updater 函数用于注入热更新的实现函数
func LoadConfig[T any](cfg *T, cfgType, cfgPath string, updater func(v *viper.Viper)) error {
	v, err := LoadViperConfig(cfgType, cfgPath)
	if err != nil {
		return err
	}

	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("配置文件解析到结构体失败：%s", err.Error())
	}

	// 如果未提供 updater 函数，则跳过热更新
	if updater != nil {
		v.WatchConfig()
		// 当配置文件发生变化时，调用 updater 函数
		v.OnConfigChange(func(e fsnotify.Event) { updater(v) })
	}

	return nil
}
