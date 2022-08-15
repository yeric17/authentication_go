package utils

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/yeric17/thullo/pkg/config"
)

func GetImageByLetter(l string) string {
	letter := strings.ToUpper(l)
	image_url := fmt.Sprintf("%s/images/users/%s.jpg", config.HOST, letter)
	return image_url
}

func RandomString(length int) string {
	stringRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, length)
	for i := range b {
		b[i] = stringRunes[rand.Intn(len(stringRunes))]
	}
	return string(b)
}
