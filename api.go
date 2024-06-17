package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"luma-api/common"
	"net/http"
	"regexp"
	"strings"
)

const (
	SubmitEndpoint     = "/api/photon/v1/generations/"
	GetTaskEndpoint    = "/api/photon/v1/generations%s"
	FileUploadEndpoint = "/api/photon/v1/generations/file_upload"
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
	var genRequest GenRequest
	err := c.BindJSON(&genRequest)
	if err != nil {
		common.WrapperLumaError(c, err, http.StatusInternalServerError)
		return
	}
	if genRequest.ImageUrl != "" && !strings.HasPrefix(genRequest.ImageUrl, "https://storage.cdn-luma.com/app_data/photon") {
		uploadRes, err := uploadFile(genRequest.ImageUrl)
		if err != nil {
			common.WrapperLumaError(c, err, http.StatusInternalServerError)
			return
		}
		common.Logger.Infow("upload file success", "uploadRes", uploadRes)
		genRequest.ImageUrl = uploadRes.PublicUrl
	}

	reqData, _ := json.Marshal(genRequest)

	resp, err := DoRequest("POST", fmt.Sprintf(common.BaseUrl+SubmitEndpoint), bytes.NewReader(reqData), nil)
	if err != nil {
		common.WrapperLumaError(c, err, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		c.Writer.WriteHeader(resp.StatusCode)
		for key, values := range resp.Header {
			for _, value := range values {
				c.Writer.Header().Add(key, value)
			}
		}
		// 读取响应体
		_, err = io.Copy(c.Writer, resp.Body)
		if err != nil {
			common.WrapperLumaError(c, err, http.StatusInternalServerError)
			return
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		common.WrapperLumaError(c, err, http.StatusInternalServerError)
		return
	}
	var res []any
	err = json.Unmarshal(body, &res)
	if err != nil {
		common.WrapperLumaError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(resp.StatusCode, res[0])
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

	url := fmt.Sprintf(common.BaseUrl+GetTaskEndpoint, action)
	if c.Request.URL.RawQuery != "" {
		url = fmt.Sprintf("%s?%s", url, c.Request.URL.RawQuery)
	}
	resp, err := DoRequest("GET", url, nil, nil)
	if err != nil {
		common.WrapperLumaError(c, err, http.StatusInternalServerError)
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
		common.WrapperLumaError(c, err, http.StatusInternalServerError)
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
		common.WrapperLumaError(c, err, http.StatusInternalServerError)
		return
	}

	res, err := uploadFile(uploadParams.Url)
	if err != nil {
		common.WrapperLumaError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, res)
	return
}

// support base64\url
func uploadFile(imgFile string) (*FileUploadResult, error) {
	signedUpload, err := getSignedUpload()
	if err != nil {
		return nil, err
	}

	presignedURL := signedUpload.PresignedUrl

	file, err := readImage(imgFile)
	if err != nil {
		return nil, err
	}

	resp, err := DoRequest(http.MethodPut, presignedURL, file, map[string]string{
		"Content-Type": "image/*",
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return signedUpload, nil
	}

	return nil, fmt.Errorf("uploud file failed: %s", resp.Status)
}

func getSignedUpload() (*FileUploadResult, error) {
	url := common.BaseUrl + FileUploadEndpoint + "?file_type=image&filename=file.jpg"

	resp, err := DoRequest(http.MethodPost, url, nil, nil)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result FileUploadResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func readImage(image string) (io.Reader, error) {
	if strings.HasPrefix(image, "data:image/") {
		return getImageFromBase64(image)
	}
	return getImageFromUrl(image)
}

var (
	reg = regexp.MustCompile(`data:image/([^;]+);base64,`)
)

func getImageFromBase64(encoded string) (data io.Reader, err error) {
	decoded, err := base64.StdEncoding.DecodeString(reg.ReplaceAllString(encoded, ""))
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(decoded), nil
}

func getImageFromUrl(url string) (data io.Reader, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	dataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return bytes.NewReader(dataBytes), nil
}
