package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"

	genai "cloud.google.com/go/ai/generativelanguage/apiv1beta/generativelanguagepb"
	openai "github.com/openai/openai-go"
	param "github.com/openai/openai-go/packages/param"
	shared "github.com/openai/openai-go/shared"
	"google.golang.org/protobuf/encoding/protojson"
)

var geminiToOpenAiModelMapping = map[string]string{
	"gemma-3-1b-it":  "gemma3:1b",
	"gemma-3-4b-it":  "gemma3:4b",
	"gemma-3-12b-it": "gemma3:12b",
	"gemma-3-27b-it": "gemma3:27b",
	"gemma-3n-e2b-it": "gemma3n:E2b",
	"gemma-3n-e4b-it": "gemma3n:E4b",
}

var openAiToGeminiModelMapping = map[string]string{
	"gemma3:1b":  "gemma-3-1b-it",
	"gemma3:4b":  "gemma-3-4b-it",
	"gemma3:12b": "gemma-3-12b-it",
	"gemma3:27b": "gemma-3-27b-it",
	"gemma3n:E2b": "gemma-3n-e2b-it",
	"gemma3n:E4b": "gemma-3n-e4b-it",
}

type ChatCompletionRequest struct {
	Stream           bool
	StreamOptions    openai.ChatCompletionStreamOptionsParam
	Model            string
	Messages         []openai.ChatCompletionMessageParamUnion
	MaxTokens        param.Opt[int64]
	Temperature      param.Opt[float64]
	TopP             param.Opt[float64]
	PresencePenalty  param.Opt[float64]
	FrequencyPenalty param.Opt[float64]
	ResponseFormat   openai.ChatCompletionNewParamsResponseFormatUnion
	Stop             openai.ChatCompletionNewParamsStopUnion
	// These are supported by openai but not supported by Ollama.
	// Fields supported by Ollama can be found in: https://github.com/ollama/ollama/blob/main/openai/openai.go#L84
	// Logprobs    param.Opt[bool]
	// TopLogprobs param.Opt[int64]
}

func ConvertRequestBody(originalBodyBytes []byte, action string, model string) ([]byte, error) {
	log.Printf("Receive req body: %s, action: %s", string(originalBodyBytes), action)
	switch action {
	case "generateContent":
		return convertGenerateContentRequestToChatCompletionRequest(originalBodyBytes, model, false)
	case "streamGenerateContent":
		return convertGenerateContentRequestToChatCompletionRequest(originalBodyBytes, model, true)
	case "generateAnswer":
		return originalBodyBytes, nil
	default:
		return originalBodyBytes, nil
	}
}

func ConvertNonStreamResponseBody(originalBodyBytes []byte, action string) ([]byte, error) {
	switch action {
	case "generateContent":
		log.Printf("original answer: %s", originalBodyBytes)
		chatCompletion := &openai.ChatCompletion{}
		err := json.Unmarshal(originalBodyBytes, chatCompletion)
		if err != nil {
			return nil, err
		}
		generateContentResponse := convertChatCompletionResponseToGenerateContentResponse(chatCompletion)
		bodyBytes, err := protojson.Marshal(generateContentResponse)
		if err != nil {
			return nil, err
		}
		return bodyBytes, nil
	case "generateAnswer":
		return originalBodyBytes, nil
	default:
		return originalBodyBytes, nil
	}
}

func ConvertStreamResponseBody(originalBody io.ReadCloser, pw *io.PipeWriter, done chan struct{}) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic in modifyStreamResponse goroutine: %v", r)
			}
			originalBody.Close()
			pw.Close()
			close(done)
		}()
		reader := bufio.NewReader(originalBody)

		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Fprintf(pw, "stream read error: %v", err)
				break
			}
			log.Printf("original line: %s", line)

			trimmed := bytes.TrimSpace(line)
			if len(trimmed) == 0 || !bytes.HasPrefix(trimmed, []byte("data: ")) {
				continue
			}

			raw := bytes.TrimSpace(bytes.TrimPrefix(trimmed, []byte("data: ")))

			if bytes.Equal(raw, []byte("[DONE]")) {
				break
			}

			var chunk openai.ChatCompletionChunk
			if err := json.Unmarshal(raw, &chunk); err != nil {
				log.Printf("unmarshal error: %v, raw: '%s'", err, raw)
				fmt.Fprintf(pw, "invalid chunk format, error: %v, raw: %s", err, string(raw))
				continue
			}
			generateContentResponse := convertChatCompletionChunkToGenerateContentResponse(chunk)
			bodyBytes, err := protojson.Marshal(generateContentResponse)
			if err != nil {
				fmt.Fprintf(pw, "failed to convert chunk, error: %v, raw: %s", err, string(raw))
				continue
			}
			pw.Write(append(bodyBytes, '\n'))
		}
	}()
}

