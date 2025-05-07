package handlers

import (
	"net/http"
	post "serverGo/data"
	"serverGo/db"
	"time"

	"github.com/gin-gonic/gin"
)

// изм имя/фамилия/отчество пользователя
func UpdateUserInfo(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден в контексте"})
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
		return
	}
	var input post.UpdateUserInfoRequest

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	// Преобразуем пустые строки в nil
	nullIfEmpty := func(s string) interface{} {
		if s == "" {
			return nil
		}
		return s
	}

	_, err := db.DB.Exec(`CALL public.update_user_info($1, $2, $3, $4)`,
		userID,
		nullIfEmpty(input.Lastname),
		nullIfEmpty(input.Firstname),
		nullIfEmpty(input.Patronymic),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления данных", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Информация обновлена"})
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////////

// изм почты
func UpdateUserEmail(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден в контексте"})
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
		return
	}
	var req post.UpdateUserEmailRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	if _, err := db.DB.Exec(`CALL public.update_user_email($1, $2)`, userID, req.NewEmail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Email успешно обновлён"})
}

// изм даты рождения
func UpdateUserBirthDate(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден в контексте"})
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
		return
	}
	var req post.UpdateUserBirthDateRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты", "details": err.Error()})
		return
	}

	if _, err := db.DB.Exec(`CALL public.update_user_birth_date($1, $2)`, userID, date); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления даты рождения", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Дата рождения обновлена"})
}

// логина(ника)
func UpdateUserLogin(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден в контексте"})
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
		return
	}
	var req post.UpdateUserLoginRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	if _, err := db.DB.Exec(`CALL public.update_user_login($1, $2)`, userID, req.NewLogin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления логина", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Логин успешно обновлён"})
}

// пороля
func UpdateUserPassword(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден в контексте"})
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
		return
	}
	var req post.UpdateUserPasswordRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	if _, err := db.DB.Exec(`CALL public.update_user_password($1, $2)`, userID, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления пароля", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Пароль успешно обновлён"})
}

// изм фото или имя чата
func UpdateGroupChat(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден в контексте"})
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
		return
	}
	var req post.UpdateGroupChatRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec(`CALL public.update_group_chat($1, $2, $3, $4)`, req.ChatID, userID, req.ChatName, req.Pfoto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении чата", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Групповой чат обновлён"})
}
