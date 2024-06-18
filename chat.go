package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"io"
	"luma-api/common"
	"luma-api/middleware"
	"net/http"
	"time"
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
		common.ReturnErr(c, err, "parse_body_failed", http.StatusBadRequest)
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

	header := map[string]string{
		"Cookie": common.COOKIE,
	}
	for k, v := range CommonHeaders {
		header[k] = v
	}
	resp, err := DoRequest("POST", fmt.Sprintf(common.BaseUrl+"/api/photon/v1/generations/"), bytes.NewBuffer(reqData), header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": map[string]any{
				"reason": err.Error(),
				"code":   1,
			},
		})
	}
	if resp.StatusCode > 201 {
		// for error

	}
	data, err := io.ReadAll(resp.Body)

	var genTasks []VideoTask
	err = json.Unmarshal(data, &genTasks)

	if len(genTasks) == 0 {

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
				common.ReturnErr(c, fmt.Errorf("get feed task timeout"), "request_timeout", 504)
			}
			return
		case <-tick:
			action := "/" + taskID
			url := fmt.Sprintf(common.BaseUrl+"/api/photon/v1/generations%s", action)
			if c.Request.URL.RawQuery != "" {
				url = fmt.Sprintf("%s?%s", url, c.Request.URL.RawQuery)
			}
			resp, err = DoRequest("GET", url, nil, header)

			data, err = io.ReadAll(resp.Body)

			var task VideoTask
			err = json.Unmarshal(data, &task)

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
				relayErr = common.WrapperErr(err, common.ErrCodeInternalError, 500)
				if requestData.Stream {
					common.SendChatData(c.Writer, model, chatID, common.GetJsonString(relayErr))
				} else {
					common.ReturnRelayErr(c, relayErr)
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