func flatten[T any](nested [][]T) []T {
	var flattened []T
	for _, innerList := range nested {
		flattened = append(flattened, innerList...)
	}
	return flattened
}

func convertContentsToMessages(contents []*genai.Content) []openai.ChatCompletionMessageParamUnion {
	nestedMessageList := make([][]openai.ChatCompletionMessageParamUnion, len(contents))
	for i, content := range contents {
		nestedMessageList[i] = convertContentToMessages(content)
	}
	return flatten(nestedMessageList)
}

func convertContentToMessages(content *genai.Content) []openai.ChatCompletionMessageParamUnion {
	openAiChatMessageList := make([]openai.ChatCompletionMessageParamUnion, len(content.GetParts()))
	if content.Role == "model" {
		for i, part := range content.GetParts() {
			openAiChatMessageList[i] = openai.ChatCompletionMessageParamUnion{
				OfAssistant: &openai.ChatCompletionAssistantMessageParam{
					Content: openai.ChatCompletionAssistantMessageParamContentUnion{
						OfString: param.NewOpt(part.GetText()),
					},
					Role: "assistant",
				},
			}
		}
		return openAiChatMessageList

	}
	for i, part := range content.GetParts() {
		openAiChatMessageList[i] = openai.ChatCompletionMessageParamUnion{
			OfUser: &openai.ChatCompletionUserMessageParam{
				Content: openai.ChatCompletionUserMessageParamContentUnion{
					OfString: param.NewOpt(part.GetText()),
				},
				Role: "user",
			},
		}
	}
	return openAiChatMessageList
}

func convertResponseMimeTypeToResponseFormat(responseFormat string) openai.ChatCompletionNewParamsResponseFormatUnion {
	switch responseFormat {
	case "text/plain":
		return openai.ChatCompletionNewParamsResponseFormatUnion{
			OfText: &shared.ResponseFormatTextParam{
				Type: "text",
			},
		}
	case "application/json":
		return openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONObject: &shared.ResponseFormatJSONObjectParam{
				Type: "json_object",
			},
		}
	case "text/x.enum":
		return openai.ChatCompletionNewParamsResponseFormatUnion{}
	default:
		return openai.ChatCompletionNewParamsResponseFormatUnion{}
	}
}

func getOptInt64(number int64) param.Opt[int64] {
	if number == 0 {
		return param.NullOpt[int64]()
	}
	return param.NewOpt(number)
}

func getOptFloat64(number float64) param.Opt[float64] {
	if number == 0 {
		return param.NullOpt[float64]()
	}
	return param.NewOpt(number)
}

func getConvertedModel(model string, modelMap map[string]string) string {
	convertedModel, ok := modelMap[model]
	if ok {
		return convertedModel
	}
	return model
}

