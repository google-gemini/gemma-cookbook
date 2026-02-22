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

var geminiToOpenAiAPIVersionMapping = map[string]string{
	"v1beta": "v1",
	"v1":     "v1",
}

const geminiAPIPathRegexPattern = `/([^/]+)/models/([^/]+):([^/]+)$`

var geminiAPIPathRegex = regexp.MustCompile(geminiAPIPathRegexPattern)

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
	pr, pw := io.Pipe()
	originalBody := resp.Body
	resp.Body = pr
	resp.Header.Del("Content-Length")
	resp.Header.Set("Transfer-Encoding", "chunked")
	ConvertStreamResponseBody(originalBody, pw, nil)
	return nil
}

type Config struct {
	Port      string
	TargetURL *url.URL
	APIKey    string
}

func loadConfig() (*Config, error) {
	envApiKey := os.Getenv("API_KEY")
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	targetURLStr := os.Getenv("OLLAMA_HOST")
	if targetURLStr == "" {
		targetURLStr = defaultTargetURL
	}

	targetURL, err := url.Parse(targetURLStr)
	if err != nil {
		return nil, fmt.Errorf("error parsing target URL %s: %v", targetURLStr, err)
	}

	return &Config{
		Port:      port,
		TargetURL: targetURL,
		APIKey:    envApiKey,
	}, nil
}

type ProxyHandler struct {
	config *Config
}

func NewProxyHandler(config *Config) *ProxyHandler {
	return &ProxyHandler{config: config}
}

func (h *ProxyHandler) setupGeminiProxy(w http.ResponseWriter, r *http.Request, proxy *httputil.ReverseProxy, matches []string) (string, bool) {
	apiVersion := matches[1]
	modelVersion := matches[2]
	action := matches[3]

	log.Printf("API Version: %s\n", apiVersion)
	log.Printf("Model Version: %s\n", modelVersion)
	log.Printf("Real Path: %s\n", action)

	convertedAPIVersion, convertedAPIVersionOk := geminiToOpenAiAPIVersionMapping[apiVersion]
	convertedAction, convertedActionOk := geminiToOpenAiActionMapping[action]

	if !(convertedAPIVersionOk && convertedActionOk) {
		log.Printf("Error occurred during conversion: convertedAPIVersionOk: %v, convertedActionOk: %v. URL path will not be converted. ", convertedAPIVersionOk, convertedActionOk)
		http.Error(w, "Failed to read request path.", http.StatusNotFound)
		return "", false
	}

	queryParams := r.URL.Query()
	requestApiKey := r.Header.Get("x-goog-api-key")
	if requestApiKey == "" {
		requestApiKey = queryParams.Get("key")
	}
	if h.config.APIKey != "" && requestApiKey != h.config.APIKey {
		log.Printf("Invalid API key provided in the request.")
		http.Error(w, "Permission denied. Invalid API Key.  Configure your SDK to use this URL, and use the API key provided at deployment (https://github.com/google-gemini/gemma-cookbook/blob/main/Demos/Gemma-on-Cloudrun/README.md)", http.StatusForbidden)
		return "", false
	}
	// 1. Read original request body
	originalBodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return "", false
	}
	// It's important to close the original body
	r.Body.Close()
	// Create a new reader with the read bytes for potential later use (like logging)
	r.Body = io.NopCloser(bytes.NewBuffer(originalBodyBytes))

	bodyBytes, err := ConvertRequestBody(originalBodyBytes, action, modelVersion)
	if err != nil {
		log.Printf("Error encountered during request body conversion: %v", err)
		http.Error(w, "Failed to convert request for target.", http.StatusInternalServerError)
		return "", false
	}
	log.Printf("Updated request body: %s", string(bodyBytes))

	proxy.Director = func(req *http.Request) {
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		req.URL.Scheme = h.config.TargetURL.Scheme
		req.URL.Host = h.config.TargetURL.Host
		req.URL.Path = fmt.Sprintf("/%s/%s", convertedAPIVersion, convertedAction)
		req.Header.Set("Content-Type", "application/json")
		req.ContentLength = int64(len(bodyBytes))
		req.Host = h.config.TargetURL.Host

		log.Printf(">>> Director: Proxying request to: %s %s", req.Method, req.URL.String())
		log.Printf(">>> Director: Outgoing Host header: %s", req.Host)
		log.Printf(">>> Director: Outgoing Content-Length: %d", req.ContentLength)
		log.Printf(">>> Director: Outgoing Content-Type: %s", req.Header.Get("Content-Type"))
	}

	return action, true
}

func (h *ProxyHandler) setupStandardProxy(w http.ResponseWriter, r *http.Request, proxy *httputil.ReverseProxy) bool {
	log.Printf("URL path does not match the expected format. No conversion applied.")

	apiKey := r.Header.Get("Authorization")
	expectedPrefix := "Bearer "
	if apiKey == "" {
		apiKey = r.URL.Query().Get("key")
		expectedPrefix = ""
	}
	if h.config.APIKey != "" && apiKey != expectedPrefix+h.config.APIKey {
		log.Printf("Invalid API key provided in the request.")
		http.Error(w, "Permission denied. Invalid API Key.  Configure your SDK to use this URL, and use the API key provided at deployment (https://github.com/google-gemini/gemma-cookbook/blob/main/Demos/Gemma-on-Cloudrun/README.md)", http.StatusForbidden)
		return false
	}

	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = h.config.TargetURL.Scheme
		req.URL.Host = h.config.TargetURL.Host
		req.Host = h.config.TargetURL.Host

		log.Printf(">>> Director: Proxying request to: %s %s", req.Method, req.URL.String())
		log.Printf(">>> Director: Outgoing Host header: %s", req.Host)
	}
	return true
}

func (h *ProxyHandler) proxyErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("!!! Proxy Error: %v", err)
	// Check for specific errors if needed, e.g., connection refused
	http.Error(w, "Proxy Error: "+err.Error(), http.StatusBadGateway)
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// --- Setup Proxy ---
	// We create a new proxy for each request to avoid race conditions on Director and ModifyResponse
	proxy := httputil.NewSingleHostReverseProxy(h.config.TargetURL)

	log.Printf("Received URL path: %s", r.URL.Path)

	matches := geminiAPIPathRegex.FindStringSubmatch(r.URL.Path)

	var (
		isCurrentRequestGeminiStyle bool   = false
		actionForResponse           string // Stores action for response modification
	)

	if len(matches) == 4 {
		isCurrentRequestGeminiStyle = true
		var ok bool
		actionForResponse, ok = h.setupGeminiProxy(w, r, proxy, matches)
		if !ok {
			return
		}

	} else {
		if !h.setupStandardProxy(w, r, proxy) {
			return
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
	proxy.ErrorHandler = h.proxyErrorHandler

	// --- Serve the request using the proxy ---
	log.Printf("Forwarding request to target: %s", h.config.TargetURL.String())
	proxy.ServeHTTP(w, r)
}

func main() {
	// --- Configuration ---
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// --- Http Handler ---
	handler := NewProxyHandler(config)

	// --- Start Server ---
	serverAddr := fmt.Sprintf(":%s", config.Port)
	log.Printf("Proxy server starting on %s", serverAddr)
	log.Printf("Proxying requests to %s", config.TargetURL.String())
	if err := http.ListenAndServe(serverAddr, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
