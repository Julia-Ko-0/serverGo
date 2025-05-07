package handlers

import (
	"net/http"
	"serverGo/db"
	"serverGo/utils" // Подключаем utils для использования функции GenerateToken
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	// Проверяем пользователя в базе данных
	var userID int
	err := db.DB.Get(&userID, `SELECT id_user FROM public.user WHERE login_us = $1 AND password = $2`, req.Login, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	// Генерация access токена (срок действия 15 минут)
	accessToken, err := utils.GenerateToken(userID, req.Login, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания access токена"})
		return
	}

	// Генерация refresh токена (срок действия 7 дней)
	refreshToken, err := utils.GenerateToken(userID, req.Login, 7*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания refresh токена"})
		return
	}

	// Устанавливаем куки с токенами (HttpOnly, Secure, SameSite)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,                    // Запрещает доступ к куки через JavaScript
		Secure:   true,                    // Cookie будет отправляться только по HTTPS
		SameSite: http.SameSiteStrictMode, // Защита от CSRF атак
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,                    // Запрещает доступ к куки через JavaScript
		Secure:   true,                    // Cookie будет отправляться только по HTTPS
		SameSite: http.SameSiteStrictMode, // Защита от CSRF атак
	})

	// Отправляем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Авторизация успешна"})
}
