package middleware

import (
	"log"
	"net/http"
	"time"
)

// ResponseWriter 包装http.ResponseWriter以捕获状态码
type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// NewResponseWriter 创建新的ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
}

// LoggingMiddleware 添加请求日志记录
func LoggingMiddleware(logger *log.Logger) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// 创建自定义ResponseWriter来捕获状态码
			wrapper := NewResponseWriter(w)

			next(wrapper, r)

			duration := time.Since(start)
			logger.Printf("%s %s %d %v %s",
				r.Method, r.URL.Path, wrapper.StatusCode, duration, r.RemoteAddr)
		}
	}
}
