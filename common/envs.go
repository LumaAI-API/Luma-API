package common

import (
	"os"
	"time"
)

var Version = "v0.0.0"
var StartTime = time.Now().Unix() // unit: second

var PProfEnabled = os.Getenv("PPROF") == "true"
var DebugEnabled = os.Getenv("DEBUG") == "true"
var LogDir = GetOrDefaultString("LOG_DIR", "./logs")
var RotateLogs = os.Getenv("ROTATE_LOGS") == "true"

var Port = GetOrDefaultString("PORT", "8000")
var Proxy = GetOrDefaultString("PROXY", "")

var BaseUrl = GetOrDefaultString("BASE_URL", "https://internal-api.virginia.labs.lumalabs.ai")
var COOKIE = GetOrDefaultString("COOKIE", "")

var ChatTemplateDir = GetOrDefaultString("CHAT_TEMPLATE_DIR", "./template")
var ChatOpenaiModel = GetOrDefaultString("CHAT_OPENAI_MODEL", "gpt-4o")
var ChatOpenaiApiBASE = GetOrDefaultString("CHAT_OPENAI_BASE", "https://api.openai.com")
var ChatOpenaiApiKey = GetOrDefaultString("CHAT_OPENAI_KEY", "")
var ChatTimeOut = GetOrDefault("CHAT_TIME_OUT", 600) // 任务超时时间
