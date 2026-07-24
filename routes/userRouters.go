package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/luminous479/JWT-Authentication-with-Gin/controllers"
	"github.com/luminous479/JWT-Authentication-with-Gin/middelwares"
)

func UserRoutes(Routes *gin.Engine) {
	Routes.Use(middelwares.Authenticate())
	Routes.GET("users/", controllers.GetAllUsers())
	Routes.GET("users/:id", controllers.GetUserById())

}
