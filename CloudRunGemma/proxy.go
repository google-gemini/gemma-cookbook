// Package proxy provides a reverse proxy for Google Generative AI API to OpenAI API.
package main

import (
	"bytes"
	genai "cloud.google.com/go/ai/generativelanguage/apiv1beta/generativelanguagepb"
	"encoding/json"
	"fmt"
	openai "github.com/openai/openai-go"
	"google.golang.org/protobuf/encoding/protojson"
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
	defaultPort       = "8080"
	defaultTargetPort = "3000" // This needs to be the same as OLLAMA_HOST which specified in the Dockerfile
	targetHost        = "localhost"
	targetScheme      = "http"
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

func convertRequestBody(originalBodyBytes []byte, action string, model string) ([]byte, error) {
	log.Printf("Receive req body: %s, action: %s", string(originalBodyBytes), action)
	switch action {
	case "generateContent":

		generateContentRequest := &genai.GenerateContentRequest{} 
		err := protojson.Unmarshal(originalBodyBytes, generateContentRequest)
		if err != nil {
			return nil, err
		}
		chatCompletionParams := ConvertGenerateContentRequestToChatCompletionRequest(generateContentRequest, model)
		bodyBytes, err := json.Marshal(chatCompletionParams)
		if err != nil {
			return nil, err
		}
		return bodyBytes, nil
	case "streamGenerateContent", "generateAnswer":
		return originalBodyBytes, nil
	default:
		return originalBodyBytes, nil
	}
}

func convertResponseBody(originalBodyBytes []byte, action string) ([]byte, error) {
	log.Printf("Receive original response body: %s, action: %s", string(originalBodyBytes), action)
	switch action {
	case "generateContent":
		chatCompletion:=&openai.ChatCompletion{}
		err := json.Unmarshal(originalBodyBytes, chatCompletion)
		if err != nil {
			return nil, err
		}
		generateContentResponse := ConvertChatCompletionResponseToGenerateContentResponse(chatCompletion)
		bodyBytes, err := json.Marshal(generateContentResponse)
		if err != nil {
			return nil, err
		}
		return bodyBytes, nil
	case "streamGenerateContent", "generateAnswer":
		return originalBodyBytes, nil
	default:
		return originalBodyBytes, nil
	}
}
func main() {
	// --- Configuration ---
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	targetURLStr := fmt.Sprintf("%s://%s:%s", targetScheme, targetHost, defaultTargetPort)
	targetURL, err := url.Parse(targetURLStr)
	if err != nil {
		log.Fatalf("Error parsing target URL %s: %v", targetURLStr, err)
	}

	// --- Setup Proxy ---
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	// --- Http Handler ---
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received URL apth: %s", r.URL.Path)

		re := regexp.MustCompile(geminiAPIPathRegexPattern)
		matches := re.FindStringSubmatch(r.URL.Path)

		if len(matches) == 4 {
			apiVersion := matches[1]
			modelVersion := matches[2]
			action := matches[3]

			log.Printf("API Version: %s\n", apiVersion)
			log.Printf("Model Version: %s\n", modelVersion)
			log.Printf("Real Path: %s\n", action)

			convertedAPIVersion, convertedAPIVersionOk := geminiToOpenAiAPIVersionhMapping[apiVersion]
			convertedAction, convertedActionOk := geminiToOpenAiActionMapping[action]

			if !(convertedAPIVersionOk && convertedActionOk) {
				log.Fatalf("Error occurred during conversion: convertedAPIVersionOk: %v, convertedActionOk: %v. URL path will not be converted. ", convertedAPIVersionOk, convertedActionOk)
				http.Error(w, "Failed to read request path.", http.StatusInternalServerError)
				return
			}

			// 1. Read original request body
			originalBodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading request body: %v", err)
				http.Error(w, "Failed to read request body", http.StatusInternalServerError)
				return
			}
			// It's important to close the original body
			r.Body.Close()
			// Create a new reader with the read bytes for potential later use (like logging)
			r.Body = io.NopCloser(bytes.NewBuffer(originalBodyBytes))

			bodyBytes, err := convertRequestBody(originalBodyBytes, action, modelVersion)
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

			// --- Modify Response ---
			proxy.ModifyResponse = func(resp *http.Response) error {

				// Handle non-200 responses - pass them through without modification
				if resp.StatusCode != http.StatusOK {
					log.Printf("<<< ModifyResponse: Target returned non-200 status (%d). Forwarding original body.", resp.StatusCode)
					return nil
				}
				targetBodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Error encountered during reading response body: %v", err)
					return fmt.Errorf("Failed to read response for action %s.", action)
				}
				// Must close the original body
				resp.Body.Close()

				finalBodyBytes, err := convertResponseBody(targetBodyBytes, action)
				if err != nil {
					log.Printf("Error encountered during response body conversion: %v", err)
					return fmt.Errorf("Failed to convert response for action %s.", action)
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

		} else {
			log.Printf("URL path does not match the expected format. No conversion applied.")
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
