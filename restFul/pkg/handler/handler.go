package handler

import (
	"service_api/pkg/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "service_api/docs"
)

type Handler struct {
	AuthService service.AuthService
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		admin := auth.Group("/admin")
		{
			admin.POST("/sign-in", h.signIn)
		}
	}
	//
	api := router.Group("/api")
	{
		task := api.Group("/task")
		{
			task.GET("/:phoneNumber", getTask)
			task.POST("/", addTask)
		}

		trafficLights := api.Group("/traffic-lights")
		{
			trafficLights.GET("/", getTrafficLights)
		}

		admin := api.Group("/admin", h.validateToken)
		{
			taskManage := admin.Group("/task-manage")
			{
				taskManage.GET("/", getTasks)
				taskManage.PUT("/:id", updateTaskStatus)
				taskManage.DELETE("/:id", deleteTask)

				paged := taskManage.Group("/paged")
				{
					paged.GET("/pagedTasks", getPagedTasks)
					paged.GET("/count", getTaskCount)
				}

			}

			telegramTrust := admin.Group("/telegram-trust")
			{
				telegramTrust.POST("/:id", addTrustedUser)
				telegramTrust.PUT("/:id", updateTrustedUser)
				telegramTrust.DELETE("/:id", deleteTrustedUser)
				telegramTrust.GET("/", getTrustedUsers)
			}

			taskStatus := admin.Group("/task-status")
			{
				taskStatus.GET("/", getTaskStatuses)
			}

			trafficLights := admin.Group("/traffic-lights")
			{
				trafficLights.PUT("/:id", updateTrafficLight)
			}

			cms := admin.Group("/cms")
			{
				services := cms.Group("/services")
				{
					services.PUT("/:id", updateService)
					services.POST("/:id", addService)
					services.DELETE("/:id", deleteService)
					services.GET("/", getService)
				}

				// employees := cms.Group("/employees")
				// {
				// 	employees.PUT("/:id", updateEmployees)
				// 	employees.DELETE("/:id", deleteEmployees)
				// 	employees.GET("/", getEmployees)
				// 	employees.POST("/:id", addEmployees)
				// }

				// reviews := cms.Group("/reviews")
				// {
				// 	reviews.POST("/:id", addReview)
				// 	reviews.PUT("/:id", updateReview)
				// 	reviews.DELETE("/:id", deleteReview)
				// 	reviews.GET("/", getReviews)
				// }
			}
		}
	}

	return router
}
