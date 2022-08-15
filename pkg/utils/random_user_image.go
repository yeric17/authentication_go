package utils

import (
	"fmt"
	"strings"

	"github.com/yeric17/thullo/pkg/config"
)

func GetImageByLetter(l string) string {
	letter := strings.ToUpper(l)
	image_url := fmt.Sprintf("%s/images/users/%s.jpg", config.HOST, letter)
	return image_url
}
