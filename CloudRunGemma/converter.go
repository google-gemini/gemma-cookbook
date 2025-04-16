package main

import (
	genai "cloud.google.com/go/ai/generativelanguage/apiv1beta/generativelanguagepb"
	openai "github.com/openai/openai-go"
	param "github.com/openai/openai-go/packages/param"
	shared "github.com/openai/openai-go/shared"
)

var geminiToOpenAiModelMapping = map[string]string{
	"gemma-3-1b-it":  "gemma3:1b",
	"gemma-3-4b-it":  "gemma3:4b",
	"gemma-3-12b-it": "gemma3:12b",
	"gemma-3-27b-it": "gemma3:27b",
}

var openAiToGeminiModelMapping = map[string]string{
	"gemma3:1b":  "gemma-3-1b-it",
	"gemma3:4b":  "gemma-3-4b-it",
	"gemma3:12b": "gemma-3-12b-it",
	"gemma3:27b": "gemma-3-27b-it",
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

func ConvertGenerateContentRequestToChatCompletionRequest(request *genai.GenerateContentRequest, model string) *openai.ChatCompletionNewParams {
	generationConfig := request.GetGenerationConfig()
	chatCompletionParams := &openai.ChatCompletionNewParams{
		Model:            geminiToOpenAiModelMapping[model],
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
	return chatCompletionParams
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

func convertCompletionUsageToUsageMetadata(usage openai.CompletionUsage) *genai.GenerateContentResponse_UsageMetadata {
	return &genai.GenerateContentResponse_UsageMetadata{
		PromptTokenCount:     int32(usage.PromptTokens),
		CandidatesTokenCount: int32(usage.CompletionTokens),
		TotalTokenCount:      int32(usage.TotalTokens),
	}
}

func ConvertChatCompletionResponseToGenerateContentResponse(response *openai.ChatCompletion) *genai.GenerateContentResponse {
	generateContentResponse := &genai.GenerateContentResponse{
		ModelVersion:  openAiToGeminiModelMapping[response.Model],
		Candidates:    convertChoicesToCandidates(response.Choices),
		UsageMetadata: convertCompletionUsageToUsageMetadata(response.Usage),
	}
	return generateContentResponse
}
