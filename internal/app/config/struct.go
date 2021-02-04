package config

import (
	"github.com/qianjunakasumi/Tennouji/internal/pkg/logger"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Client 客户端
type Client struct {
	Number   uint64 `yaml:"number"`
	Password string `yaml:"password"`
}

// GetConfig 获取配置
func GetConfig() (conf *Client, err error) {

	f, err := ioutil.ReadFile(".tennouji/config.yml")
	if err != nil {
		logger.Error("无法读取配置文件", zap.Error(err))
		return nil, err
	}

	err = yaml.Unmarshal(f, conf)
	if err != nil {
		logger.Error("无法解析配置文件或损坏的配置文件", zap.Error(err))
		return nil, err
	}

	return
}
