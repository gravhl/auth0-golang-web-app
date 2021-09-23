package auth

import (
	"context"
	"log"

	"app"

	"golang.org/x/oauth2"

	oidc "github.com/coreos/go-oidc"
)

type Authenticator struct {
	Provider *oidc.Provider
	Config   oauth2.Config
	Ctx      context.Context
}

func NewAuthenticator() (*Authenticator, error) {
	ctx := context.Background()

	// check that we have the AUTH0_DOMAIN var

	provider, err := oidc.NewProvider(ctx, "https://"+app.Auth0Domain+"/")
	if err != nil {
		app.Log.Errorf("Failed to get provider: %v", err)
		log.Printf("failed to get provider: %v", err)
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     app.Auth0ClientID,
		ClientSecret: app.Auth0ClientSecret,
		RedirectURL:  app.Auth0CallbackURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile"},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Ctx:      ctx,
	}, nil
}
