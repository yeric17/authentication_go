package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/yeric17/thullo/pkg/config"
	"github.com/yeric17/thullo/pkg/models"
	"github.com/yeric17/thullo/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

func RegisterByEmail(ctx *gin.Context) {
	user := &models.User{}

	if err := ctx.BindJSON(user); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: fmt.Sprintf("Bad Request Map User: %s", err),
		})
		fmt.Println(err)
		return
	}

	if err := user.Create("email"); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: fmt.Sprintf("Bad Request Create User: %s", err),
		})
		fmt.Println(err)
		return
	}

	ctx.JSON(http.StatusCreated, utils.DefaultResponse{
		Data:    user,
		Message: "User is register, check your email for confirmate the account",
	})
}

func LoginByEmail(ctx *gin.Context) {
	user := &models.User{}

	//Parse request body to User struct
	if err := ctx.BindJSON(user); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: fmt.Sprintf("Bad Request Map User: %s", err),
		})
		fmt.Println(err)
		return
	}

	//Validate the required fields for login by email
	if user.Password == "" || user.Email == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Bad Request Credentials: password and email is required",
		})
		return
	}

	prevPass := user.Password

	//Get stored user by request email
	if err := user.GetByEmail(true); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: fmt.Sprintf("Bad Request Get User: %s", err),
		})
		fmt.Println(err)
		return
	}

	//Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(prevPass)); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: fmt.Sprintf("Bad Request Credentials: %s", err),
		})
		fmt.Println(err)
		return
	}

	//Validate if user has confirmed the email account
	if user.Status == int(models.UserStatusWaitingConfirmation) {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Email not confirmed",
		})
		return
	}

	user.Password = ""

	//Get token according with the user data
	tokenString, err := GetTokenByUser(*user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Internal Sever Error",
		})
		fmt.Println(err)
		return
	}

	user.Token = tokenString

	//Create a refresh token for update info user without login
	if err := user.CreateRefreshToken(); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Internal Sever Error",
		})
		fmt.Println(err)
		return
	}

	ctx.JSON(http.StatusOK, utils.DefaultResponse{
		Message: "Login successfully!",
		Data:    user,
	})

}

//AuthByToken: used into app when user already login
func AuthByToken(ctx *gin.Context) {

	token := ctx.Request.Header.Get("Authorization")

	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Token is required in Authorization Header",
		})
		return
	}

	user := &models.User{
		Token: token,
	}

	if err := user.ValidToken(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse{
			Message: "Token is no valid",
		})
		return
	}

	ctx.Next()
}

//AuthByRefreshToken: used into app for get updated user information
func AuthByRefreshToken(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")

	if token == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: "Token is required in Authorization Header",
		})
		return
	}

	user := &models.User{
		RefreshToken: token,
	}

	if err := user.ValidRefreshToken(); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: fmt.Sprintf("Bad Request: %s", err),
		})
		return
	}

	if err := user.CreateRefreshToken(); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Message: "Internal Server Error",
		})
		return
	}
	ctx.JSON(http.StatusOK, utils.DefaultResponse{
		Data:    user,
		Message: "Authenticated!",
	})
}

func GetTokenByUser(user models.User) (string, error) {
	claim := models.UserClaims{
		DataUser: user,
	}

	claim.Id = user.ID
	claim.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := token.SignedString([]byte(config.JWT_SECRET_KEY))

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return tokenString, nil
}

func UpdateUser(ctx *gin.Context) {
	user := &models.User{}

	if err := ctx.BindJSON(user); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
			Message: fmt.Sprintf("Bad request: %s", err),
		})
		return
	}

	if err := user.Update(); err != nil {
		if errors.Is(err, utils.ErrorIdRequired) {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse{
				Message: fmt.Sprintf("Bad request: %s", err),
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{
				Message: "Internal server error",
			})
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.DefaultResponse{
		Message: "User updated!",
	})
}
