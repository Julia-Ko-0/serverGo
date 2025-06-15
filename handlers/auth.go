package handlers

import (
	"fmt"
	"net/http"
	"serverGo/db"
	"serverGo/utils" // Подключаем utils для использования функции GenerateToken
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"` // <- тоже SHA256 от клиента
}

// Авторизация пользователя
func LoginUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	var userID int
	var dbHashedPassword string

	// Ищем пользователя по логину
	err := db.DB.QueryRow(
		`SELECT id_user, password FROM public.user WHERE login_us = $1`,
		req.Login,
	).Scan(&userID, &dbHashedPassword)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	// Сравниваем SHA256-пароль с bcrypt-хешем
	err = bcrypt.CompareHashAndPassword([]byte(dbHashedPassword), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	// Генерируем токены
	accessToken, err := utils.GenerateToken(userID, req.Login, 1*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания access токена"})
		return
	}

	refreshToken, err := utils.GenerateToken(userID, req.Login, 7*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания refresh токена"})
		return
	}

	// Устанавливаем куки
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	c.JSON(http.StatusOK, gin.H{"message": "Авторизация успешна"})
}

// Функция для проверки авторизации
func CheckAuth(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err == nil {
		if claims, err := utils.ValidateToken(accessToken); err == nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "Пользователь авторизован",
				"user_id": claims.UserID,
				"login":   claims.Login,
			})
			return
		}
	}

	// Попробовать refresh_token
	refreshToken, err := c.Cookie("refresh_token")
	fmt.Println("refresh_token:", refreshToken)
	if err == nil {
		if claims, err := utils.ValidateToken(refreshToken); err == nil {
			// Выдать новый access_token
			newAccessToken, err := utils.GenerateToken(claims.UserID, claims.Login, 1*time.Hour)
			if err == nil {
				http.SetCookie(c.Writer, &http.Cookie{
					Name:     "access_token",
					Value:    newAccessToken,
					Path:     "/",
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteStrictMode,
				})

				c.JSON(http.StatusOK, gin.H{
					"message": "Пользователь авторизован (по refresh токену)",
					"user_id": claims.UserID,
					"login":   claims.Login,
				})
				return
			}
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизован"})
}
