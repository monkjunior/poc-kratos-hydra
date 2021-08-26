package kratos

import (
	"context"
	"net/http"
	"net/http/cookiejar"

	"github.com/monkjunior/poc-kratos-hydra/models"
	kratosClient "github.com/ory/kratos-client-go"
)

var (
	cfgKratos = kratosClient.Configuration{
		Host:   "127.0.0.1:4455",
		Scheme: "http",
		Debug:  true,
		Servers: []kratosClient.ServerConfiguration{
			{
				URL: "/.ory/kratos/public",
			},
		},
	}
)

type ClientService interface {
	InitSelfServiceRegistrationFlow(ctx context.Context) (*kratosClient.SelfServiceRegistrationFlow, *http.Response, error)
	GetRegistrationFlow(ctx context.Context, flowID string, cookie string) (*kratosClient.SelfServiceRegistrationFlow, *http.Response, error)
	CreateIdentity(identity *models.Identity) error
}

type Client struct {
	API *kratosClient.APIClient
}

func NewClient() *Client {
	cj, _ := cookiejar.New(nil)
	cfgKratos.HTTPClient = &http.Client{Jar: cj}
	return &Client{
		API: kratosClient.NewAPIClient(&cfgKratos),
	}
}

func (c *Client) InitSelfServiceRegistrationFlow(ctx context.Context) (*kratosClient.SelfServiceRegistrationFlow, *http.Response, error) {
	flow, res, err := c.API.V0alpha1Api.InitializeSelfServiceRegistrationFlowForBrowsers(ctx).Execute()
	LogOnError(err, res)
	PrintJSONPretty(flow)

	return flow, res, err
}

func (c *Client) CreateIdentity(identity *models.Identity) error {
	req := c.API.V0alpha1Api.AdminCreateIdentity(context.Background())
	req.AdminCreateIdentityBody(kratosClient.AdminCreateIdentityBody{
		SchemaId: "https://schemas.ory.sh/presets/kratos/quickstart/email-password/identity.schema.json",
		Traits:   map[string]interface{}{},
	})
	return nil
}

func (c *Client) GetRegistrationFlow(ctx context.Context, flowID string, cookie string) (*kratosClient.SelfServiceRegistrationFlow, *http.Response, error) {
	return c.API.V0alpha1Api.GetSelfServiceRegistrationFlow(ctx).Cookie(cookie).Id(flowID).Execute()
}
