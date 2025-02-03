package httprest

import (
	"github.com/gin-gonic/gin"
)

type errResponse struct {
	Error string `json:"error"`
}

func Unauthorized(c *gin.Context) {
	c.JSON(401, errResponse{Error: "Unauthorized"})
	c.Abort()
}

func Forbidden(c *gin.Context) {
	c.JSON(403, errResponse{Error: "Forbidden"})
	c.Abort()
}

func Error(c *gin.Context, msg string) {
	c.JSON(500, errResponse{Error: msg})
	c.Abort()
}

func Ok(c *gin.Context, msg any) {
	c.JSON(200, msg)
}

func Created(c *gin.Context, msg any) {
	if msg == nil {
		c.Status(201)
		c.Abort()
		return
	}
	c.JSON(201, msg)
}
