package kratos

import (
	"context"
	"github.com/monkjunior/poc-kratos-hydra/models"
	kratosClient "github.com/ory/kratos-client-go"
)

var (
	cfgKratos = kratosClient.Configuration{
		Host:   "127.0.0.1:4455",
		Scheme: "http",
		DefaultHeader: map[string]string{
			"Env": "dev",
		},
		UserAgent: "monk_junior",
		Debug:     true,
		Servers: []kratosClient.ServerConfiguration{
			kratosClient.ServerConfiguration{
				URL: "/.ory/kratos/public",
			},
		},
	}
)

type ClientService interface {
	InitSelfServiceRegistrationFlow() error

	CreateIdentity(identity *models.Identity) error
}

type Client struct {
	API *kratosClient.APIClient
}

func NewClient() *Client {
	return &Client{
		API: kratosClient.NewAPIClient(&cfgKratos),
	}
}

func (k *Client) InitSelfServiceRegistrationFlow() error {
	return nil
}

func (k *Client) CreateIdentity(identity *models.Identity) error {
	req := k.API.V0alpha1Api.AdminCreateIdentity(context.Background())
	req.AdminCreateIdentityBody(kratosClient.AdminCreateIdentityBody{
		SchemaId: "https://schemas.ory.sh/presets/kratos/quickstart/email-password/identity.schema.json",
		Traits:   map[string]interface{}{},
	})
	return nil
}
