package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"luma-api/common"
	"luma-api/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

const (
	ModelLuma = "luma"
)

func ChatCompletions(c *gin.Context) {
	err := checkChatConfig()
	if err != nil {
		common.ReturnOpenAIError(c, err, "config_invalid", http.StatusInternalServerError)
		return
	}

	chatSubmitTemp := common.Templates["chat_stream_submit"]
	chatTickTemp := common.Templates["chat_stream_tick"]
	chatRespTemp := common.Templates["chat_resp"]

	var requestData common.GeneralOpenAIRequest
	err = c.ShouldBindJSON(&requestData)
	if err != nil {
		common.ReturnOpenAIError(c, err, "parse_body_failed", http.StatusBadRequest)
		return
	}
	isStream := requestData.Stream
	model := ModelLuma

	chatID := "chatcmpl-" + c.GetString(middleware.RequestIdKey)

	prompt := requestData.Messages[len(requestData.Messages)-1].Content

	// do openai tools get params
	params := make(map[string]any)
	params["aspect_ratio"] = "16:9"
	params["expand_prompt"] = "true"
	params["user_prompt"] = prompt

	reqData, _ := json.Marshal(params)

	resp, err := DoRequest("POST", fmt.Sprintf(common.BaseUrl+SubmitEndpoint), bytes.NewBuffer(reqData), nil)
	if err != nil {
		common.ReturnOpenAIError(c, err, "parse_body_failed", http.StatusBadRequest)
		return
	}
	if resp.StatusCode >= 400 {
		// for error
		common.ReturnOpenAIError(c, err, "parse_body_failed", http.StatusBadRequest)
		return
	}
	data, err := io.ReadAll(resp.Body)

	var genTasks []VideoTask
	err = json.Unmarshal(data, &genTasks)

	if len(genTasks) == 0 {
		common.ReturnOpenAIError(c, err, "parse_body_failed", http.StatusBadRequest)
		return
	}
	taskID := genTasks[0].ID

	timeout := time.After(time.Duration(common.ChatTimeOut) * time.Second)
	tick := time.Tick(5 * time.Second)

	if isStream {
		setSSEHeaders(c)
		c.Writer.WriteHeader(http.StatusOK)
	}
	first := false
	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-timeout:
			if isStream {
				common.SendChatData(c.Writer, model, chatID, "timeout")
			} else {
				common.ReturnOpenAIError(c, fmt.Errorf("get feed task timeout"), "request_timeout", 504)
			}
			return
		case <-tick:
			task, err := fetchTask(taskID)
			if err != nil {
				if requestData.Stream {
					common.SendChatData(c.Writer, model, chatID, common.GetJsonString(common.OpenAIError{
						Error: common.Error{
							Message: err.Error(),
							Code:    "fetch_task_failed",
						},
					}))
				} else {
					common.ReturnOpenAIError(c, err, "fetch_task_failed", 500)
				}
				return
			}
			if isStream {
				if !first {
					if chatSubmitTemp != nil {
						var byteBuf bytes.Buffer
						err = chatSubmitTemp.Execute(&byteBuf, task)
						message := byteBuf.String()
						common.SendChatData(c.Writer, model, chatID, message)
					}
				} else {
					if chatTickTemp != nil {
						var byteBuf bytes.Buffer
						err = chatTickTemp.Execute(&byteBuf, task)
						message := byteBuf.String()
						common.SendChatData(c.Writer, model, chatID, message)
					}
				}
				first = true
			}

			if task.State != "completed" {
				continue
			}

			var byteBuf bytes.Buffer
			err = chatRespTemp.Execute(&byteBuf, task)
			if err != nil {
				if requestData.Stream {
					common.SendChatData(c.Writer, model, chatID, common.GetJsonString(common.OpenAIError{
						Error: common.Error{
							Message: err.Error(),
							Code:    "exec_chat_resp_template_failed",
						},
					}))
				} else {
					common.ReturnOpenAIError(c, err, "exec_chat_resp_template_failed", 500)
				}
				return
			}

			message := byteBuf.String()
			if isStream {
				common.SendChatData(c.Writer, model, chatID, message)
				common.SendChatDone(c.Writer)
				return
			} else {
				responses := openai.ChatCompletionResponse{
					ID:      chatID,
					Object:  "chat.completion",
					Created: time.Now().Unix(),
					Model:   model,
					Choices: []openai.ChatCompletionChoice{
						{
							Index:        0,
							FinishReason: "stop",
							Message: openai.ChatCompletionMessage{
								Content: message,
								Role:    "assistant",
							},
						},
					},
				}
				c.JSON(http.StatusOK, responses)
				return
			}
		}
	}
}

func fetchTask(taskID string) (task *VideoTask, err error) {
	action := "/" + taskID
	url := fmt.Sprintf(common.BaseUrl+GetTaskEndpoint, action)
	resp, err := DoRequest("GET", url, nil, nil)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &task)
	if err != nil {
		return nil, err
	}
	return
}

func checkChatConfig() error {
	if common.ChatOpenaiApiBASE == "" {
		return fmt.Errorf("CHAT_OPENAI_BASE is empty")
	}
	if common.ChatOpenaiApiKey == "" {
		return fmt.Errorf("CHAT_OPENAI_KEY is empty")
	}
	_, ok := common.Templates["chat_resp"]
	if !ok {
		return fmt.Errorf("chat_resp template not found")
	}
	return nil
}

func setSSEHeaders(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
}
