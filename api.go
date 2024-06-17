package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"luma-api/common"
	"net/http"
)

var CommonHeaders = map[string]string{
	"Content-Type": "application/json",
	"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	"Referer":      "https://lumalabs.ai/",
	"Origin":       "https://lumalabs.ai",
	"Accept":       "*/*",
}

// @Summary Submit luma generate video task
// @Schemes
// @Description
// @Accept json
// @Produce json
// @Param body body GenRequest true "submit generate video"
// @Success 200 {object} []VideoTask "generate result"
// @Router /luma/generations [post]
func Generations(c *gin.Context) {
	header := map[string]string{
		"Cookie": common.COOKIE,
	}
	for k, v := range CommonHeaders {
		header[k] = v
	}

	resp, err := DoRequest("POST", fmt.Sprintf(common.BaseUrl+"/api/photon/v1/generations/"), c.Request.Body, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": map[string]any{
				"reason": err.Error(),
				"code":   1,
			},
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": map[string]any{
				"reason": err.Error(),
				"code":   1,
			},
		})
		return
	}
}

// @Summary Submit luma generate video task
// @Schemes
// @Description
// @Accept json
// @Produce json
// @Success 200 {object} []VideoTask "video tasks"
// @Router /luma/generations/{id} [get]
// @Router /luma/generations/ [get]
func Fetch(c *gin.Context) {
	action := c.Param("action")

	header := map[string]string{
		"Cookie": common.COOKIE,
	}
	for k, v := range CommonHeaders {
		header[k] = v
	}
	url := fmt.Sprintf(common.BaseUrl+"/api/photon/v1/generations%s", action)
	if c.Request.URL.RawQuery != "" {
		url = fmt.Sprintf("%s?%s", url, c.Request.URL.RawQuery)
	}
	resp, err := DoRequest("GET", url, nil, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": map[string]any{
				"reason": err.Error(),
				"code":   1,
			},
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": map[string]any{
				"reason": err.Error(),
				"code":   1,
			},
		})
		return
	}
	return
}
