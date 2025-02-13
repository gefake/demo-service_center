package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authHeader  = "Authorization"
	userContext = "userID"
)

func (h *Handler) validateToken(c *gin.Context) {
	header := c.GetHeader(authHeader)

	//fmt.Println("auth started")

	if len(header) == 0 {
		newError(c, http.StatusUnauthorized, "empty header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newError(c, http.StatusUnauthorized, "invalid header")
		return
	}

	if headerParts[0] != "Bearer" {
		newError(c, http.StatusUnauthorized, "invalid header")
		return
	}

	userID, err := h.AuthService.ParseToken(headerParts[1])

	if err != nil {
		newError(c, http.StatusUnauthorized, err.Error())

		return
	}

	//fmt.Println("auth finished")

	c.Set(userContext, userID)
}
