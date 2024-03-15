package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Jeffail/gabs/v2"
	"github.com/lestrrat-go/jwx/jwk"
)

func fetchOIDCConfig() (err error) {
	uri := config.OIDCIssuer + "/.well-known/openid-configuration"
	var resp *http.Response
	resp, err = http.Get(uri)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}
	var oidcConfigParsed *gabs.Container
	oidcConfigParsed, err = gabs.ParseJSON(body)
	if err != nil {
		log.Printf("ERROR: %v", err)
		return
	}
	config.OIDCTokenURL = oidcConfigParsed.Path("token_endpoint").Data().(string)
	config.OIDCAuthURL = oidcConfigParsed.Path("authorization_endpoint").Data().(string)
	set, err := jwk.FetchHTTP(oidcConfigParsed.Path("jwks_uri").Data().(string))
	if err != nil {
		return
	}
	config.Set = *set
	return
}

// func fetchKeys(uri string) (err error) {
// 	var resp *http.Response
// 	resp, err = http.Get(uri)
// 	if err != nil {
// 		log.Printf("ERROR: %v", err)
// 		return
// 	}
// 	var body []byte
// 	body, err = ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Printf("ERROR: %v", err)
// 		return
// 	}

// 	var keys jwksResponse
// 	err = json.Unmarshal(body, &keys)
// 	if err != nil {
// 		log.Printf("ERROR: %v", err)
// 	}
// 	config.Keys = keys.Keys
// 	return
// }
