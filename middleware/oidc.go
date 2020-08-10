package middleware

import (
	"context"
	"encoding/json"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"strings"
)

func OIDC(oauth2Config *oauth2.Config, ctx context.Context, verifier *oidc.IDTokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			rawAccessToken := r.Header.Get("Authorization")
			if rawAccessToken == "" {
				http.Redirect(w, r, oauth2Config.AuthCodeURL("some_TODO_state"), http.StatusFound)
				return
			}

			parts := strings.Split(rawAccessToken, " ")
			if len(parts) != 2 {
				w.WriteHeader(400)
				return
			}
			_, err := verifier.Verify(ctx, parts[1])

			if err != nil {
				http.Redirect(w, r, oauth2Config.AuthCodeURL("some_TODO_state"), http.StatusFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func OIDCMiddleware(oauth2Config *oauth2.Config, ctx context.Context, verifier *oidc.IDTokenVerifier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

			log.Println("##### Data: #####")
			log.Println(data)
			log.Println("##### ##### #####")
			w.Write(data)
		})
	}
}
