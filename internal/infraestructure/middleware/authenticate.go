package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/inter-hubly/keeper/internal/infraestructure/config"
	"github.com/inter-hubly/keeper/internal/infraestructure/httprest"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/server"
)

// Authentication middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			httprest.Unauthorized(c)
			return
		}

		if !strings.HasPrefix(tokenString, "Bearer ") {
			httprest.Unauthorized(c)
			return
		}

		tokenString = tokenString[7:]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return []byte(server.GetEnvironment().HashEncrypt), nil
		})

		if err != nil || !token.Valid {
			httprest.Unauthorized(c)
		}
		claims := token.Claims.(jwt.MapClaims)
		if loggedUserCtx, ok := claims["loggedUser"].(map[string]interface{}); ok {
			ctx := hctx.LoggedUser.New(hctx.Logged{
				Username: loggedUserCtx["username"].(string),
				Tenant:   loggedUserCtx["tenant"].(string),
			})
			c.Set(config.LoggedUserContextKey, ctx)
			c.Next()
			return
		}

		httprest.Unauthorized(c)
		return
	}
}

func GetLoggedUser(ctx *gin.Context) (context.Context, *hctx.Logged) {
	loggedUserCtx, exists := ctx.Get(config.LoggedUserContextKey)
	if !exists {
		httprest.Unauthorized(ctx)
		return nil, nil
	}

	if loggedUser, ok := loggedUserCtx.(context.Context); ok {
		logged := loggedUser.Value(hctx.LoggedUser).(hctx.Logged)
		return hctx.LoggedUser.New(logged), &logged
	}
	httprest.Unauthorized(ctx)
	return nil, nil
}
