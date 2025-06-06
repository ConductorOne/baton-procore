package client

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"golang.org/x/oauth2/clientcredentials"
)

type Client struct {
	*uhttp.BaseHttpClient
}

func New(ctx context.Context, clientId, clientSecret string) (*Client, error) {
	config := &clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     "https://login.procore.com/oauth/token",
	}

	client, err := uhttp.NewBaseHttpClientWithContext(ctx, config.Client(ctx))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP client: %w", err)
	}

	return &Client{
		client,
	}, nil
}
