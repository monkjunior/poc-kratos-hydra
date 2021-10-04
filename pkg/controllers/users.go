package controllers

import (
	"errors"
	"go.uber.org/zap"
	"net/http"

	"github.com/monkjunior/poc-kratos-hydra/pkg/common"
	"github.com/monkjunior/poc-kratos-hydra/pkg/config"
	"github.com/monkjunior/poc-kratos-hydra/pkg/log"
	"github.com/monkjunior/poc-kratos-hydra/pkg/views"
	kratosClient "github.com/ory/kratos-client-go"
)

var (
	ErrConfirmPasswordMismatched = errors.New("controller: confirm password mismatched")
)

func NewUsers(k *kratosClient.APIClient) *Users {
	KratosPublicURL = config.Cfg.BaseURL + config.Cfg.Kratos.PublicBasePath
	return &Users{
		LoginView:        views.NewView("bootstrap", "login"),
		RegistrationView: views.NewView("bootstrap", "registration"),
		kratosClient:     k,
	}
}

// Users controller handles traditions authentication flows, includes: registration, login, logout and so on
// It interacts with Ory Kratos, an opensource Identity Provider.
type Users struct {
	LoginView        *views.View
	RegistrationView *views.View
	CallbackView     *views.View
	kratosClient     *kratosClient.APIClient
}

// LoginForm stores data for rendering Login form and submit a Login flow
type LoginForm struct {
	SubmitMethod string
	Action       string
	CsrfToken    string `schema:"csrf_token"`
	FlowID       string
	Email        string `schema:"password_identifier"`
	Password     string `schema:"password"`
}

// GetLogin requires flow params, if the flow is not set, it will redirect to Kratos to browse a new one.
// Kratos will create a new flow and redirect back to /auth/login with the param was set in the URL.
// GetLogin will use this id to fetch data from Kratos to render submit form.
//
// GET /auth/login/?flow=<flow_id>
func (u *Users) GetLogin(w http.ResponseWriter, r *http.Request) {
	// TODO: logging
	flow := r.URL.Query().Get("flow")
	if flow == "" {
		http.Redirect(w, r, KratosPublicURL+KratosSSLoginBrowserPath, http.StatusFound)
		return
	}
	flowObject, res, err := u.kratosClient.V0alpha1Api.GetSelfServiceLoginFlow(r.Context()).Id(flow).Cookie(r.Header.Get("Cookie")).Execute()
	if err != nil || res == nil || res.StatusCode != http.StatusOK {
		// TODO: handle error when received wrong flow id, should create a new flow
		return
	}

	data := views.Data{
		Yield: LoginForm{
			CsrfToken:    flowObject.Ui.GetNodes()[0].Attributes.UiNodeInputAttributes.Value.(string),
			FlowID:       flow,
			SubmitMethod: flowObject.Ui.Method,
			Action:       flowObject.Ui.Action,
		},
	}
	u.LoginView.Render(w, r, data)
}

// RegistrationForm stores data for rendering Registration form and submit a Registration flow
type RegistrationForm struct {
	RegistrationMethod string
	SubmitMethod       string
	Action             string
	CsrfToken          string `schema:"csrf_token"`
	FlowID             string
	Email              string `schema:"traits.email"`
	Password           string `schema:"password"`
}

// GetRegistration requires flow params to render Registration screen
// if flow param is not found, it will redirect to Kratos /self-service/registration/browser
// to browse a new flow_id.
// Kratos then redirect back to this path with a flow param in the URL.
//
// GET /auth/registration/?flow=<flow_id>
func (u *Users) GetRegistration(w http.ResponseWriter, r *http.Request) {
	flow := r.URL.Query().Get("flow")
	if flow == "" {
		http.Redirect(w, r, KratosPublicURL+KratosSSRegistrationBrowserPath, http.StatusFound)
		return
	}
	flowObject, res, err := u.kratosClient.V0alpha1Api.GetSelfServiceRegistrationFlow(r.Context()).Id(flow).Cookie(r.Header.Get("Cookie")).Execute()
	if err != nil || res == nil || res.StatusCode != http.StatusOK {
		common.LogOnError(err, res)
		return
	}
	//kratos.PrintJSONPretty(flowObject)
	data := views.Data{
		Yield: RegistrationForm{
			CsrfToken:    flowObject.Ui.GetNodes()[0].Attributes.UiNodeInputAttributes.Value.(string),
			FlowID:       flow,
			SubmitMethod: flowObject.Ui.Method,
			Action:       flowObject.Ui.Action,
		},
	}
	u.RegistrationView.Render(w, r, data)
}

