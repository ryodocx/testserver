package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// GET /
func handler(w http.ResponseWriter, req *http.Request) {
	if accessLog {
		log.Printf("%s %s %s", req.RemoteAddr, req.Method, req.RequestURI)
	}
	time.Sleep(responseSleep)
	_, _ = w.Write(responseBody)
}

// GET /echo
func echoHandler(w http.ResponseWriter, req *http.Request) {
	if accessLog {
		log.Printf("%s %s %s", req.RemoteAddr, req.Method, req.RequestURI)
	}

	// Parse Authorization header
	var authzInfo = map[string]any{}
	{
		authzHeader := req.Header.Get("Authorization")
		if len(authzHeader) <= 7 {
			// empty or not have bearer token
			goto breakAuthz
		}

		// JWT
		jwtInfo := map[string]any{}
		jwtAll := strings.Split(strings.TrimPrefix(authzHeader, "Bearer "), ".")
		if strings.HasPrefix(authzHeader, "Bearer ") && len(jwtAll) == 3 {
			parseJWT := func(jwtBodyB64Encoded, mapKey string) error {
				jwtBody := map[string]any{}
				jwtBodyDecoded, err := base64.RawURLEncoding.DecodeString(jwtBodyB64Encoded)
				if err != nil {
					return err
				}
				if err := json.Unmarshal(jwtBodyDecoded, &jwtBody); err != nil {
					return err
				}
				jwtInfo[mapKey] = jwtBody
				return nil
			}

			if err := parseJWT(jwtAll[0], "header"); err != nil {
				log.Println(err)
				goto breakAuthz
			}
			if err := parseJWT(jwtAll[1], "payload"); err != nil {
				log.Println(err)
				goto breakAuthz
			}

			authzInfo["jwt"] = jwtInfo
		}
	}
breakAuthz:

	respMap := map[string]interface{}{
		"Header":        req.Header,
		"Form":          req.Form,
		"Proto":         req.Proto,
		"Method":        req.Method,
		"Host":          req.Host,
		"RequestURI":    req.RequestURI,
		"RemoteAddr":    req.RemoteAddr,
		"Authorization": authzInfo,
	}

	resp, err := json.MarshalIndent(respMap, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(resp)
}
