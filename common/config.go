package common

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	hydraSDK "github.com/ory/hydra-client-go/client"
	kratosSDK "github.com/ory/kratos-client-go"
)

type Config struct {
	Env     string       `json:"env"`
	BaseURL string       `json:"base_url"`
	Kratos  KratosConfig `json:"kratos"`
	Hydra   HydraConfig  `json:"hydra"`
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

// DefaultConfig for develop environment
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

// IsProd read the flag -prod value
// when -prod is set we must read config from file at /etc/config/config.json
func (c Config) IsProd() bool {
	env := strings.ToLower(c.Env)
	if env == "prod" || env == "production" {
		return true
	}
	return false
}

// GetAuthStackCfg read config file at /etc/config/config.json
// then return clients to interact Kratos public, Hydra public, Hydra admin.
func GetAuthStackCfg(configReq bool) (kratosSDK.Configuration, hydraSDK.TransportConfig, hydraSDK.TransportConfig) {
	if !configReq {
		log.Println("Using default config")
		return DefaultConfig()
	}
	c := parseConfig()
	log.Println("Successfully loaded /etc/config/config.json file")
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

// GetKratosPublicBaseURL read config file at /etc/config/config.json
// then return public base URL of Kratos.
func GetKratosPublicBaseURL() string {
	c := parseConfig()
	return c.BaseURL + c.Kratos.PublicBasePath
}

// parseConfig will read the config file at /etc/config/config.json
// and return the appropriate Config.
func parseConfig() Config {
	f, err := os.Open("/etc/config/config.json")
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
	return c
}