type ChangePasswordForm struct {
	Token           string `schema:"token"`
	CurrentPassword string `schema:"current_password"`
	NewPassword     string `schema:"new_password"`
	ConfirmPassword string `schema:"confirm_password"`
}

// PostChangePassword handles request from front-end app to change password of current user
func (u *Users) PostChangePassword(w http.ResponseWriter, r *http.Request) {
	uEmail := r.Header.Get("X-User-Email")
	logger := log.GetLogger().With(
		zap.String("receiver", "User"),
		zap.String("method", "PostChangePassword"),
		zap.String("user_email", uEmail),
	)
	var form ChangePasswordForm
	err := parseForm(r, &form)
	if err != nil {
		logger.Error("failed to parse form", zap.Error(err))
		return
	}
	loginFlow, res, err := u.kratosClient.V0alpha1Api.InitializeSelfServiceLoginFlowWithoutBrowser(r.Context()).Execute()
	if err != nil || res == nil || res.StatusCode != http.StatusOK || loginFlow == nil {
		if res != nil {
			logger.With(zap.Int("response_status_code", res.StatusCode))
		}
		logger.Error("failed to init login flow",
			zap.Error(err),
		)
		return
	}
	loginResult, res, err := u.kratosClient.V0alpha1Api.SubmitSelfServiceLoginFlow(r.Context()).Flow(loginFlow.Id).SubmitSelfServiceLoginFlowBody(
		kratosClient.SubmitSelfServiceLoginFlowWithPasswordMethodBodyAsSubmitSelfServiceLoginFlowBody(&kratosClient.SubmitSelfServiceLoginFlowWithPasswordMethodBody{
			Method:             "password",
			Password:           form.CurrentPassword,
			PasswordIdentifier: uEmail,
		}),
	).Execute()
	if err != nil || res == nil || res.StatusCode != http.StatusOK || loginResult == nil {
		if res != nil {
			logger.With(zap.Int("response_status_code", res.StatusCode))
		}
		logger.Error("failed to validate current password", zap.Error(err))
		return
	}
	if form.ConfirmPassword != form.NewPassword {
		logger.Error("confirm password is not the same with new password",
			zap.Error(ErrConfirmPasswordMismatched),
		)
		return
	}
	sessionToken := loginResult.GetSessionToken()
	changePwFlow, res, err := u.kratosClient.V0alpha1Api.InitializeSelfServiceSettingsFlowWithoutBrowserExecute(
		kratosClient.V0alpha1ApiApiInitializeSelfServiceSettingsFlowWithoutBrowserRequest{}.XSessionToken(sessionToken),
	)
	if err != nil || res == nil || res.StatusCode != http.StatusOK || changePwFlow == nil {
		if res != nil {
			logger.With(zap.Int("response_status_code", res.StatusCode))
		}
		logger.Error("failed to init change password flow",
			zap.Error(err),
		)
		return
	}

	changePwResult, res, err := u.kratosClient.V0alpha1Api.SubmitSelfServiceSettingsFlow(r.Context()).Flow(changePwFlow.Id).XSessionToken(sessionToken).SubmitSelfServiceSettingsFlowBody(
		kratosClient.SubmitSelfServiceSettingsFlowWithPasswordMethodBodyAsSubmitSelfServiceSettingsFlowBody(&kratosClient.SubmitSelfServiceSettingsFlowWithPasswordMethodBody{
			Method:   "password",
			Password: form.NewPassword,
		}),
	).Execute()
	if err != nil || res == nil || res.StatusCode != http.StatusOK || changePwResult == nil {
		if res != nil {
			logger.With(zap.Int("response_status_code", res.StatusCode))
		}
		logger.Error("failed to change password password", zap.Error(err))
		return
	}
	logger.Info("change password successfully")
}
