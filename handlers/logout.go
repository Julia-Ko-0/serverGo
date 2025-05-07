// handlers/logout.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogoutUser(c *gin.Context) {
	// Очистка access_token и refresh_token
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Вы вышли из аккаунта"})
}
