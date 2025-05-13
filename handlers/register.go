// handlers/user.go
package handlers

import (
	"net/http"
	"serverGo/db"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Login       string `json:"login"`
	Password    string `json:"password"` // <- это уже SHA256-хеш от клиента
	Email       string `json:"email"`
	Lastname    string `json:"lastname"`
	Firstname   string `json:"firstname"`
	Patronymic  string `json:"patronymic"`
	DateBirth   string `json:"date_birth"` // формат: "2006-01-02"
	Description string `json:"description"`
}

// Регистрация пользователя
func RegisterUser(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	// Хешируем полученный SHA256-пароль с помощью bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка хеширования пароля"})
		return
	}

	var birthDate time.Time
	if req.DateBirth != "" {
		birthDate, err = time.Parse("2006-01-02", req.DateBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты рождения"})
			return
		}
	}

	// Сохраняем пользователя
	_, err = db.DB.Exec(`CALL public.add_user($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		req.Login,
		string(hashedPassword),
		req.Email,
		req.Lastname,
		req.Firstname,
		req.Patronymic,
		birthDate,
		req.Description,
		nil,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка регистрации"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно зарегистрирован"})
}
