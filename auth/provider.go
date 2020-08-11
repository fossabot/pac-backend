package auth

import (
	"context"
	"encoding/json"
	"github.com/coreos/go-oidc"
	"github.com/hashicorp/go-hclog"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
)

type OauthConfig struct {
	Issuer       string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

type OauthProvider struct {
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	context      context.Context
	logger       hclog.Logger
}

const oauthState = "myState"

func NewProvider(config OauthConfig, logger hclog.Logger) (*OauthProvider, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, config.Issuer)
	if err != nil {
		return nil, err
	}

	// Create OAuth2 configuration
	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		// Discovery returns the OAuth2 endpoints.
		Endpoint:    provider.Endpoint(),
		RedirectURL: config.RedirectURL,
		// "openid" is a required scope for OpenID Connect flows.
		Scopes: config.Scopes,
	}

	oidcConfig := &oidc.Config{
		ClientID: config.ClientID,
	}

	// Verifier to verify JWTs
	verifier := provider.Verifier(oidcConfig)

	return &OauthProvider{
		oauth2Config: oauth2Config,
		verifier:     verifier,
		context:      ctx,
		logger:       logger,
	}, nil
}

func (p *OauthProvider) Middleware(next http.Handler) http.Handler {
	oauth2Config := p.oauth2Config
	verifier := p.verifier
	ctx := p.context

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rawAccessToken := r.Header.Get("Authorization")
		if rawAccessToken == "" {
			http.Redirect(w, r, oauth2Config.AuthCodeURL(oauthState), http.StatusFound)
			return
		}

		parts := strings.Split(rawAccessToken, " ")
		if len(parts) != 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err := verifier.Verify(ctx, parts[1])

		if err != nil {
			http.Redirect(w, r, oauth2Config.AuthCodeURL(oauthState), http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (p *OauthProvider) CallbackHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oauth2Config := p.oauth2Config
		verifier := p.verifier
		ctx := p.context

		oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			return
		}
		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp := struct {
			OAuth2Token   *oauth2.Token
			IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
		}{oauth2Token, new(json.RawMessage)}

		if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := json.MarshalIndent(resp, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		p.logger.Debug("Received OAuth2 Token: ", "token", hclog.Fmt("%s", string(data)))
		w.Write(data)
	})
}
