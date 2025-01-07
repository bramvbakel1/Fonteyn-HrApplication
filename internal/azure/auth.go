package azure

import (
	"context"
	"fmt"

	"golang.org/x/oauth2/clientcredentials"
)

func GetAccessToken(tenantID, clientID, clientSecret string) (string, error) {
	conf := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     "https://login.microsoftonline.com/" + tenantID + "/oauth2/v2.0/token",
		Scopes:       []string{"https://graph.microsoft.com/.default"},
	}

	token, err := conf.Token(context.Background())
	if err != nil {
		return "", fmt.Errorf("unable to get token: %v", err)
	}

	return token.AccessToken, nil
}
