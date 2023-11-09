package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/rickywei/sparrow/project/logger"
)

func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.L().Debug("recover", zap.Any("panic", r))
			}
		}()

		ctx.Next()
	}
}
