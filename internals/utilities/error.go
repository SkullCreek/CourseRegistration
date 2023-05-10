package utilities

import (
	"fmt"
	"strings"
)

func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// check if admin token is valid
func IsTokenValid(auth string) bool {
	if auth != "" && strings.HasPrefix(auth, "Bearer ") {
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == "admin123" {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
