package main

import (
	"io"
	"luma-api/common"
	"net/http"
	"net/url"

	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

var TlsHTTPClient tls_client.HttpClient
var HTTPClient *http.Client

func init() {
	InitTlsHTTPClient()
}

func InitTlsHTTPClient() {
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(360),
		tls_client.WithClientProfile(profiles.Chrome_120),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}
	client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	TlsHTTPClient = client

	HTTPClient = &http.Client{}
	if common.Proxy != "" {
		TlsHTTPClient.SetProxy(common.Proxy)

		u, err := url.Parse(common.Proxy)
		if err != nil {
			common.Logger.Error("Invalid proxy URL: " + common.Proxy)
			panic(err)
		}
		HTTPClient = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(u),
			},
		}
	}
}

func DoRequest(method, url string, body io.Reader, otherHeaders map[string]string) (*http.Response, error) {
	headers := map[string]string{
		"Cookie": common.GetLumaAuth(),
	}
	for k, v := range CommonHeaders {
		headers[k] = v
	}
	for k, v := range otherHeaders {
		headers[k] = v
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	defer func() {
		if body != nil {
			req.Body.Close()
		}
	}()

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
