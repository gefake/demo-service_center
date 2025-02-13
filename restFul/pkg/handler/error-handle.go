package handler

import (
	"service_api/pkg/logger"

	"github.com/gin-gonic/gin"
)

type error struct {
	Msg string `json:"message"`
}

func newError(c *gin.Context, status int, msg string) {
	logger.Log.Error(msg)
	c.AbortWithStatusJSON(status, error{Msg: msg})
}
