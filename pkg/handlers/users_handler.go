package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yeric17/thullo/pkg/models"
	"github.com/yeric17/thullo/pkg/utils"
)

func RegisterByEmail(ctx *gin.Context) {
	user := &models.User{}

	if err := ctx.BindJSON(user); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			Message:   fmt.Sprintf("Bad Request Map User: %s", err),
		})
		fmt.Println(err)
		return
	}

	if err := user.Create(); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			ErrorCode: http.StatusBadRequest,
			Message:   fmt.Sprintf("Bad Request Create User: %s", err),
		})
		fmt.Println(err)
		return
	}

	ctx.JSON(http.StatusCreated, utils.DefaultResponse{
		Data:    user,
		Message: "User is register, check your email for confirmate the account",
	})
}
