package cmd

import (
	"context"
	"crypto/tls"
	"math"
	"net/http"
	"time"

	"github.com/monkjunior/poc-kratos-hydra/pkg/config"
	"github.com/monkjunior/poc-kratos-hydra/pkg/log"
	hydraSDK "github.com/ory/hydra-client-go/client"
	hydraAdmin "github.com/ory/hydra-client-go/client/admin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// createClientCmd represents the createClient command
var createClientCmd = &cobra.Command{
	Use:   "create-client",
	Short: "Register an oauth2 client in Authorization Server (Hydra).",
	Run:   runCreateClientCmd,
}

func init() {
	oauth2Cmd.AddCommand(createClientCmd)
}

func runCreateClientCmd(cmd *cobra.Command, args []string) {
	logger := log.GetLogger().With(zap.String("cmd", "oauth2 create-client"))

	_, _, hAdmCfg := config.GetAuthStackCfg()
	hAdm := hydraSDK.NewHTTPClientWithConfig(nil, &hAdmCfg)

	for i := 0; i < 6; i++ {
		isAlive, err := hAdm.Admin.IsInstanceAlive(&hydraAdmin.IsInstanceAliveParams{
			Context: context.Background(),
		})
		waitTime := math.Pow(2, float64(i))
		if err != nil || isAlive == nil {
			if waitTime > 30 {
				logger.Fatal("timeout, hydra is not ready")
			}
			time.Sleep(time.Duration(waitTime) * time.Second)
			continue
		}
		logger.Info("hydra is ready, starting registration an OAuth2 client",
			zap.Float64("waited_seconds", waitTime))
		break
	}

	params := hydraAdmin.CreateOAuth2ClientParams{
		Body:    config.Cfg.GetHydraOauth2Config(),
		Context: context.Background(),
		HTTPClient: &http.Client{
			Transport: newTransport(oauth2Cmd),
		},
	}
	result, err := hAdm.Admin.CreateOAuth2Client(&params)
	if err != nil || result == nil {
		logger.Fatal("failed to register a new oauth2 client", zap.Error(err))
	}
	logger.Info("successfully register a new oauth2 client")
}

type transport struct {
	Transport http.RoundTripper
	cmd       *cobra.Command
}

func newTransport(cmd *cobra.Command) *transport {
	skipTLSVerify, err := cmd.PersistentFlags().GetBool("skip-tls-verify")
	if err != nil {
		skipTLSVerify = false
	}
	return &transport{
		cmd: cmd,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLSVerify},
		},
	}
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	fakeTLSTerm, err := t.cmd.PersistentFlags().GetBool("fake-tls-termination")
	if err != nil {
		fakeTLSTerm = false
	}

	if fakeTLSTerm {
		req.Header.Set("X-Forwarded-Proto", "https")
	}
	return t.Transport.RoundTrip(req)
}
