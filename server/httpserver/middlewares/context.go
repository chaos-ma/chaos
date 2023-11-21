package middlewares

/**
* created by mengqi on 2023/11/21
 */

import "github.com/gin-gonic/gin"

const (
	UsernameKey = "username"
	KeyUserID   = "userid"
	UserIP      = "ip"
)

// Context 为每个请求添加上下文
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO 自定义扩展
		c.Next()
	}
}
