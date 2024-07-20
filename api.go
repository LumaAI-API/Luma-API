package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"luma-api/common"
	"net/http"
)

// @Summary Submit luma generate video task
// @Schemes
// @Description
// @Accept json
// @Produce json
// @Param body body GenRequest true "submit generate video"
// @Success 200 {object} []VideoTask "generate result"
// @Router /luma/generations [post]
func Generations(c *gin.Context) {
	doGeneration(c)
}

// @Summary Submit luma extend generate video task
// @Schemes
// @Description
// @Accept json
// @Produce json
// @Param task_id path string true "extend task id"
// @Param body body GenRequest true "submit generate video"
// @Success 200 {object} []VideoTask "generate result"
// @Router /luma/generations/:task_id/extend [post]
func ExtentGenerations(c *gin.Context) {
	doGeneration(c)
}

// @Summary Get luma generate video task
// @Schemes
// @Description
// @Accept json
// @Produce json
// @Param task_id path string true "fetch single task by id"
// @Success 200 {object} VideoTask "video single task"
// @Router /luma/generations/{task_id} [get]
func FetchByID(c *gin.Context) {
	id := c.Param("task_id")
	url := fmt.Sprintf(common.BaseUrl+GetTaskEndpoint, "/"+id)
	if c.Request.URL.RawQuery != "" {
		url = fmt.Sprintf("%s?%s", url, c.Request.URL.RawQuery)
	}
	resp, err := DoRequest("GET", url, nil, nil)
	if err != nil {
		WrapperLumaError(c, err, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	c.Writer.WriteHeader(resp.StatusCode)
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}
	// 读取响应体
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		WrapperLumaError(c, err, http.StatusInternalServerError)
		return
	}
	return
}

// @Summary Upload image to luma
// @Schemes
// @Description
// @Accept json
// @Produce json
// @Param body body UploadReq true "Upload image params"
// @Success 200 {object} []FileUploadResult "upload result"
// @Router /luma/generations/file_upload [post]
func Upload(c *gin.Context) {
	var uploadParams UploadReq
	err := c.BindJSON(&uploadParams)
	if err != nil {
		WrapperLumaError(c, err, http.StatusBadRequest)
		return
	}

	res, relayErr := uploadFile(uploadParams.Url)
	if relayErr != nil {
		ReturnLumaError(c, relayErr.ErrorResp, relayErr.StatusCode)
		return
	}
	c.JSON(http.StatusOK, res)
	return
}

// @Summary Get video url without watermark
// @Schemes
// @Description
// @Accept json
// @Produce json
// @Param task_id path string true "fetch by id"
// @Success 200 {object} object "url"
// @Router /luma/generations/{task_id}/download_video_url [post]
func GetDownloadUrl(c *gin.Context) {
	taskID := c.Param("task_id")

	resp, err := DoRequest(http.MethodGet, fmt.Sprintf(common.BaseUrl+DownloadEndpoint, taskID), nil, nil)
	if err != nil {
		WrapperLumaError(c, err, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	c.Writer.WriteHeader(resp.StatusCode)
	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}
	// 读取响应体
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		WrapperLumaError(c, err, http.StatusInternalServerError)
		return
	}
	return
}

// @Summary Get current user info
// @Schemes
// @Description
// @Accept json
// @Produce json
// @Success 200 {object} object "user info"
// @Router /luma/users/me [get]
func Me(c *gin.Context) {
	res, err := getMe()
	if err != nil {
		WrapperLumaError(c, err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, res)
}

// @Summary Get current user subscription usage
// @Schemes
// @Description
// @Accept json
// @Produce json
// @Success 200 {object} object "subscription info"
// @Router /luma/subscription/usage [get]
func Usage(c *gin.Context) {
	res, err := getUsage()
	if err != nil {
		WrapperLumaError(c, err, http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, res)
}
