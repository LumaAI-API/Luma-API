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
