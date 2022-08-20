package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yeric17/thullo/pkg/config"
	"github.com/yeric17/thullo/pkg/data"
	"github.com/yeric17/thullo/pkg/handlers"
	"github.com/yeric17/thullo/pkg/utils"
)

func main() {

	port := config.PORT

	defer data.Connection.Close()

	if port == "" {
		log.Fatal("Not found env variable PORT")
	}

	router := gin.New()
	router.SetTrustedProxies([]string{"192.168.1.2"})
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Bienvenido")
	})
	router.Static("/images/users", "./public/images/users")

	router.POST("/register/email", handlers.RegisterByEmail)
	router.POST("/login/email", handlers.LoginByEmail)
	router.PUT("/users", handlers.AuthByToken, handlers.UpdateUser)
	router.GET("/auth/token", handlers.AuthByToken, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, utils.DefaultResponse{
			Message: "Authenticated!",
		})
	})
	router.GET("/auth/refresh-token", handlers.AuthByRefreshToken)
	router.GET("/auth/google", handlers.GoogleAuth)
	router.GET("/auth/google/callback", handlers.GoogleCallback)

	router.Run(":" + port)
}
