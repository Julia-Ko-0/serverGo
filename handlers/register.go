// handlers/user.go
package handlers

import (
	"encoding/base64"
	"net/http"
	"serverGo/db"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Login          string `json:"login"`
	Password       string `json:"password"` // SHA256 от клиента
	Email          string `json:"email"`
	Lastname       string `json:"lastname"`
	Firstname      string `json:"firstname"`
	Patronymic     string `json:"patronymic"`
	DateBirth      string `json:"date_birth"` // формат: "2006-01-02"
	Description    string `json:"description"`
	ProfilePicture string `json:"profile_picture"` // base64 (data:image/...;base64,...)
}

// Регистрация пользователя
func RegisterUser(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	// Хешируем пароль (уже хеширован SHA256 на клиенте)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка хеширования пароля"})
		return
	}

	// Парсим дату рождения
	var birthDate time.Time
	if req.DateBirth != "" {
		birthDate, err = time.Parse("2006-01-02", req.DateBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты рождения"})
			return
		}
	}

	// Обработка изображения
	var imageBytesParam interface{} = nil
	if req.ProfilePicture != "" {
		// Удаляем префикс base64 типа "data:image/...;base64,"
		if commaIdx := strings.Index(req.ProfilePicture, ","); commaIdx != -1 {
			req.ProfilePicture = req.ProfilePicture[commaIdx+1:]
		}

		imageBytes, err := base64.StdEncoding.DecodeString(req.ProfilePicture)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка декодирования изображения", "details": err.Error()})
			return
		}

		// Только если изображение не пустое — передаём его
		if len(imageBytes) > 0 {
			imageBytesParam = imageBytes
		}
	}

	// Вызов процедуры добавления пользователя
	_, err = db.DB.Exec(`CALL public.add_user($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		req.Login,
		string(hashedPassword),
		req.Email,
		req.Lastname,
		req.Firstname,
		req.Patronymic,
		birthDate,
		req.Description,
		imageBytesParam,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка регистрации", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно зарегистрирован"})
}
