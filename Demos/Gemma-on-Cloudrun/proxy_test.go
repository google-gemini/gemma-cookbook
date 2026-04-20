package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

const apiKey = "some key"

func TestProxy_GenerateContent(t *testing.T) {
	expectedRequestBody := `{
		"Stream":false,
		"StreamOptions": {},
		"Model":"gemma3:1b",
		"Messages":[
			{"role":"user","content":"Hello"}
		],
		"MaxTokens":100,
		"Temperature":null,
		"TopP":null,
		"PresencePenalty":null,
		"FrequencyPenalty":null,
		"Stop":["\\n\\n"],
		"ResponseFormat":{"type":"text"}
	}`

	targetServerResponseBody := `{
		"id": "chatcmpl-someid",
		"object": "chat.completion",
		"created": 1687888509,
		"model": "gemma3:1b",
		"choices": [
			{
				"index": 0,
				"message": {
					"role": "assistant",
					"content": "Hello from OpenAI!"
				},
				"finish_reason": "stop"
			}
		],
		"usage": {
			"prompt_tokens": 10,
			"completion_tokens": 5,
			"total_tokens": 15
		}
	}`

	// 1. Set up a mock target server
	targetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/chat/completions" && r.Method == http.MethodPost {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Error reading target request body: %v", err)
			}

			var actual map[string]interface{}
			var expected map[string]interface{}

			if err := json.Unmarshal(body, &actual); err != nil {
				t.Fatalf("failed to unmarshal actual request body: %v", err)
			}
			if err := json.Unmarshal([]byte(expectedRequestBody), &expected); err != nil {
				t.Fatalf("failed to unmarshal expected request body: %v", err)
			}

			if !reflect.DeepEqual(actual, expected) {
				t.Errorf("ConvertRequestBody returned incorrect request body. Got: %v, Want: %v", actual, expected)
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(targetServerResponseBody))
		} else {
			t.Errorf("Unexpected target request: Method: %s, Path: %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	})
	mockTargetServer := httptest.NewServer(targetHandler)
	defer mockTargetServer.Close()

	// 2. Set up the proxy server to point to the mock target
	os.Setenv("PORT", "8085")                      // Use a test port
	os.Setenv("OLLAMA_HOST", mockTargetServer.URL) // Ensure proxy targets the mock
	os.Setenv("API_KEY", apiKey)

	// Start the proxy server in a goroutine
	go main()

	proxyURL := fmt.Sprintf("http://localhost:%s", os.Getenv("PORT"))

	// 3. Make a request to the proxy server
	reqBody := `{
		"contents": [
			{
				"parts": [
					{"text": "Hello"}
				]
			}
		],
		"generationConfig": {
			"maxOutputTokens": 100,
			"stopSequences": ["\\n\\n"],
			"responseMimeType": "text/plain"
		}
	}`
	req, err := http.NewRequest(http.MethodPost, proxyURL+"/v1beta/models/gemma-3-1b-it:generateContent", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("x-goog-api-key", apiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error making request to proxy: %v", err)
	}
	defer resp.Body.Close()

	// 4. Assert the response from the proxy
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected Content-Type to contain %q, got %q", "application/json", contentType)
	}
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body from proxy: %v", err)
	}
	expectedResponseBody := `{
		"modelVersion": "gemma-3-1b-it",
		"candidates": [
			{
				"content": {
					"role": "model",
					"parts": [
						{
							"text": "Hello from OpenAI!"
						}
					]
				},
				"index": 0,
				"finishReason": "STOP"
			}
		],
		"usageMetadata": {
			"promptTokenCount": 10,
			"candidatesTokenCount": 5,
			"totalTokenCount": 15
		}
	}`

	var actual map[string]interface{}
	var expected map[string]interface{}

	if err := json.Unmarshal(respBodyBytes, &actual); err != nil {
		t.Fatalf("failed to unmarshal actual response body: %v", err)
	}
	if err := json.Unmarshal([]byte(expectedResponseBody), &expected); err != nil {
		t.Fatalf("failed to unmarshal expected response body: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ConvertResponseBody returned incorrect response body. Got: %v, Want: %v", actual, expected)
	}
}
