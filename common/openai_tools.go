package common

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"net/http"
	"time"
)

type GeneralOpenAIRequest struct {
	Model            string      `json:"model,omitempty"`
	Messages         []Message   `json:"messages,omitempty"`
	Stream           bool        `json:"stream,omitempty"`
	MaxTokens        uint        `json:"max_tokens,omitempty"`
	Temperature      float64     `json:"temperature,omitempty"`
	TopP             float64     `json:"top_p,omitempty"`
	TopK             int         `json:"top_k,omitempty"`
	FunctionCall     interface{} `json:"function_call,omitempty"`
	FrequencyPenalty float64     `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64     `json:"presence_penalty,omitempty"`
	ToolChoice       string      `json:"tool_choice,omitempty"`
	Tools            []Tool      `json:"tools,omitempty"`
}

type Message struct {
	Role    string  `json:"role"`
	Content string  `json:"content"`
	Name    *string `json:"name,omitempty"`
}

type Function struct {
	Url         string    `json:"url,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Parameters  Parameter `json:"parameters"`
	Arguments   string    `json:"arguments,omitempty"`
}
type Parameter struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

type Tool struct {
	Id       string   `json:"id,omitempty"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type ChatCompletionsStreamResponse struct {
	Id      string                                `json:"id"`
	Object  string                                `json:"object"`
	Created interface{}                           `json:"created"`
	Model   string                                `json:"model"`
	Choices []ChatCompletionsStreamResponseChoice `json:"choices"`
}

type ChatCompletionsStreamResponseChoice struct {
	Index int `json:"index"`
	Delta struct {
		Content string `json:"content"`
		Role    string `json:"role,omitempty"`
	} `json:"delta"`
	FinishReason *string `json:"finish_reason,omitempty"`
}

type OpenAIError struct {
	Error Error `json:"error"`
}

type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    any    `json:"code"`
}

func ConstructChatCompletionStreamReponse(model, answerID string, answer string) openai.ChatCompletionStreamResponse {
	resp := openai.ChatCompletionStreamResponse{
		ID:      answerID,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []openai.ChatCompletionStreamChoice{
			{
				Index: 0,
				Delta: openai.ChatCompletionStreamChoiceDelta{
					Content: answer,
				},
			},
		},
	}
	return resp
}

func SendChatData(w http.ResponseWriter, model string, chatID, data string) {
	dataRune := []rune(data)
	for _, d := range dataRune {
		respData := ConstructChatCompletionStreamReponse(model, chatID, string(d))
		byteData, _ := json.Marshal(respData)
		_, _ = fmt.Fprintf(w, "data: %s\n\n", string(byteData))
		w.(http.Flusher).Flush()
		time.Sleep(1 * time.Millisecond)
	}
}

func SendChatDone(w http.ResponseWriter) {
	_, _ = fmt.Fprintf(w, "data: [DONE]")
	w.(http.Flusher).Flush()
}

func WrapperOpenAIError(err error, code string) *OpenAIError {
	return &OpenAIError{
		Error: Error{
			Message: err.Error(),
			Code:    code,
		},
	}
}

func ReturnOpenAIError(c *gin.Context, err error, code string, statusCode int) {
	c.JSON(statusCode, OpenAIError{
		Error: Error{
			Message: err.Error(),
			Code:    code,
		},
	})
	return
}
