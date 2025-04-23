package main

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestConvertRequestBody_GenerateContent(t *testing.T) {
	originalBody := `{
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
	model := "gemma-3-1b-it"
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

	actualRequestBody, err := ConvertRequestBody([]byte(originalBody), "generateContent", model)
	if err != nil {
		t.Fatalf("ConvertRequestBody failed: %v", err)
	}

	var actual map[string]interface{}
	var expected map[string]interface{}

	if err := json.Unmarshal(actualRequestBody, &actual); err != nil {
		t.Fatalf("failed to unmarshal actual request body: %v", err)
	}
	if err := json.Unmarshal([]byte(expectedRequestBody), &expected); err != nil {
		t.Fatalf("failed to unmarshal expected request body: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ConvertRequestBody returned incorrect request body. Got: %v, Want: %v", actual, expected)
	}
}

func TestConvertRequestBody_StreamGenerateContent(t *testing.T) {
	originalBody := `{
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
			"responseMimeType": "application/json"
		}
	}`
	model := "gemma-3-4b-it"
	expectedRequestBody := `{
		"StreamOptions":{ "include_usage": true },
		"Stream":true,
		"Model":"gemma3:4b",
		"Messages":[
			{"role":"user","content":"Hello"}
		],
		"MaxTokens":100,
		"Temperature": null,
		"TopP":null,
		"PresencePenalty":null,
		"FrequencyPenalty":null,
		"Stop":["\\n\\n"],
		"ResponseFormat":{"type":"json_object"}
	}`

	actualRequestBody, err := ConvertRequestBody([]byte(originalBody), "streamGenerateContent", model)
	if err != nil {
		t.Fatalf("ConvertRequestBody failed: %v", err)
	}

	var actual map[string]interface{}
	var expected map[string]interface{}

	if err := json.Unmarshal(actualRequestBody, &actual); err != nil {
		t.Fatalf("failed to unmarshal actual request body: %v", err)
	}
	if err := json.Unmarshal([]byte(expectedRequestBody), &expected); err != nil {
		t.Fatalf("failed to unmarshal expected request body: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ConvertRequestBody returned incorrect request body. Got: %v, Want: %v", actual, expected)
	}
}

func TestConvertRequestBody_OtherActions(t *testing.T) {
	originalBody := []byte(`{"some": "data"}`)
	model := "gemma-3-1b-it"
	action := "generateAnswer" // Or any other action

	actualRequestBody, err := ConvertRequestBody(originalBody, action, model)
	if err != nil {
		t.Fatalf("ConvertRequestBody failed: %v", err)
	}

	if !bytes.Equal(actualRequestBody, originalBody) {
		t.Errorf("ConvertRequestBody should return original body for action '%s'. Got: %s, Want: %s", action, string(actualRequestBody), string(originalBody))
	}
}

func TestConvertNonStreamResponseBody_GenerateContent(t *testing.T) {
	originalBody := `{
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

	actualResponseBodyBytes, err := ConvertNonStreamResponseBody([]byte(originalBody), "generateContent")
	if err != nil {
		t.Fatalf("ConvertNonStreamResponseBody failed: %v", err)
	}

	var actual map[string]interface{}
	var expected map[string]interface{}

	if err := json.Unmarshal(actualResponseBodyBytes, &actual); err != nil {
		t.Fatalf("failed to unmarshal actual response body: %v", err)
	}
	if err := json.Unmarshal([]byte(expectedResponseBody), &expected); err != nil {
		t.Fatalf("failed to unmarshal expected response body: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ConvertResponseBody returned incorrect response body. Got: %v, Want: %v", actual, expected)
	}
}

func TestConvertNonStreamResponseBody_OtherActions(t *testing.T) {

	originalBody := []byte(`{"answer": "42"}`)
	action := "generateAnswer"

	actualResponseBody, err := ConvertNonStreamResponseBody(originalBody, action)
	if err != nil {
		t.Fatalf("ConvertNonStreamResponseBody failed: %v", err)
	}

	if !bytes.Equal(actualResponseBody, originalBody) {
		t.Errorf("ConvertNonStreamResponseBody should return original body for action '%s'. Got: %s, Want: %s", action, string(actualResponseBody), string(originalBody))
	}
}

func TestConvertStreamResponseBody(t *testing.T) {

	streamResponse := `
	data: {"id":"chatcmpl-someid1","object":"chat.completion.chunk","created":1687888510,"model":"gemma3:4b","choices":[{"delta":{"content":"Hello"},"index":0,"finish_reason":null}],"usage":{}}
    data: {"id":"chatcmpl-someid1","object":"chat.completion.chunk","created":1687888510,"model":"gemma3:4b","choices":[{"delta":{"content":" from"},"index":0,"finish_reason":null}],"usage":{}}
    data: {"id":"chatcmpl-someid1","object":"chat.completion.chunk","created":1687888510,"model":"gemma3:4b","choices":[{"delta":{"content":" stream!"},"index":0,"finish_reason":"stop"}],"usage":{}}
    data: [DONE]
	`

	expectedResponse := `{"candidates":[{"index":0, "content":{"parts":[{"text":"Hello"}], "role":"model"}}], "usageMetadata":{}, "modelVersion":"gemma-3-4b-it"}
{"candidates":[{"index":0, "content":{"parts":[{"text":" from"}], "role":"model"}}], "usageMetadata":{}, "modelVersion":"gemma-3-4b-it"}
{"candidates":[{"index":0, "content":{"parts":[{"text":" stream!"}], "role":"model"}, "finishReason":"STOP"}], "usageMetadata":{}, "modelVersion":"gemma-3-4b-it"}`

	actualResponse := ``

	reader := io.NopCloser(strings.NewReader(streamResponse))
	pr, pw := io.Pipe()

	convertDone := make(chan struct{})
	readDone := make(chan struct{})
	ConvertStreamResponseBody(reader, pw, convertDone)

	go func() {
		defer close(readDone)
		var buf bytes.Buffer
		_, err := io.Copy(&buf, pr)
		if err != nil {
			t.Errorf("failed to read from pipe: %v", err)
		}
		actualResponse = buf.String()
	}()

	<-convertDone
	<-readDone

	actualLines := strings.Split(strings.TrimSpace(actualResponse), "\n")
	expectedLines := strings.Split(strings.TrimSpace(expectedResponse), "\n")

	if len(actualLines) != len(expectedLines) {
		t.Errorf("ConvertStreamResponseBody returned incorrect number of responses. Got %d, Want %d", len(actualLines), len(expectedLines))
		t.Errorf("Got:\n%v\nWant:\n%v", strings.Join(actualLines, "\n"), strings.Join(expectedLines, "\n"))
		return
	}

	for i := range expectedLines {
		var actual map[string]interface{}
		var expected map[string]interface{}

		if err := json.Unmarshal([]byte(actualLines[i]), &actual); err != nil {
			t.Fatalf("failed to unmarshal actual response body: %v", err)
		}
		if err := json.Unmarshal([]byte(expectedLines[i]), &expected); err != nil {
			t.Fatalf("failed to unmarshal expected response body: %v", err)
		}

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("ConvertResponseBody returned incorrect response body. Got: %v, Want: %v", actual, expected)
		}
	}
}
