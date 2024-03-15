package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Jeffail/gabs/v2"
	"github.com/dgrijalva/jwt-go"
)

func extractUser(tokenString string) (user string, err error) {
	// claims := jwt.MapClaims{}
	var token *jwt.Token
	token, _ = jwt.Parse(tokenString, nil)

	// _, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
	// 	kid, ok := token.Header["kid"].(string)
	// 	fmt.Println(kid)
	// 	if !ok {
	// 		return nil, errors.New("expecting JWT header to have string kid")
	// 	}

	// 	if keys := config.Set.LookupKeyID(kid); len(keys) == 1 {
	// 		var pKey interface{}
	// 		if err := keys[0].Raw(&pKey); err != nil {
	// 			return nil, errors.New("failed to create public key")
	// 		}
	// 		return pKey, nil
	// 	}
	// 	return nil, errors.New("kid for token didn't match any keys issued by issuer")
	// })
	// if err != nil {
	// 	return
	// }
	// do something with decoded claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		for key, val := range claims {
			if key == "upn" {
				emailParts := strings.Split(val.(string), "@")
				user = emailParts[0]
				return
			}
		}
	}
	return
}

func getToken(authCode string) (token string, err error) {
	method := "POST"

	payload := strings.NewReader("client_id=" + config.OIDCClientID + "&client_secret=" + config.OIDCSecret + "&scope=openid&code=" + authCode + "&grant_type=authorization_code&redirect_uri=" + url.QueryEscape(config.ExternalURL))

	client := &http.Client{}
	var req *http.Request
	req, err = http.NewRequest(method, config.OIDCTokenURL, payload)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	var tokenResp *gabs.Container
	tokenResp, err = gabs.ParseJSON(body)
	if err != nil {
		return
	}

	if tokenResp.ExistsP("access_token") {
		token = tokenResp.Path("access_token").Data().(string)
	}
	return
}
