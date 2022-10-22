package common

import (
	"context"
	"net/http"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

func GetHTTPClient() *http.Client {
	c, err := idtoken.NewClient(context.Background(), "")
	if err != nil {
		panic(err)
	}
	return c
}

// MakeAuthenticatedHTTPClient blah blah.
func MakeAuthenticatedHTTPClient() (*http.Client, error) {
	_, err := google.FindDefaultCredentials(context.Background())
	if err != nil {
		panic(err)
	}

	return nil, nil
}
