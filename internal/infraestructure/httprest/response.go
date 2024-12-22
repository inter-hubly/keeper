package httprest

import (
	"github.com/gin-gonic/gin"
)

type errResponse struct {
	Error string `json:"error"`
}

func Unauthorized(c *gin.Context) {
	c.JSON(401, errResponse{Error: "Unauthorized"})
}

func Forbidden(c *gin.Context) {
	c.JSON(403, errResponse{Error: "Forbidden"})
}

func Error(c *gin.Context, msg string) {
	c.JSON(500, errResponse{Error: msg})
}

func Ok(c *gin.Context, msg any) {
	c.JSON(200, msg)
}
func Created(c *gin.Context, msg any) {
	c.JSON(201, msg)
}
