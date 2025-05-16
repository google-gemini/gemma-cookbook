// Package proxy provides a reverse proxy for Google Generative AI API to OpenAI API.
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strconv"
)

// --- Main Application ---

const (
	defaultPort      = "8080"
	defaultTargetURL = "http://localhost:3000"
)

var geminiToOpenAiActionMapping = map[string]string{
	"generateContent":       "chat/completions",
	"streamGenerateContent": "chat/completions",
	"generateAnswer":        "chat/completions",
}

var geminiToOpenAiAPIVersionhMapping = map[string]string{
	"v1beta": "v1",
	"v1":     "v1",
}

const geminiAPIPathRegexPattern = `/([^/]+)/models/([^/]+):([^/]+)$`

func modifyNonStreamResponse(resp *http.Response, action string) error {
	targetBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error encountered during reading response body: %v", err)
		return fmt.Errorf("failed to read response for action %s", action)
	}
	// Must close the original body
	resp.Body.Close()

	finalBodyBytes, err := ConvertNonStreamResponseBody(targetBodyBytes, action)
	if err != nil {
		log.Printf("Error encountered during response body conversion: %v", err)
		return fmt.Errorf("failed to convert response for action %s", action)
	}
	log.Printf("Updated response body: %s", string(finalBodyBytes))

	resp.Body = io.NopCloser(bytes.NewReader(finalBodyBytes))
	resp.Header.Set("Content-Length", strconv.Itoa(len(finalBodyBytes)))
	resp.Header.Set("Content-Type", "application/json") // Ensure correct content type
	// as we are sending plain JSON now.
	resp.Header.Del("Content-Encoding")

	// Setting content length is crucial
	resp.ContentLength = int64(len(finalBodyBytes))
	return nil
}

func modifyStreamResponse(resp *http.Response) error {
	done := make(chan struct{})
	pr, pw := io.Pipe()
	originalBody := resp.Body
	resp.Body = pr
	resp.Header.Del("Content-Length")
	resp.Header.Set("Transfer-Encoding", "chunked")
	ConvertStreamResponseBody(originalBody, pw, done)
	return nil
}

func getApiKey() (string, error) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("environment variable 'API_KEY' is required")
	}
	return apiKey, nil
}

