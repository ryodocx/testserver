package main

import (
	"encoding/json"
	"log"
	"net/http"
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
	respMap := map[string]interface{}{
		"Header":     req.Header,
		"Form":       req.Form,
		"Proto":      req.Proto,
		"Method":     req.Method,
		"Host":       req.Host,
		"RequestURI": req.RequestURI,
		"RemoteAddr": req.RemoteAddr,
	}
	resp, err := json.MarshalIndent(respMap, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(resp)
}
