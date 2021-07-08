package nacos_config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/zhangsq-ax/nacos-helper-go"
	"github.com/zhangsq-ax/nacos-helper-go/options"
	"gopkg.in/yaml.v2"
)

// NacosConfig Nacos 配置对象
type NacosConfig struct {
	configClient *config_client.IConfigClient
	dataId       string
	group        string
}

// NewNacosConfig 创建 NacosConfig 对象
func NewNacosConfig(opts *options.NacosOptions, dataId string, group string) (*NacosConfig, error) {
	configClient, err := nacos_helper.GetConfigClient(opts)
	if err != nil {
		return nil, err
	}

	return &NacosConfig{
		configClient: configClient,
		dataId:       dataId,
		group:        group,
	}, nil
}

// GetConfigString 以字符串形式获取配置
func (nc *NacosConfig) GetConfigString() (string, error) {
	return (*nc.configClient).GetConfig(vo.ConfigParam{
		DataId: nc.dataId,
		Group:  nc.group,
	})
}

// GetConfigJSON 以 JSON 格式获取配置
func (nc *NacosConfig) GetConfigJSON(target interface{}) error {
	content, err := nc.GetConfigString()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(content), target)
	if err != nil {
		return errors.New(fmt.Sprintf("the config information is not valid JSON data: %s", content))
	}
	return nil
}

// GetConfigYAML 以 YAML 格式获取配置
func (nc *NacosConfig) GetConfigYAML(target interface{}) error {
	content, err := nc.GetConfigString()
	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(content), target)
	if err != nil {
		return errors.New(fmt.Sprintf("the config information is not valid YAML data: %s", content))
	}
	return nil
}