func main() {
	// --- Configuration ---
	envApiKey, err := getApiKey()
	if err != nil {
		log.Fatalf("Error reading API_KEY env var: %v", err)
	}

	port := os.Getenv("PORT")
	targetURLStr := os.Getenv("OLLAMA_HOST")
	if port == "" {
		port = defaultPort
	}
	if targetURLStr == "" {
		targetURLStr = defaultTargetURL
	}

	targetURL, err := url.Parse(targetURLStr)
	if err != nil {
		log.Printf("Error parsing target URL %s: %v", targetURLStr, err)
	}

	// --- Setup Proxy ---
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	// --- Http Handler ---
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received URL apth: %s", r.URL.Path)

		re := regexp.MustCompile(geminiAPIPathRegexPattern)
		matches := re.FindStringSubmatch(r.URL.Path)

		var (
			isCurrentRequestGeminiStyle bool = false
			actionForResponse           string // Stores action for response modification
		)

		if len(matches) == 4 {
			apiVersion := matches[1]
			modelVersion := matches[2]
			action := matches[3]

			log.Printf("API Version: %s\n", apiVersion)
			log.Printf("Model Version: %s\n", modelVersion)
			log.Printf("Real Path: %s\n", action)

			isCurrentRequestGeminiStyle = true
			actionForResponse = action

			convertedAPIVersion, convertedAPIVersionOk := geminiToOpenAiAPIVersionhMapping[apiVersion]
			convertedAction, convertedActionOk := geminiToOpenAiActionMapping[action]

			if !(convertedAPIVersionOk && convertedActionOk) {
				log.Printf("Error occurred during conversion: convertedAPIVersionOk: %v, convertedActionOk: %v. URL path will not be converted. ", convertedAPIVersionOk, convertedActionOk)
				http.Error(w, "Failed to read request path.", http.StatusNotFound)
				return
			}

			queryParams := r.URL.Query()
			requestApiKey := r.Header.Get("x-goog-api-key")
			if requestApiKey == "" {
				requestApiKey = queryParams.Get("key")
			}
			if requestApiKey != envApiKey {
				log.Printf("Invalid API key provided in the request.")
				http.Error(w, "Permission denied. Invalid API Key.  Configure your SDK to use this URL, and use the API key provided at deployment (https://github.com/google-gemini/gemma-cookbook/blob/main/Demos/Gemma-on-Cloudrun/README.md)", http.StatusForbidden)
				return
			}
			// 1. Read original request body
			originalBodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading request body: %v", err)
				http.Error(w, "Failed to read request body", http.StatusBadRequest)
				return
			}
			// It's important to close the original body
			r.Body.Close()
			// Create a new reader with the read bytes for potential later use (like logging)
			r.Body = io.NopCloser(bytes.NewBuffer(originalBodyBytes))

			bodyBytes, err := ConvertRequestBody(originalBodyBytes, action, modelVersion)
			if err != nil {
				log.Printf("Error encountered during request body conversion: %v", err)
				http.Error(w, "Failed to convert request for target.", http.StatusInternalServerError)
				return
			}
			log.Printf("Updated request body: %s", string(bodyBytes))

			proxy.Director = func(req *http.Request) {
				req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
				req.URL.Scheme = targetURL.Scheme
				req.URL.Host = targetURL.Host
				req.URL.Path = fmt.Sprintf("/%s/%s", convertedAPIVersion, convertedAction)
				req.Header.Set("Content-Type", "application/json")
				req.ContentLength = int64(len(bodyBytes))
				req.Host = targetURL.Host

				log.Printf(">>> Director: Proxying request to: %s %s", req.Method, req.URL.String())
				log.Printf(">>> Director: Outgoing Host header: %s", req.Host)
				log.Printf(">>> Director: Outgoing Content-Length: %d", req.ContentLength)
				log.Printf(">>> Director: Outgoing Content-Type: %s", req.Header.Get("Content-Type"))
			}

		} else {
			log.Printf("URL path does not match the expected format. No conversion applied.")
			isCurrentRequestGeminiStyle = false

			apiKey := r.Header.Get("Authorization")
			expectedPrefix := "Bearer "
			if apiKey == "" {
				apiKey = r.URL.Query().Get("key")
				expectedPrefix = ""
			}
			if apiKey != expectedPrefix + envApiKey {
				log.Printf("Invalid API key provided in the request.")
				http.Error(w, "Permission denied. Invalid API Key.  Configure your SDK to use this URL, and use the API key provided at deployment (https://github.com/google-gemini/gemma-cookbook/blob/main/Demos/Gemma-on-Cloudrun/README.md)", http.StatusForbidden)
				return
			}

			proxy.Director = func(req *http.Request) {
				req.URL.Scheme = targetURL.Scheme
				req.URL.Host = targetURL.Host
				req.Host = targetURL.Host

				log.Printf(">>> Director: Proxying request to: %s %s", req.Method, req.URL.String())
				log.Printf(">>> Director: Outgoing Host header: %s", req.Host)
			}
		}

		// --- Modify Response ---
		proxy.ModifyResponse = func(resp *http.Response) error {
			// Handle non-GenAI style requests - pass them through without modification
			if !isCurrentRequestGeminiStyle {
				return nil
			}

			// Handle non-200 responses - pass them through without modification
			if resp.StatusCode != http.StatusOK {
				log.Printf("<<< ModifyResponse: Target returned non-200 status (%d). Forwarding original body.", resp.StatusCode)
				return nil
			}

			switch actionForResponse {
			case "generateContent", "generateAnswer":
				modifyNonStreamResponse(resp, actionForResponse)
			case "streamGenerateContent":
				modifyStreamResponse(resp)
			default:
				return fmt.Errorf("unexpected action %s", actionForResponse)
			}

			return nil
		}

		// --- Error Handling for the Proxy ---
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("!!! Proxy Error: %v", err)
			// Check for specific errors if needed, e.g., connection refused
			http.Error(w, "Proxy Error: "+err.Error(), http.StatusBadGateway)
		}

		// --- Serve the request using the proxy ---
		log.Printf("Forwarding request to target: %s", targetURL.String())
		proxy.ServeHTTP(w, r)
	})

	// --- Start Server ---
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Proxy server starting on %s", serverAddr)
	log.Printf("Proxying requests to %s", targetURL.String())
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
