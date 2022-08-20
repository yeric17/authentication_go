package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/yeric17/thullo/pkg/config"
	"github.com/yeric17/thullo/pkg/models"
	"github.com/yeric17/thullo/pkg/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		ClientID:     config.GOOGLE_CLIENT_ID,
		ClientSecret: config.GOOGLE_CLIENT_SECRET,
		RedirectURL:  fmt.Sprintf("%s/auth/google/callback", config.HOST),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
)

type GoogleUser struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

func GoogleAuth(ctx *gin.Context) {

	url := googleOauthConfig.AuthCodeURL("randomstate")

	//redirect to google login page
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(ctx *gin.Context) {
	state := ctx.Query("state")
	if state != "randomstate" {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Message: "State don't match",
		})
		return
	}

	code := ctx.Query("code")

	token, err := googleOauthConfig.Exchange(context.Background(), code)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Message: "Code-Token Exchange Failed",
		})
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Message: "User data fetch Failed",
		})
		return
	}

	defer resp.Body.Close()

	userData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Message: "User data fetch Failed",
		})
		return
	}

	googleUser := &GoogleUser{}
	json.Unmarshal(userData, &googleUser)

	user := &models.User{}

	user.Email = googleUser.Email

	if user.Email == "" {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse{
			Message: "Provider has not email for login o register",
		})
		fmt.Println("Provider has not email for login o register")
		return
	}

	if err := user.GetByEmail(false); err != nil {
		user.Name = googleUser.Name
		user.Avatar = googleUser.Picture
		user.Password = utils.RandomString(18)
		user.UniqueName = utils.RandomString(5)
		if err := user.Create("google"); err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse{
				Message: "Provider create error",
			})
			fmt.Println(err)
			return
		}
	}

	tokenString, err := GetTokenByUser(*user)

	if err != nil {

		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Message: "Can not get Token for user",
		})
		fmt.Println(err)
		return
	}

	user.Token = tokenString

	// ctx.JSON(http.StatusOK, utils.DefaultResponse{
	// 	Message: "Authenticated!",
	// 	Data:    user,
	// })

	urlValues := url.Values{}
	urlValues.Set("token", tokenString)
	location := url.URL{Path: config.REDIRECT_GOOGLE_URL, RawQuery: urlValues.Encode()}
	ctx.Redirect(http.StatusSeeOther, location.RequestURI())

}
