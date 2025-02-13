package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type signInInput struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Авторизация в CMS
// @Tags auth
// @Description Создает аккаунт
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body database.User true "Информация об аккаунте"
// @Success 200 {integer} integer 1
// @Failure 500 {object} handler.error
// @Failure default {object} handler.error
// @Router /auth/admin/sign-in	 [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newError(c, http.StatusInternalServerError, err.Error())

		return
	}

	if token, err := h.AuthService.GenerateToken(input.Name, input.Password); err != nil {
		newError(c, http.StatusInternalServerError, err.Error())

		return
	} else {
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
