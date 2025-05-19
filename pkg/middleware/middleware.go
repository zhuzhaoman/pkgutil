package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhuzhaoman/pkgutil/pkg/log"
	"io/ioutil"
	"strings"
)

// 调试 打印body 信息
func HealthCheckMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		data, err := ctx.GetRawData()
		if err != nil {
			log.Error(err.Error())
		}
		//str := strings.Replace(string(data),"\n","",-1)
		log.Error("data:", string(data))
		//log.Error("data:", str)
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data)) // 关键点
		//ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(str))) // 关键点
		ctx.Next()
	}
}

// NoMethodHandler 未找到请求方法的处理函数
//func NoMethodHandler() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		c.JSON(405, gin.H{"message": "方法不被允许"})
//	}
//}

// SkipperFunc 定义中间件跳过函数
type SkipperFunc func(*gin.Context) bool

// 检查请求路径是否包含指定的前缀，如果包含则跳过
func AllowPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

// 检查请求路径是否包含指定的前缀，如果包含则不跳过
func AllowPathPrefixNoSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := c.Request.URL.Path
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return false
			}
		}
		return true
	}
}

// 检查请求方法和路径是否包含指定的前缀，如果包含则跳过
func AllowMethodAndPathPrefixSkipper(prefixes ...string) SkipperFunc {
	return func(c *gin.Context) bool {
		path := JoinRouter(c.Request.Method, c.Request.URL.Path)
		pathLen := len(path)

		for _, p := range prefixes {
			if pl := len(p); pathLen >= pl && path[:pl] == p {
				return true
			}
		}
		return false
	}
}

// JoinRouter 拼接路由
func JoinRouter(method, path string) string {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return fmt.Sprintf("%s%s", strings.ToUpper(method), path)
}
