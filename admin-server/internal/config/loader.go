package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var (
	// 默认配置文件路径（根据操作系统动态设置）
	DefaultRedisConfigPath      string
	DefaultMiddlewareConfigPath string
)

func init() {
	// 根据操作系统设置默认路径
	if runtime.GOOS == "windows" {
		DefaultRedisConfigPath = "C:\\work\\redis.json"
		DefaultMiddlewareConfigPath = "C:\\work\\middleware.yaml"
	} else {
		// Linux 或其他 Unix 系统
		DefaultRedisConfigPath = "/etc/work/redis.json"
		DefaultMiddlewareConfigPath = "/etc/work/middleware.yaml"
	}
}

// LoadRedisConfig 从文件加载 Redis 配置
func LoadRedisConfig(configPath string) (*RedisConf, error) {
	if configPath == "" {
		configPath = DefaultRedisConfigPath
	}

	// 检查文件是否存在（必须存在）
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.Errorf("Redis配置文件不存在: %s，请确保文件存在", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, errors.Wrapf(err, "读取Redis配置文件失败: %s", configPath)
	}

	var redisConfig struct {
		Host        string `json:"host"`
		Port        string `json:"port"`
		Password    string `json:"password"`
		Database    string `json:"database"`
		Timeout     int    `json:"timeout"`
		DialTimeout int    `json:"dialTimeout"`
	}

	if err := json.Unmarshal(data, &redisConfig); err != nil {
		return nil, errors.Wrapf(err, "解析Redis配置文件失败: %s", configPath)
	}

	// 解析数据库编号
	var db int
	if redisConfig.Database != "" {
		// 解析失败时 db 保持零值默认（Redis DB 0），和字段本就为空时的行为一致。
		_, _ = fmt.Sscanf(redisConfig.Database, "%d", &db)
	}

	address := fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port)

	// 设置默认值
	timeout := redisConfig.Timeout
	if timeout == 0 {
		timeout = 5
	}
	dialTimeout := redisConfig.DialTimeout
	if dialTimeout == 0 {
		dialTimeout = 5
	}

	return &RedisConf{
		Address:     address,
		Password:    redisConfig.Password,
		DB:          db,
		Timeout:     timeout,
		DialTimeout: dialTimeout,
	}, nil
}

// LoadMiddlewareConfig 从文件加载中间件配置（限流等）
func LoadMiddlewareConfig(configPath string) (*RateLimitConf, error) {
	if configPath == "" {
		configPath = DefaultMiddlewareConfigPath
	}

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.Errorf("中间件配置文件不存在: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, errors.Wrapf(err, "读取中间件配置文件失败: %s", configPath)
	}

	var middlewareConfig struct {
		RateLimit RateLimitConf `json:"rateLimit" yaml:"rateLimit"`
	}

	if err := yaml.Unmarshal(data, &middlewareConfig); err != nil {
		return nil, errors.Wrapf(err, "解析中间件配置文件失败: %s", configPath)
	}

	return &middlewareConfig.RateLimit, nil
}

// MergeExternalConfig 合并外部配置到主配置。gateway 拆分后不再直连任何 MySQL，只需要
// 合并 Redis（必须存在）+ 中间件配置（可选）。
func MergeExternalConfig(c *Config, redisConfigPath, middlewareConfigPath string) error {
	// 从外部文件加载 Redis 配置（必须存在）
	if redisConfigPath == "" {
		redisConfigPath = DefaultRedisConfigPath
	}
	redisConf, err := LoadRedisConfig(redisConfigPath)
	if err != nil {
		return errors.Wrap(err, "加载Redis配置失败，配置文件必须存在")
	}

	// 直接使用外部文件的配置（包含所有Redis相关参数）
	c.Redis = *redisConf

	// 从外部文件加载中间件配置（如果存在则使用，不存在则使用配置文件中的）
	if middlewareConfigPath == "" {
		middlewareConfigPath = DefaultMiddlewareConfigPath
	}
	middlewareConf, err := LoadMiddlewareConfig(middlewareConfigPath)
	if err != nil {
		// 中间件配置可选，如果不存在则使用配置文件中的
		if !os.IsNotExist(err) {
			return errors.Wrap(err, "加载中间件配置失败")
		}
		// 使用配置文件中的中间件配置
	} else {
		// 使用外部文件的配置
		c.RateLimit = *middlewareConf
	}

	return nil
}

// GetConfigPath 获取配置文件路径（支持相对路径和绝对路径）
func GetConfigPath(configFile string) string {
	if filepath.IsAbs(configFile) {
		return configFile
	}
	// 如果是相对路径，尝试从当前工作目录或可执行文件目录查找
	if _, err := os.Stat(configFile); err == nil {
		return configFile
	}
	// 尝试从可执行文件所在目录查找
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		absPath := filepath.Join(exeDir, configFile)
		if _, err := os.Stat(absPath); err == nil {
			return absPath
		}
	}
	return configFile
}
