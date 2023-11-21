package middlewares

/**
* created by mengqi on 2023/11/21
* 这里使用策略模式， 通过设置不同的策略， 可以实现不同的鉴权方式
 */

import (
	"github.com/gin-gonic/gin"
)

type AuthStrategy interface {
	AuthFunc() gin.HandlerFunc
}

type AuthOperator struct {
	strategy AuthStrategy
}

func (ao *AuthOperator) SetStrategy(strategy AuthStrategy) {
	ao.strategy = strategy
}

func (ao *AuthOperator) AuthFunc() gin.HandlerFunc {
	return ao.strategy.AuthFunc()
}
