package authenticator

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

var httpClient = http.Client{
	Timeout: time.Second * 5,
}

type OpenIdConfiguration struct {
	JwksUri string `json:"jwks_uri"`
	Issuer  string `json:"issuer"`
}

func fetchOpenIDConfiguration(url string) (*OpenIdConfiguration, error) {
	res, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	c := OpenIdConfiguration{}
	jsonErr := json.Unmarshal(body, &c)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &c, nil
}
