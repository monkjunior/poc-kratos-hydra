package cmd

import (
	hydraSDK "github.com/ory/hydra-client-go/client"
	kratosSDK "github.com/ory/kratos-client-go"
)

type Config struct {
	Env     string       `json:"env"`
	BaseURL string       `json:"baseUrl"`
	Kratos  KratosConfig `json:"kratos"`
	Hydra   HydraConfig  `json:"hydra"`
}

type KratosConfig struct {
	Host           string `json:"host"`
	Scheme         string `json:"scheme"`
	Debug          bool   `json:"debug"`
	PublicBasePath string `json:"publicBasePath"`
}

type HydraConfig struct {
	Public HydraTransportConfig
	Admin  HydraTransportConfig
}

type HydraTransportConfig struct {
	Host     string   `json:"host"`
	BasePath string   `json:"basePath"`
	Schemes  []string `json:"schemes"`
}

// GetAuthStackCfg read config file at
// then return clients to interact Kratos public, Hydra public, Hydra admin.
func GetAuthStackCfg() (kratosSDK.Configuration, hydraSDK.TransportConfig, hydraSDK.TransportConfig) {
	c := Cfg
	return kratosSDK.Configuration{
			Host:   c.Kratos.Host,
			Scheme: c.Kratos.Scheme,
			Debug:  c.Kratos.Debug,
			Servers: []kratosSDK.ServerConfiguration{
				{
					URL: c.Kratos.PublicBasePath,
				},
			},
		},
		hydraSDK.TransportConfig{
			Host:     c.Hydra.Public.Host,
			BasePath: c.Hydra.Public.BasePath,
			Schemes:  c.Hydra.Public.Schemes,
		},
		hydraSDK.TransportConfig{
			Host:     c.Hydra.Admin.Host,
			BasePath: c.Hydra.Admin.BasePath,
			Schemes:  c.Hydra.Admin.Schemes,
		}
}
