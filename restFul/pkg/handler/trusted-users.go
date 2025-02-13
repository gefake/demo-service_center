package handler

import (
	"net/http"
	database "service_api/pkg/db"
	"service_api/pkg/hooks"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func validateTelegramID(telegramID string) (bool, string) {
	if len(telegramID) == 0 {
		return false, "Telegram ID is required"
	}

	if strings.Index(telegramID, "@") != 0 {
		return false, "Invalid Telegram ID"
	}

	return true, ""
}

// @Summary Добавить доверенный чат Telegram
// @Tags telegram-trust
// @Description Создать новую запись о доверенном чате Telegram
// @ID add-telegram-trust
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param input body database.TrustedTelegramUsers true "Информация о доверенном чате Telegram"
// @Success 200 {integer} integer 1
// @Failure 400 {object} handler.error
// @Router /api/admin/telegram-trust [post]
func addTrustedUser(c *gin.Context) {
	var input database.TrustedTelegramUsers

	telegramID := c.Param("id")

	if ok, err := validateTelegramID(telegramID); !ok {
		newError(c, http.StatusBadRequest, err)
		return
	}

	input.TelegramID = telegramID

	database.DataSource.Context.Create(&input)

	tst := input

	hooks.PostMessageToBot("newTrustedUser", &tst)
}

// @Summary Удалить доверенный чат Telegram
// @Tags telegram-trust
// @Description Удалить информацию о доверенном чате Telegram
// @ID delete-telegram-trust
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID доверенного чата Telegram"
// @Success 200 {integer} integer 1
// @Failure 400 {object} handler.error
// @Failure 500 {object} handler.error
// @Router /api/admin/telegram-trust/{id} [delete]
func getTrustedUsers(c *gin.Context) {
	var trustedArray []*database.TrustedTelegramUsers

	database.DataSource.Context.Find(&trustedArray)

	c.JSON(http.StatusOK, trustedArray)
}

// @Summary Получить список доверенных чатов Telegram
// @Tags telegram-trust
// @Description Получить список всех доверенных чатов Telegram
// @ID get-telegram-trust
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} database.TrustedTelegramUsers "Массив информации о доверенных чатах Telegram"
// @Router /api/admin/telegram-trust [get]
func deleteTrustedUser(c *gin.Context) {
	var input database.TrustedTelegramUsers

	telegramID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		newError(c, http.StatusBadRequest, err.Error())
		return
	}

	input.ID = telegramID

	var user database.TrustedTelegramUsers
	_ = database.DataSource.Context.First(&user, "id = ?", telegramID)
	input.TelegramID = user.TelegramID

	hooks.PostMessageToBot("removeTrustedUser", &input)

	ctx := database.DataSource.Context.Delete(&input)

	if ctx.Error != nil {
		newError(c, http.StatusBadRequest, ctx.Error.Error())
		return
	}
}

// @Summary Обновить информацию о доверенном чате Telegram
// @Tags telegram-trust
// @Description Обновить информацию о доверенном чате Telegram
// @ID update-telegram-trust
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID доверенного чата Telegram"
// @Param input body database.TrustedTelegramUsers true "Информация о доверенном чате Telegram"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} handler.error
// @Failure 500 {object} handler.error
// @Router /api/admin/telegram-trust/{id} [put]
func updateTrustedUser(c *gin.Context) {
	var input database.TrustedTelegramUsers

	telegramID := c.Param("id")

	if ok, err := validateTelegramID(telegramID); !ok {
		newError(c, http.StatusBadRequest, err)
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		newError(c, http.StatusBadRequest, err.Error())
		return
	}

	result := database.DataSource.Context.Model(&database.TrustedTelegramUsers{}).Where("telegram_id = ?", telegramID).Updates(&input)
	if result.Error != nil {
		newError(c, http.StatusInternalServerError, "Failed to update record")
		return
	}
}
