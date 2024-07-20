package main

import (
	"luma-api/common"
	"time"
)

func StartAllKeepAlive() {
	for {
		time.Sleep(10 * time.Second)
		if GetLumaAuth() == "" {
			continue
		}
		_, err := getMe()
		if err != nil {
			common.Logger.Errorf("Luma Keep-alive 失败， err: %s", err.Error())
		}
	}
}
