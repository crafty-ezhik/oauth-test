package utils

import (
	"log/slog"
	"net/url"
	"os"
)

func GenerateGoogleAuthURL() (string, error) {
	authUrl, err := url.Parse("https://accounts.google.com/o/oauth2/v2/auth")
	if err != nil {
		slog.Info("Error parsing auth url")
		return "", err
	}
	q := authUrl.Query()
	q.Set("client_id", os.Getenv("OAUTH_GOOGLE_CLIENT_ID"))
	q.Set("redirect_uri", os.Getenv("OAUTH_GOOGLE_REDIRECT_URI"))
	q.Set("response_type", "code")
	q.Set("scope", "openid profile email https://www.googleapis.com/auth/drive https://www.googleapis.com/auth/calendar")
	q.Set("state", os.Getenv("OAUTH_GOOGLE_STATE"))

	authUrl.RawQuery = q.Encode()
	return authUrl.String(), nil
}
