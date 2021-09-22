package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	hydraSDK "github.com/ory/hydra-client-go/client"
	kratosSDK "github.com/ory/kratos-client-go"
)

type Config struct {
	Env    string       `json:"env"`
	Kratos KratosConfig `json:"kratos"`
	Hydra  HydraConfig  `json:"hydra"`
}

type KratosConfig struct {
	Host           string `json:"host"`
	Scheme         string `json:"scheme"`
	Debug          bool   `json:"debug"`
	PublicBasePath string `json:"public_base_path"`
}

type HydraConfig struct {
	Public HydraTransportConfig
	Admin  HydraTransportConfig
}

type HydraTransportConfig struct {
	Host     string   `json:"host"`
	BasePath string   `json:"base_path"`
	Schemes  []string `json:"schemes"`
}

func DefaultConfig() (kratosSDK.Configuration, hydraSDK.TransportConfig, hydraSDK.TransportConfig) {
	return kratosSDK.Configuration{
			Host:   "oathkeeper:4455",
			Scheme: "http",
			Debug:  true,
			Servers: []kratosSDK.ServerConfiguration{
				{
					URL: "/.ory/kratos/public",
				},
			},
		},
		hydraSDK.TransportConfig{
			Host:     "hydra:4444",
			BasePath: "/",
			Schemes:  []string{"http"},
		},
		hydraSDK.TransportConfig{
			Host:     "hydra:4445",
			BasePath: "/",
			Schemes:  []string{"http"},
		}
}

func (c Config) IsProd() bool {
	env := strings.ToLower(c.Env)
	if env == "prod" || env == "production" {
		return true
	}
	return false
}

func LoadConfig(configReq bool) (kratosSDK.Configuration, hydraSDK.TransportConfig, hydraSDK.TransportConfig) {
	if !configReq {
		return DefaultConfig()
	}
	f, err := os.Open("config.json")
	if err != nil {
		log.Println(err)
		panic(err)
	}
	var c Config
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&c)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println("Successfully loaded config.json file")
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
