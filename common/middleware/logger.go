package middleware

import (
	"bufio"
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/baowk/dilu-core/common/utils"
	"github.com/baowk/dilu-core/common/utils/ips"
	"github.com/baowk/dilu-core/core"
	"github.com/gin-gonic/gin"
)

// LogLayout 日志layout
// type LogLayout struct {
// 	//Metadata  map[string]interface{} // 存储自定义原数据
// 	Method    string //方法
// 	Path      string // 访问路径
// 	Query     string // 携带query
// 	Body      string // 携带body数据
// 	IP        string // ip地址
// 	UserAgent string // 代理
// 	Error     string // 错误
// 	Cost      string // 花费时间
// 	Source    string // 来源
// }

// LoggerToFile 日志记录到文件
func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 处理请求
		var body string
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodGet, http.MethodDelete:
			bf := bytes.NewBuffer(nil)
			wt := bufio.NewWriter(bf)
			_, err := io.Copy(wt, c.Request.Body)
			if err != nil {
				slog.Warn("copy body error", "warn", err)
				err = nil
			}
			rb, _ := io.ReadAll(bf)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(rb))
			body = string(rb)
		}

		c.Next()

		writeLog(startTime, body, c)
	}
}

func writeLog(startTime time.Time, body string, c *gin.Context) {
	// 结束时间
	if c.Request.Method == http.MethodOptions {
		return
	}
	cost := time.Since(startTime)

	if cost.Milliseconds() < 200 {
		core.Log.Info("request", "ip", ips.GetIP(c), "method", c.Request.Method, "path", c.Request.RequestURI,
			"cost", cost, "userAgent", c.Request.UserAgent(), "query", c.Request.URL.RawQuery,
			"body", body, "source", core.Cfg.Server.Name, "reqId", utils.GetReqId(c))
		//,"error", strings.TrimRight(c.Errors.ByType(gin.ErrorTypePrivate).String(, "\n")))
	} else {
		core.Log.Warn("request", "ip", ips.GetIP(c), "method", c.Request.Method, "path", c.Request.RequestURI,
			"cost", cost, "userAgent", c.Request.UserAgent(), "query", c.Request.URL.RawQuery,
			"body", body, "source", core.Cfg.Server.Name, "reqId", utils.GetReqId(c))
		//,"error", strings.TrimRight(c.Errors.ByType(gin.ErrorTypePrivate).String(, "\n")))
	}
}