func convertGenerateContentRequestToChatCompletionRequest(originalBodyBytes []byte, model string, isStreamRequest bool) ([]byte, error) {
	request := &genai.GenerateContentRequest{}
	err := protojson.Unmarshal(originalBodyBytes, request)
	if err != nil {
		return nil, err
	}
	generationConfig := request.GetGenerationConfig()
	chatCompletionRequest := &ChatCompletionRequest{
		Model:            getConvertedModel(model, geminiToOpenAiModelMapping),
		Messages:         convertContentsToMessages(request.GetContents()),
		MaxTokens:        getOptInt64(int64(generationConfig.GetMaxOutputTokens())),
		Temperature:      getOptFloat64(float64(generationConfig.GetTemperature())),
		TopP:             getOptFloat64(float64(generationConfig.GetTopP())),
		PresencePenalty:  getOptFloat64(float64(generationConfig.GetPresencePenalty())),
		FrequencyPenalty: getOptFloat64(float64(generationConfig.GetFrequencyPenalty())),
		ResponseFormat:   convertResponseMimeTypeToResponseFormat(generationConfig.GetResponseMimeType()),
		Stop: openai.ChatCompletionNewParamsStopUnion{
			OfChatCompletionNewsStopArray: generationConfig.GetStopSequences(),
		},
		// These are supported by openai but not supported by Ollama.
		// Fields supported by Ollama can be found in: https://github.com/ollama/ollama/blob/main/openai/openai.go#L84
		// Logprobs:			 param.NewOpt(generationConfig.GetResponseLogprobs()),
		// TopLogprobs:          param.Newopt(int64(generationConfig.GetLogProbs())),
	}
	if isStreamRequest {
		chatCompletionRequest.Stream = true
		chatCompletionRequest.StreamOptions = openai.ChatCompletionStreamOptionsParam{
			IncludeUsage: param.NewOpt(true),
		}
	}
	bodyBytes, err := json.Marshal(chatCompletionRequest)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func convertMessageToContent(message openai.ChatCompletionMessage) *genai.Content {
	return &genai.Content{
		Role: "model",
		Parts: []*genai.Part{
			{
				Data: &genai.Part_Text{
					Text: message.Content,
				},
			},
		},
	}
}

func convertChunkChoiceDeltaToContent(delta openai.ChatCompletionChunkChoiceDelta) *genai.Content {
	return &genai.Content{
		Role: "model",
		Parts: []*genai.Part{
			{
				Data: &genai.Part_Text{
					Text: delta.Content,
				},
			},
		},
	}
}

func convertOpenaiFinishReasontoGeminiFinishReason(openAiReason string) genai.Candidate_FinishReason {
	switch openAiReason {
	case "stop":
		return genai.Candidate_STOP
	case "length":
		return genai.Candidate_MAX_TOKENS
	case "":
		return genai.Candidate_FINISH_REASON_UNSPECIFIED
	default:
		return genai.Candidate_FINISH_REASON_UNSPECIFIED
	}
}

func convertChoicesToCandidates(choices []openai.ChatCompletionChoice) []*genai.Candidate {
	candidates := make([]*genai.Candidate, len(choices))
	for i, choice := range choices {
		index := int32(choice.Index)
		candidates[i] = &genai.Candidate{
			Content:      convertMessageToContent(choice.Message),
			Index:        &index,
			FinishReason: convertOpenaiFinishReasontoGeminiFinishReason(choice.FinishReason),
		}
	}
	return candidates
}

func convertChunkChoicesToCandidates(choices []openai.ChatCompletionChunkChoice) []*genai.Candidate {
	candidates := make([]*genai.Candidate, len(choices))
	for i, choice := range choices {
		index := int32(choice.Index)
		candidates[i] = &genai.Candidate{
			Content:      convertChunkChoiceDeltaToContent(choice.Delta),
			Index:        &index,
			FinishReason: convertOpenaiFinishReasontoGeminiFinishReason(choice.FinishReason),
		}
	}
	return candidates
}

func convertCompletionUsageToUsageMetadata(usage openai.CompletionUsage) *genai.GenerateContentResponse_UsageMetadata {
	return &genai.GenerateContentResponse_UsageMetadata{
		PromptTokenCount:     int32(usage.PromptTokens),
		CandidatesTokenCount: int32(usage.CompletionTokens),
		TotalTokenCount:      int32(usage.TotalTokens),
	}
}

func convertChatCompletionResponseToGenerateContentResponse(response *openai.ChatCompletion) *genai.GenerateContentResponse {
	generateContentResponse := &genai.GenerateContentResponse{
		ModelVersion:  getConvertedModel(response.Model, openAiToGeminiModelMapping),
		Candidates:    convertChoicesToCandidates(response.Choices),
		UsageMetadata: convertCompletionUsageToUsageMetadata(response.Usage),
	}
	return generateContentResponse
}

func convertChatCompletionChunkToGenerateContentResponse(chunk openai.ChatCompletionChunk) *genai.GenerateContentResponse {
	generateContentResponse := &genai.GenerateContentResponse{
		ModelVersion:  getConvertedModel(chunk.Model, openAiToGeminiModelMapping),
		Candidates:    convertChunkChoicesToCandidates(chunk.Choices),
		UsageMetadata: convertCompletionUsageToUsageMetadata(chunk.Usage),
	}
	return generateContentResponse
}
