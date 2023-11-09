package middleware

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"

	"github.com/rickywei/sparrow/project/logger"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		wb := &bodyWriter{
			body:           &bytes.Buffer{},
			ResponseWriter: ctx.Writer,
		}
		ctx.Writer = wb

		ctx.Next()

		bs := wb.body.Bytes()
		r := make(map[string]any)
		_ = json.Unmarshal(bs, &r)

		logger.L().Info(
			"HTTP",
			zap.String("clientIP", ctx.ClientIP()),
			zap.String("method", ctx.Request.Method),
			zap.String("url", ctx.Request.URL.String()),
			zap.Duration("duration", time.Since(start)),
			zap.Int("status", ctx.Writer.Status()),
			zap.Int("code", cast.ToInt(r["code"])),
			zap.String("msg", cast.ToString(r["msg"])),
			// zap.Any("data", r["data"]),
		)

		wb.ResponseWriter.Write(bs)
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}
