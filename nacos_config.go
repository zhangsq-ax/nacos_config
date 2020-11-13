package nacos_config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/nacos_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/http_agent"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v2"
)

type NacosConfig struct {
	configClient *config_client.ConfigClient
	dataId       string
	group        string
}

type NacosOptions struct {
	Host        string `json:"host"`
	Port        uint64 `json:"port"`
	NamespaceId string `json:"namespace_id"`
	Group       string `json:"group"`
	DataId      string `json:"data_id"`
}

func NewNacosConfig(opts NacosOptions) (*NacosConfig, error) {
	clientConfig := constant.ClientConfig{
		NamespaceId:         opts.NamespaceId,
		TimeoutMs:           10 * 1000,
		BeatInterval:        5 * 1000,
		ListenInterval:      300 * 1000,
		NotLoadCacheAtStart: true,
	}

	serverConfig := constant.ServerConfig{
		IpAddr:      opts.Host,
		Port:        opts.Port,
		ContextPath: "/nacos",
	}

	nc := nacos_client.NacosClient{}
	err := nc.SetServerConfig([]constant.ServerConfig{serverConfig})
	if err != nil {
		return nil, err
	}
	err = nc.SetClientConfig(clientConfig)
	if err != nil {
		return nil, err
	}
	err = nc.SetHttpAgent(&http_agent.HttpAgent{})
	if err != nil {
		return nil, err
	}

	client, err := config_client.NewConfigClient(&nc)
	if err != nil {
		return nil, err
	}

	return &NacosConfig{
		configClient: &client,
		dataId:       opts.DataId,
		group:        opts.Group,
	}, nil
}

func (nc *NacosConfig) GetConfigString() (string, error) {
	return nc.configClient.GetConfig(vo.ConfigParam{
		DataId: nc.dataId,
		Group:  nc.group,
	})
}

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
