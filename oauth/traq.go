package oauth

import (
	"context"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

const (
	clientID = "nmVeJT08KHXIdB8xlrCIwa6YJTkISrP5zWzm"
)

var (
	authURL, _  = url.Parse("https://q.trap.jp/api/v3/oauth2/authorize")
	tokenURL, _ = url.Parse("https://q.trap.jp/api/v3/oauth2/token")
)

type OauthClient struct {
	conf   oauth2.Config
	client *http.Client
}

func New() *OauthClient {
	conf := oauth2.Config{
		ClientID: clientID,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL.String(),
			TokenURL: tokenURL.String(),
		},
	}
	return &OauthClient{
		conf: conf,
	}
}

func (c *OauthClient) GetToken(code string, codeVerifier string) (*oauth2.Token, error) {
	tok, err := c.conf.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", codeVerifier))
	if err != nil {
		return nil, err
	}
	return tok, nil
}
