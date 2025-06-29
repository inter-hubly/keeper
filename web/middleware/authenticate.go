package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/inter-hubly/keeper/internal/config"
	"github.com/inter-hubly/keeper/web/httprest"
	"github.com/inter-hubly/pilot/hctx"
	"github.com/inter-hubly/pilot/server"
)

// AuthMiddleware middleware
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
			tenantId := loggedUserCtx["tenant"].(string)
			ctx := hctx.LoggedUser.New(hctx.Logged{
				Username: loggedUserCtx["username"].(string),
				Tenant:   tenantId,
				UserId:   loggedUserCtx["userId"].(string),
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
	// Retrieve the logged user context
	loggedUserCtx, exists := ctx.Get(config.LoggedUserContextKey)
	if !exists {
		httprest.Unauthorized(ctx)
		return nil, nil
	}

	// Check if the logged user is of the expected type
	if loggedUser, ok := loggedUserCtx.(context.Context); ok {
		// Extract the Logged data from the context
		logged := loggedUser.Value(hctx.LoggedUser).(hctx.Logged)

		// You can create a new context here, but do not overwrite the *gin.Context
		newCtx := hctx.LoggedUser.New(logged)
		newCtx = hctx.Tenant.New(logged.Tenant)

		// Pass the new context back, but keep the original gin.Context intact
		ctx.Set(config.LoggedUserContextKey, newCtx)

		return newCtx, &logged
	}

	httprest.Unauthorized(ctx)
	return nil, nil
}
