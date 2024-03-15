package main

import "github.com/lestrrat-go/jwx/jwk"

type configOptions struct {
	OIDCIssuer    string   `json:"oidcIssuer"`
	OIDCClientID  string   `json:"oidcClientId"`
	OIDCSecret    string   `json:"oidcSecret"`
	OIDCAuthURL   string   `json:"oidcAuthUrl"`
	OIDCTokenURL  string   `json:"oidcTokenUrl"`
	ExternalURL   string   `json:"externalUrl"`
	ListenAddress string   `json:"listenAddress"`
	TargetURL     string   `json:"targetUrl"`
	ExclusionList []string `json:"exclusionList"`
	LogLevel      string   `json:"LogLevel"`
	Set           jwk.Set
}

// type jwksResponse struct {
// 	Keys []rawJSONWebKey `json:"keys"`
// }

// type rawJSONWebKey struct {
// 	Use string `json:"use,omitempty"`
// 	Kty string `json:"kty,omitempty"`
// 	Kid string `json:"kid,omitempty"`
// 	Crv string `json:"crv,omitempty"`
// 	Alg string `json:"alg,omitempty"`
// 	// Certificates
// 	X5c       []string `json:"x5c,omitempty"`
// 	X5u       *url.URL `json:"x5u,omitempty"`
// 	X5tSHA1   string   `json:"x5t,omitempty"`
// 	X5tSHA256 string   `json:"x5t#S256,omitempty"`
// }
