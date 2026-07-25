package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/luminous479/JWT-Authentication-with-Gin/controllers"
)

func AuthRoutes(Routes *gin.Engine) {
	Routes.POST("/users/signUp", controllers.SignUp())
	Routes.POST("/users/signIn", controllers.SignIn())
}
