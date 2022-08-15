package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yeric17/thullo/pkg/config"
	"github.com/yeric17/thullo/pkg/handlers"
)

func main() {

	port := config.PORT

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

	router.Run(":" + port)
}
