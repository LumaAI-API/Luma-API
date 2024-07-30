package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"luma-api/common"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	SubmitEndpoint     = "/api/photon/v1/generations/"
	GetTaskEndpoint    = "/api/photon/v1/generations%s"
	FileUploadEndpoint = "/api/photon/v1/generations/file_upload"
	UsageEndpoint      = "/api/photon/v1/subscription/usage"
	UserinfoEndpoint   = "/api/users/v1/me"
	DownloadEndpoint   = "/api/photon/v1/generations/%s/download_video_url"
)

var CommonHeaders = map[string]string{
	"Content-Type": "application/json",
	"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	"Referer":      "https://lumalabs.ai/",
	"Origin":       "https://lumalabs.ai",
	"Accept":       "*/*",
}

func WrapperLumaError(c *gin.Context, err error, statusCode int) {
	common.Logger.Errorw("wrapper luma error", "statusCode", statusCode, "err", err)
	c.JSON(statusCode, gin.H{
		"detail": map[string]any{
			"reason": err.Error(),
			"code":   1,
		},
	})
}

func GenLumaError(err error, statusCode int) *WrapperErrorResp {
	return &WrapperErrorResp{
		StatusCode: statusCode,
		ErrorResp: ErrorResp{
			Detail: err.Error(),
		},
	}
}

func ReturnLumaError(c *gin.Context, err ErrorResp, statusCode int) {
	common.Logger.Errorw("wrapper luma error", "statusCode", statusCode, "err", err)
	c.JSON(statusCode, err)
}

func GetLumaAuth() string {
	return common.COOKIE
}

var mu sync.Mutex

func HandleRespCookies(resp *http.Response) {
	setCookies := resp.Header.Values("Set-Cookie")
	if len(setCookies) == 0 {
		return
	}
	mu.Lock()
	defer mu.Unlock()

	common.Logger.Infow("Luma账号触发返回cookie")

	//get old cookies
	cookies := make(map[string]string)
	for _, cookie := range strings.Split(common.COOKIE, ";") {
		kv := strings.Split(cookie, "=")
		if len(kv) == 2 {
			cookies[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	for _, cookie := range setCookies {
		names := strings.Split(cookie, "; ")
		kv := strings.Split(names[0], "=")
		if len(kv) == 2 {
			cookies[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	var cookieArr []string
	for k, v := range cookies {
		cookieArr = append(cookieArr, k+"="+v)
	}
	common.COOKIE = strings.Join(cookieArr, "; ")
	return
}

func doGeneration(c *gin.Context) {
	var genRequest GenRequest
	err := c.BindJSON(&genRequest)
	if err != nil {
		WrapperLumaError(c, err, http.StatusBadRequest)
		return
	}
	if genRequest.ImageUrl != "" && !strings.HasPrefix(genRequest.ImageUrl, "https://storage.cdn-luma.com/app_data/photon") {
		uploadRes, relayErr := uploadFile(genRequest.ImageUrl)
		if relayErr != nil {
			ReturnLumaError(c, relayErr.ErrorResp, relayErr.StatusCode)
			return
		}
		common.Logger.Infow("upload file success", "uploadRes", uploadRes)
		genRequest.ImageUrl = uploadRes.PublicUrl
	}

	if genRequest.ImageEndUrl != "" && !strings.HasPrefix(genRequest.ImageEndUrl, "https://storage.cdn-luma.com/app_data/photon") {
		uploadRes, relayErr := uploadFile(genRequest.ImageEndUrl)
		if relayErr != nil {
			ReturnLumaError(c, relayErr.ErrorResp, relayErr.StatusCode)
			return
		}
		common.Logger.Infow("upload file success", "uploadRes", uploadRes)
		genRequest.ImageEndUrl = uploadRes.PublicUrl
	}

	reqData, _ := json.Marshal(genRequest)
	url := common.BaseUrl + SubmitEndpoint
	if strings.HasSuffix(c.Request.URL.Path, "/extend") {
		paths := strings.Split(c.Request.URL.Path, "/generations/")
		url += paths[1]
	}
	resp, err := DoRequest("POST", url, bytes.NewReader(reqData), nil)
	if err != nil {
		WrapperLumaError(c, err, http.StatusInternalServerError)
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
			WrapperLumaError(c, err, http.StatusInternalServerError)
			return
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		WrapperLumaError(c, err, http.StatusInternalServerError)
		return
	}
	var res []any
	err = json.Unmarshal(body, &res)
	if err != nil {
		WrapperLumaError(c, err, http.StatusInternalServerError)
		return
	}

	c.JSON(resp.StatusCode, res[0])
}

// support base64\url
func uploadFile(imgFile string) (*FileUploadResult, *WrapperErrorResp) {
	signedUpload, relayErr := getSignedUpload()
	if relayErr != nil {
		return nil, relayErr
	}

	presignedURL := signedUpload.PresignedUrl

	file, err := readImage(imgFile)
	if err != nil {
		return nil, GenLumaError(err, http.StatusInternalServerError)
	}

	resp, err := DoRequest(http.MethodPut, presignedURL, file, map[string]string{
		"Content-Type": "image/*",
	})
	if err != nil {
		return nil, GenLumaError(err, http.StatusInternalServerError)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return signedUpload, nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, GenLumaError(err, http.StatusInternalServerError)
	}
	return nil, GenLumaError(fmt.Errorf(string(body)), resp.StatusCode)
}

func getSignedUpload() (*FileUploadResult, *WrapperErrorResp) {
	url := common.BaseUrl + FileUploadEndpoint + "?file_type=image&filename=file.jpg"

	resp, err := DoRequest(http.MethodPost, url, nil, nil)
	if err != nil {
		return nil, GenLumaError(err, http.StatusInternalServerError)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, GenLumaError(err, http.StatusInternalServerError)
	}

	if resp.StatusCode >= 400 {
		return nil, GenLumaError(fmt.Errorf(string(body)), resp.StatusCode)
	}

	var result FileUploadResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, GenLumaError(err, http.StatusInternalServerError)
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

func getMe() (map[string]any, error) {
	resp, err := DoRequest(http.MethodGet, common.BaseUrl+UserinfoEndpoint, nil, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	var res = make(map[string]any)
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func getUsage() (map[string]any, error) {
	resp, err := DoRequest(http.MethodGet, common.BaseUrl+UsageEndpoint, nil, nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	var res = make(map[string]any)
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
