package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var TrafficLight int = 1

func getTrafficLights(c *gin.Context) {
	c.JSON(200, TrafficLight)
}

func updateTrafficLight(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		newError(c, http.StatusBadRequest, "Неверный ID")
		return
	}

	if id > 4 {
		newError(c, http.StatusBadRequest, "Неверный ID")
		return
	}

	TrafficLight = id
}
