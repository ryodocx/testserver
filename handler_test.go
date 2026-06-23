package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEchoHandler(t *testing.T) {
	// A helper to perform requests and parse the response map
	performRequest := func(authHeader string) map[string]interface{} {
		req := httptest.NewRequest(http.MethodGet, "/echo", nil)
		if authHeader != "" {
			req.Header.Set("Authorization", authHeader)
		}
		w := httptest.NewRecorder()
		echoHandler(w, req)

		res := w.Result()
		defer res.Body.Close()

		var respMap map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&respMap); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		return respMap
	}

	t.Run("Valid JWT", func(t *testing.T) {
		// Valid base64 encoding of {"alg": "none"} and {"user": "test"}
		// eyJhbGciOiAibm9uZSJ9 = {"alg": "none"}
		// eyJ1c2VyIjogInRlc3QifQ == {"user": "test"}
		validHeader := "eyJhbGciOiAibm9uZSJ9"
		validPayload := "eyJ1c2VyIjogInRlc3QifQ"
		authHeader := "Bearer " + validHeader + "." + validPayload + ".signature"

		respMap := performRequest(authHeader)

		authzInfo, ok := respMap["Authorization"].(map[string]interface{})
		if !ok {
			t.Fatalf("Authorization missing or invalid type")
		}

		jwtInfo, ok := authzInfo["jwt"].(map[string]interface{})
		if !ok {
			t.Fatalf("jwt missing or invalid type")
		}

		header, ok := jwtInfo["header"].(map[string]interface{})
		if !ok {
			t.Fatalf("jwt header missing or invalid type")
		}
		if header["alg"] != "none" {
			t.Errorf("expected alg 'none', got %v", header["alg"])
		}

		payload, ok := jwtInfo["payload"].(map[string]interface{})
		if !ok {
			t.Fatalf("jwt payload missing or invalid type")
		}
		if payload["user"] != "test" {
			t.Errorf("expected user 'test', got %v", payload["user"])
		}
	})

	t.Run("Invalid Base64 JWT", func(t *testing.T) {
		// Invalid base64 characters (e.g. !@#$)
		authHeader := "Bearer !@#$.!@#$.!@#$"

		respMap := performRequest(authHeader)

		authzInfo, ok := respMap["Authorization"].(map[string]interface{})
		if !ok {
			t.Fatalf("Authorization missing or invalid type")
		}

		// Should not contain "jwt" key because parsing failed
		if _, exists := authzInfo["jwt"]; exists {
			t.Errorf("expected jwt to not be present due to parsing error, but got %v", authzInfo["jwt"])
		}
	})

	t.Run("Invalid JSON JWT", func(t *testing.T) {
		// Valid base64 encoding of "not json" -> bm90IGpzb24
		authHeader := "Bearer bm90IGpzb24.bm90IGpzb24.signature"

		respMap := performRequest(authHeader)

		authzInfo, ok := respMap["Authorization"].(map[string]interface{})
		if !ok {
			t.Fatalf("Authorization missing or invalid type")
		}

		// Should not contain "jwt" key because JSON parsing failed
		if _, exists := authzInfo["jwt"]; exists {
			t.Errorf("expected jwt to not be present due to json parsing error, but got %v", authzInfo["jwt"])
		}
	})
}
