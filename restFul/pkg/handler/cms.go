package handler

import (
	"net/http"
	database "service_api/pkg/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Добавить услугу
// @Tags cms/services
// @Description Создать услугу
// @Security ApiKeyAuth
// @ID create-service
// @Accept  json
// @Produce  json
// @Param input body database.Service true "Информация об услуге"
// @Success 200 {integer} integer 1
// @Failure 400 {object} handler.error
// @Router /api/admin/cms/services [post]
func addService(c *gin.Context) {

	var input database.Service

	if err := c.ShouldBindJSON(&input); err != nil {
		newError(c, http.StatusBadRequest, err.Error())
		return
	}

	database.DataSource.Context.Create(&input)
}

// @Summary Удалить услугу
// @Tags cms/services
// @Description Удалить информацию об услуге
// @ID delete-service
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "ID услуги"
// @Success 200 {integer} integer 1
// @Failure 400 {object} handler.error
// @Failure 500 {object} handler.error
// @Failure default {object} handler.error
// @Router /api/admin/cms/services/{id} [delete]
func deleteService(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if id == 0 || err != nil {
		newError(c, http.StatusBadRequest, "id is not valid")
		return
	}
	var service database.Service
	if err := database.DataSource.Context.First(&service, id).Error; err != nil {
		newError(c, http.StatusBadRequest, "record not found")
		return
	}
	if err := database.DataSource.Context.Delete(&service).Error; err != nil {
		newError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Service deleted successfully"})
}

// @Summary Получить список услуг
// @Tags cms/services
// @Description Получить актуальный список услуг с прайс-листом
// @ID get-service
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {array} database.Service "Массив информации об услугах"
// @Router /api/admin/cms/services [get]
func getService(c *gin.Context) {
	var servicesArray []*database.Service
	database.DataSource.Context.Find(&servicesArray)

	//fmt.Println(servicesArray)

	c.JSON(http.StatusOK, servicesArray)
}

// @Summary Обновить информацию об услуге
// @Tags cms/services
// @Description Обновить информацию об услуге
// @ID update-service
// @Accept json
// @Security ApiKeyAuth
// @Produce json
// @Param id path int true "ID услуги"
// @Param input body database.Service true "Информация об услуге"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} handler.error
// @Failure 500 {object} handler.error
// @Failure default {object} handler.error
// @Router /api/admin/cms/services/{id} [put]
func updateService(c *gin.Context) {
	var input database.Service
	id, err := strconv.Atoi(c.Param("id"))

	if id == 0 || err != nil {
		newError(c, http.StatusBadRequest, "id is not valid")
		return
	}

	// fmt.Println(id)
	if err := c.ShouldBindJSON(&input); err != nil {
		newError(c, http.StatusBadRequest, err.Error())
		return
	}

	var service database.Service
	if err := database.DataSource.Context.First(&service, id).Error; err != nil {
		newError(c, http.StatusBadRequest, "record not found")
		return
	}

	input.ID = id

	database.DataSource.Context.Save(&input)
}

// // Employees

// func addEmployees(c *gin.Context) {}

// func deleteEmployees(c *gin.Context) {}

// func getEmployees(c *gin.Context) {}

// func updateEmployees(c *gin.Context) {}

// // Reviews

// func addReview(c *gin.Context) {}

// func getReviews(c *gin.Context) {}

// func updateReview(c *gin.Context) {}

// func deleteReview(c *gin.Context) {}
