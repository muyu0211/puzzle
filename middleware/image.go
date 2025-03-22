package middleware

import (
	"fmt"
	"github.com/gogf/gf/v2/net/ghttp"
	"net/http"
)

// ImageSplitValid 校验从前端传回来的分割参数
func ImageSplitValid(r *ghttp.Request) {
	var (
		data = r.GetMap()
	)
	fmt.Println("中间件:", data)
	r.Middleware.Next()
	fmt.Println("Response Headers:", r.Response.Header())
}

// AllowCORS 设置跨域
func AllowCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
