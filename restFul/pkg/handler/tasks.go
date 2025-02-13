package handler

import (
	"net/http"
	database "service_api/pkg/db"
	"service_api/pkg/helpers"
	"service_api/pkg/hooks"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary Добавить задачу
// @Tags tasks
// @Description Создать новую задачу
// @ID add-task
// @Accept  json
// @Produce  json
// @Param input body database.ApplicationForCall true "Информация о заявке"
// @Success 200 {integer} integer 1
// @Failure 400 {object} handler.error
// @Router /api/task [post]
func addTask(c *gin.Context) {
	var input database.ApplicationForCall
	if err := c.ShouldBindJSON(&input); err != nil {
		newError(c, http.StatusBadRequest, err.Error())
		return
	}

	if !helpers.ValidatePhoneNumber(input.PhoneNumber) {
		newError(c, http.StatusBadRequest, "Неверный формат номера телефона")
		return
	}

	// var existingApplication database.ApplicationForCall
	// if err := database.DataSource.Context.First(&existingApplication, "phone_number = ?", input.PhoneNumber).Error; err == nil {
	// 	newError(c, http.StatusBadRequest, "Заявка с таким номером телефона уже существует")
	// 	return
	// }

	input.Date = int(time.Now().Unix())
	hooks.PostMessageToBot("newTask", &input)
	database.DataSource.Context.Create(&input)
}

// @Summary Получить задачу
// @Tags tasks
// @Description Получить информацию о задаче по номеру телефона
// @ID get-task
// @Accept json
// @Produce json
// @Param phoneNumber path string true "Номер телефона"
// @Success 200 {array} database.ApplicationForCall "Массив информации о заявках"
// @Router /api/task/{phoneNumber} [get]
func getTask(c *gin.Context) {
	var results []*database.ApplicationForCall
	phoneNum := c.Param("phoneNumber")

	if len(phoneNum) == 0 {
		newError(c, http.StatusBadRequest, "Номер телефона пуст")
		return
	}

	if !helpers.ValidatePhoneNumber(phoneNum) {
		newError(c, http.StatusBadRequest, "Неверный формат номера телефона")
		return
	}

	database.DataSource.Context.Find(&results, "phone_number = ?", phoneNum)
	c.JSON(http.StatusOK, results)
}

// @Summary Удалить задачу
// @Tags tasks
// @Description Удалить задачу по ID
// @ID delete-task
// @Accept json
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "ID задачи"
// @Success 200 {integer} integer 1
// @Failure 400 {object} handler.error
// @Router /api/admin/tasks-manage/{id} [delete]
func deleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if id == 0 || err != nil {
		newError(c, http.StatusBadRequest, "Неверный ID")
		return
	}

	if err := database.DataSource.Context.Delete(&database.ApplicationForCall{}, id).Error; err != nil {
		newError(c, http.StatusBadRequest, "Запись не найдена")
		return
	}
}

// @Summary Обновить статус задачи
// @Tags tasks
// @Description Обновить статус задачи по ID
// @ID update-task-status
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID задачи"
// @Param input body database.ApplicationForCall true "Информация о задаче"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} handler.error
// @Router /api/admin/tasks-manage/{id} [put]
func updateTaskStatus(c *gin.Context) {
	var input database.ApplicationForCall
	id, err := strconv.Atoi(c.Param("id"))
	if id == 0 || err != nil {
		newError(c, http.StatusBadRequest, "Неверный ID")
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		newError(c, http.StatusBadRequest, err.Error())
		return
	}

	var service database.ApplicationForCall
	if err := database.DataSource.Context.First(&service, id).Error; err != nil {
		newError(c, http.StatusBadRequest, "Запись не найдена")
		return
	}

	input.ID = id
	database.DataSource.Context.Save(&input)
}

// @Summary Получить все задачи
// @Tags tasks
// @Description Получить список всех задач
// @Security ApiKeyAuth
// @ID get-tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} database.ApplicationForCall "Массив информации о задачах"
// @Router /api/admin/tasks-manage [get]
func getTasks(c *gin.Context) {
	var results []*database.ApplicationForCall
	database.DataSource.Context.Find(&results)
	c.JSON(http.StatusOK, results)
}

// @Summary Получить задачи с пагинацией
// @Description Получить список задач с использованием пагинации на основе параметров страницы и лимита
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query string false "Номер страницы (по умолчанию 1)"
// @Param limit query string false "Лимит элементов на странице (по умолчанию 10)"
// @Success 200 {array} database.ApplicationForCall "Массив информации о задачах"
// @Failure 400 {object} map[string]string "Ошибка при неверном параметре"
// @Router /api/admin/tasks-manage/paged/pagedTasks [get]
func getPagedTasks(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	var results []*database.ApplicationForCall
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	database.DataSource.Context.Offset((page - 1) * limit).Limit(limit).Find(&results)
	c.JSON(http.StatusOK, results)
}

// @Summary Получить количество задач
// @Description Получить общее количество задач
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]int "Количество задач"
// @Router /api/admin/task-count/paged/count [get]
func getTaskCount(c *gin.Context) {
	var count int64
	database.DataSource.Context.Model(&database.ApplicationForCall{}).Count(&count)
	c.JSON(http.StatusOK, gin.H{"count": count})
}

func getTaskStatuses(c *gin.Context) {
	// TODO: Закончить по факту дополнений в тз
}
