package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
)

const (
	FacebookInfoURL = "https://graph.facebook.com/v10.0/me?fields=id%2Cname%2Cemail"
	GoogleInfoURL   = "https://openidconnect.googleapis.com/v1/userinfo"
)

type OAuth struct {
	Conf    *oauth2.Config
	InfoURL string
}

func (a *OAuth) GetAuthURL(state string) string {
	url := a.Conf.AuthCodeURL(state)
	return url
}

func (a *OAuth) GetUserData(ctx context.Context, code string) (map[string]interface{}, error) {
	token, err := a.Conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("error: could not exchange code for token: %w", err)
	}
	client := a.Conf.Client(ctx, token)
	resp, err := client.Get(a.InfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// sub(id), name, email
	data := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	// fix for data from google, where "sub" field is used instead of "id"
	if _, ok := data["id"]; !ok {
		data["id"] = data["sub"]
		delete(data, "sub")
	}
	return data, nil
}
