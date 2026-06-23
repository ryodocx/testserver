package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEchoHandler_ValidJWT(t *testing.T) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"1234567890","name":"John Doe","iat":1516239022}`))
	signature := "dummy"

	jwt := header + "." + payload + "." + signature

	req := httptest.NewRequest(http.MethodGet, "/echo", nil)
	req.Header.Set("Authorization", "Bearer "+jwt)

	w := httptest.NewRecorder()
	echoHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", res.Status)
	}

	var respMap map[string]any
	if err := json.NewDecoder(res.Body).Decode(&respMap); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	authzInfo, ok := respMap["Authorization"].(map[string]any)
	if !ok {
		t.Fatalf("Authorization is missing or not a map")
	}

	jwtInfo, ok := authzInfo["jwt"].(map[string]any)
	if !ok {
		t.Fatalf("jwt is missing in Authorization info")
	}

	if jwtInfo["header"] == nil || jwtInfo["payload"] == nil {
		t.Errorf("Expected header and payload in jwt info, got %v", jwtInfo)
	}
}

func TestEchoHandler_InvalidJSONJWT(t *testing.T) {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{invalid_json}`))
	signature := "dummy"

	jwt := header + "." + payload + "." + signature

	req := httptest.NewRequest(http.MethodGet, "/echo", nil)
	req.Header.Set("Authorization", "Bearer "+jwt)

	w := httptest.NewRecorder()
	echoHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", res.Status)
	}

	var respMap map[string]any
	if err := json.NewDecoder(res.Body).Decode(&respMap); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	authzInfo, ok := respMap["Authorization"].(map[string]any)
	if !ok {
		t.Fatalf("Authorization is missing or not a map")
	}

	_, ok = authzInfo["jwt"]
	if ok {
		t.Errorf("Expected jwt to be missing due to invalid JSON payload")
	}
}
