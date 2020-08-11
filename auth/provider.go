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
	Enabled      bool
	Issuer       string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

type OauthProvider struct {
	enabled      bool
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	context      context.Context
	logger       hclog.Logger
}

const oauthState = "myState"

func NewProvider(config OauthConfig, logger hclog.Logger) (*OauthProvider, error) {
	if !config.Enabled {
		return &OauthProvider{enabled: false}, nil
	}

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
	if !p.enabled {
		// return no-op Middleware
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(rw, r)
		})
	}

	oauth2Config := p.oauth2Config
	verifier := p.verifier
	ctx := p.context

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		rawAccessToken := r.Header.Get("Authorization")
		if rawAccessToken == "" {
			http.Redirect(rw, r, oauth2Config.AuthCodeURL(oauthState), http.StatusFound)
			return
		}

		parts := strings.Split(rawAccessToken, " ")
		if len(parts) != 2 {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err := verifier.Verify(ctx, parts[1])

		if err != nil {
			http.Redirect(rw, r, oauth2Config.AuthCodeURL(oauthState), http.StatusFound)
			return
		}

		next.ServeHTTP(rw, r)
	})
}

func (p *OauthProvider) CallbackHandler() http.Handler {
	if !p.enabled {
		// Return unimplemented callback handler
		return http.HandlerFunc(func(rw http.ResponseWriter, request *http.Request) {
			http.Error(rw, "OAuth2 callback disabled!", http.StatusNotImplemented)
		})
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		oauth2Config := p.oauth2Config
		verifier := p.verifier
		ctx := p.context

		oauth2Token, err := oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(rw, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			http.Error(rw, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			return
		}
		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			http.Error(rw, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resp := struct {
			OAuth2Token   *oauth2.Token
			IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
		}{oauth2Token, new(json.RawMessage)}

		if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := json.MarshalIndent(resp, "", "    ")
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		p.logger.Debug("Received OAuth2 Token: ", "token", hclog.Fmt("%s", string(data)))
		rw.Write(data)
	})
}
