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

// Context 为每个请求添加上下文, django
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO 大家自己去扩展
		c.Next()
	}
}
