package pkg

import (
	"net/url"
	"os"
)

func GetEndpoint(path string) string {
	var baseURL string
	baseURL = "http://0.0.0.0:8080/"
	env := os.Getenv("APP_ENV")
	if env == "stage" {
		baseURL = "http://stage.localhost:8080/"
	}
	p, _ := url.Parse(path)
	b, _ := url.Parse(baseURL)
	return b.ResolveReference(p).String()
}
