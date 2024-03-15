package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var config configOptions

func main() {
	var configBytes []byte
	var err error
	if len(os.Args) > 1 {
		configBytes, err = ioutil.ReadFile(os.Args[1])
	} else {
		configBytes, err = ioutil.ReadFile("config.json")
	}

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		serveReverseProxy(res, req)
	})
	err = fetchOIDCConfig()
	if err != nil {
		panic(errors.New("Unable to fetch config from OIDC Issuer"))
	}
	log.Fatal(http.ListenAndServeTLS(config.ListenAddress, "server.crt", "server.key", nil))
}

// Serve a reverse proxy for a given url
func serveReverseProxy(res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(config.TargetURL)
	if strings.ToUpper(config.LogLevel) == "DEBUG" {
		log.Printf("Got request for Path %s", req.URL.Path)
	}
	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)
	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Host = url.Host
	req.Header.Set("Origin", config.TargetURL)
	req.Header.Set("Referer", config.TargetURL+"/")
	if !inExclusionList(req.URL.Path) {
		authFlow(res, req)
		// Note that ServeHttp is non blocking and uses a go routine under the hood
	}
	proxy.ModifyResponse = rewriteBody
	proxy.ServeHTTP(res, req)
}

func inExclusionList(activePath string) bool {
	for _, pathPattern := range config.ExclusionList {
		var path = regexp.MustCompile(pathPattern)
		if path.MatchString(activePath) {
			if strings.ToUpper(config.LogLevel) == "DEBUG" {
				log.Printf("By-passing proxy for pattern match %s", pathPattern)
			}
			return true
		}
	}
	return false
}

func authFlow(res http.ResponseWriter, req *http.Request) {
	authCode := req.URL.Query().Get("code")
	if authCode != "" {
		token, err := getToken(authCode)
		if err == nil {
			res.Header().Add("Set-Cookie", "access_token="+token+";path=/")
			originalPath, err := req.Cookie("original_path")
			if err == nil && originalPath.Value != "" {
				res.Header().Add("Set-Cookie", "original_path=")
				http.Redirect(res, req, originalPath.Value, 303)
			} else {
				http.Redirect(res, req, config.ExternalURL, 303)
			}
		}
	} else {
		tokenString, err := req.Cookie("access_token")
		if err != nil {
			originalPath, err := req.Cookie("original_path")
			if err != nil || originalPath.Value == "" {
				capturedPath := config.ExternalURL + req.URL.RequestURI()
				res.Header().Set("Set-Cookie", "original_path="+capturedPath+";path=/")
				http.Redirect(res, req, "/", 303)
			}
			http.Redirect(res, req, config.OIDCAuthURL+"?response_type=code&scope=openid&client_id="+config.OIDCClientID+"&redirect_uri="+config.ExternalURL, 303)
		} else {
			var user string
			user, err = extractUser(tokenString.Value)
			if err != nil {
				panic(err)
			}
			req.Header.Add("X-Alfresco-Remote-User", user)
		}
	}
}

func rewriteBody(resp *http.Response) (err error) {
	loc := resp.Header.Get("Location")
	if strings.Contains(loc, config.TargetURL) {
		loc = strings.ReplaceAll(loc, config.TargetURL, "")
		loc = strings.Join([]string{config.ExternalURL, loc}, "")
		resp.Header.Set("Location", loc)
	}
	return nil
}

// func refererForURL(lastReq, newReq *url.URL) string {
// 	// https://tools.ietf.org/html/rfc7231#section-5.5.2
// 	//   "Clients SHOULD NOT include a Referer header field in a
// 	//    (non-secure) HTTP request if the referring page was
// 	//    transferred with a secure protocol."
// 	if lastReq.Scheme == "https" && newReq.Scheme == "http" {
// 		return ""
// 	}
// 	return TargetURL
// }
