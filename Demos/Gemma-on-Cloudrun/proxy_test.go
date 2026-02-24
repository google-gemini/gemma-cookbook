package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	mockTargetURL, err := url.Parse(mockTargetServer.URL)
	if err != nil {
		t.Fatalf("Failed to parse mock target URL: %v", err)
	}

	config := &Config{
		Port:      "8085",
		TargetURL: mockTargetURL,
		APIKey:    apiKey,
	}
	handler := NewProxyHandler(config)
	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	proxyURL := proxyServer.URL

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

func TestProxy_StreamGenerateContent(t *testing.T) {
	expectedRequestBody := `{
		"Stream":true,
		"StreamOptions": {"include_usage": true},
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

	targetServerResponseBody := `data: {"id":"chatcmpl-someid1","object":"chat.completion.chunk","created":1687888510,"model":"gemma3:1b","choices":[{"delta":{"content":"Hello"},"index":0,"finish_reason":null}],"usage":{}}

data: {"id":"chatcmpl-someid1","object":"chat.completion.chunk","created":1687888510,"model":"gemma3:1b","choices":[{"delta":{"content":" stream!"},"index":0,"finish_reason":"stop"}],"usage":{}}

data: [DONE]
`

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

			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(targetServerResponseBody))
		} else {
			t.Errorf("Unexpected target request: Method: %s, Path: %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	})
	mockTargetServer := httptest.NewServer(targetHandler)
	defer mockTargetServer.Close()

	// 2. Set up the proxy server
	mockTargetURL, err := url.Parse(mockTargetServer.URL)
	if err != nil {
		t.Fatalf("Failed to parse mock target URL: %v", err)
	}

	config := &Config{
		Port:      "8085",
		TargetURL: mockTargetURL,
		APIKey:    apiKey,
	}
	handler := NewProxyHandler(config)
	proxyServer := httptest.NewServer(handler)
	defer proxyServer.Close()

	proxyURL := proxyServer.URL

	// 3. Make a streaming request
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
	req, err := http.NewRequest(http.MethodPost, proxyURL+"/v1beta/models/gemma-3-1b-it:streamGenerateContent", bytes.NewBufferString(reqBody))
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

	// 4. Assert the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
	// contentType := resp.Header.Get("Content-Type")
	// Note: The proxy sets "Transfer-Encoding: chunked" but doesn't explicitly set Content-Type for streaming
	// unless explicitly handled. Wait, in `proxy.go`:
	// ConvertStreamResponseBody writes raw JSON lines. It doesn't set "text/event-stream" or "application/json".
	// The original ModifyResponse sets `resp.Header.Set("Transfer-Encoding", "chunked")`.
	// Let's verify what we get. The original logic returned JSON objects separated by newlines?
	// `ConvertStreamResponseBody` writes `json_marshaled_chunk` + `\n`.

	reader := bufio.NewReader(resp.Body)

	expectedLines := []string{
		`{"candidates":[{"index":0, "content":{"parts":[{"text":"Hello"}], "role":"model"}}], "usageMetadata":{}, "modelVersion":"gemma-3-1b-it"}`,
		`{"candidates":[{"index":0, "content":{"parts":[{"text":" stream!"}], "role":"model"}, "finishReason":"STOP"}], "usageMetadata":{}, "modelVersion":"gemma-3-1b-it"}`,
	}

	for i, expectedLine := range expectedLines {
		line, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Error reading line %d: %v", i, err)
		}
		line = strings.TrimSpace(line)

		// Parse both as JSON to compare structurally
		var actual map[string]interface{}
		var expected map[string]interface{}

		if err := json.Unmarshal([]byte(line), &actual); err != nil {
			t.Fatalf("Line %d is not valid JSON: %s", i, line)
		}
		if err := json.Unmarshal([]byte(expectedLine), &expected); err != nil {
			t.Fatalf("Expected line %d is not valid JSON: %s", i, expectedLine)
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Line %d mismatch.\nGot:  %v\nWant: %v", i, actual, expected)
		}
	}
}
