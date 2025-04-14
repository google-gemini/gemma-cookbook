package main

import (
	genai "cloud.google.com/go/ai/generativelanguage/apiv1beta/generativelanguagepb"
	openai "github.com/openai/openai-go"
	param "github.com/openai/openai-go/packages/param"
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

func ConvertGenerateContentRequestToChatCompletionRequest(request *genai.GenerateContentRequest, model string) *openai.ChatCompletionNewParams {
	chatCompletionParams := &openai.ChatCompletionNewParams{
		Model:    geminiToOpenAiModelMapping[model],
		Messages: convertContentsToMessages(request.GetContents()),
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

func ConvertChatCompletionResponseToGenerateContentResponse(response *openai.ChatCompletion) *genai.GenerateContentResponse {
	generateContentResponse := &genai.GenerateContentResponse{
		ModelVersion: openAiToGeminiModelMapping[response.Model],
		Candidates:   convertChoicesToCandidates(response.Choices),
	}
	return generateContentResponse
}
