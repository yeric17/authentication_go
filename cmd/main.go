package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yeric17/thullo/pkg/config"
)

func main() {

	port := config.PORT

	if port == "" {
		log.Fatal("Not found env variable PORT")
	}

	router := gin.New()

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Bienvenido")
	})

	router.Run(":" + port)
}
