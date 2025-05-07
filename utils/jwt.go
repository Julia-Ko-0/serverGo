package utils

import (
	"errors"
	"os"
	"time"

	"log"

	"github.com/golang-jwt/jwt/v4"
)

// Генерация токена
func GenerateToken(userID int, login string, duration time.Duration) (string, error) {
	// Загружаем секрет из переменной окружения
	jwtKey := os.Getenv("JWT_SECRET")

	// Проверка на пустую переменную окружения
	if jwtKey == "" {
		log.Println("❌ JWT_SECRET is EMPTY! Please check your .env file or environment variables.")
		return "", errors.New("JWT_SECRET is not set in environment variables")
	}

	// Создаем claims для токена
	claims := jwt.MapClaims{
		"user_id": userID,
		"login":   login,
		"exp":     time.Now().Add(duration).Unix(),
	}

	// Создаем новый токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен с использованием секретного ключа
	signedToken, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		log.Println("❌ Error signing token:", err)
		return "", err
	}

	return signedToken, nil
}
