package middleware

import (
	"net/http"
	"strconv"

	"github.com/rickywei/sparrow/project/conf"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		defer func() {
			if err != nil {
				ctx.AbortWithError(http.StatusUnauthorized, err)
			}
		}()

		tokenStr, err := ctx.Cookie("token")
		if err != nil {
			return
		}

		claims := &jwt.MapClaims{}
		if _, err := jwt.ParseWithClaims(
			tokenStr,
			claims,
			func(t *jwt.Token) (interface{}, error) { return conf.String("jwt.secret"), nil },
			jwt.WithIssuedAt(),
		); err != nil {
			return
		}

		sub, err := claims.GetSubject()
		if err != nil {
			return
		}
		uid, err := strconv.ParseInt(sub, 10, 64)
		if err != nil {
			return
		}

		ctx.Set("uid", uid)

		ctx.Next()
	}
}
