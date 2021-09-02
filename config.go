package main

import (
	hydraSDK "github.com/ory/hydra-client-go/client"
	kratosSDK "github.com/ory/kratos-client-go"
)

var (
	CfgKratos = kratosSDK.Configuration{
		Host:   "oathkeeper:4455",
		Scheme: "http",
		Debug:  true,
		Servers: []kratosSDK.ServerConfiguration{
			{
				URL: "/.ory/kratos/public",
			},
		},
	}
	ConfigHydraClient = hydraSDK.TransportConfig{
		Host:     "hydra:4444",
		BasePath: "/",
		Schemes:  []string{"http"},
	}
	ConfigHydraAdmin = hydraSDK.TransportConfig{
		Host:     "hydra:4445",
		BasePath: "/",
		Schemes:  []string{"http"},
	}
)
