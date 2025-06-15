package middleware

import (
	"net/http"
	"os"
	"serverGo/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware — авторизация с поддержкой обновления токенов
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Пытаемся получить access_token
		accessToken, err := c.Cookie("access_token")
		if err == nil {
			if claims, valid := validateJWT(accessToken); valid {
				c.Set("user_id", int(claims["user_id"].(float64)))
				c.Set("login", claims["login"])
				c.Next()
				return
			}
		}

		// 2. access_token невалиден — проверяем refresh_token
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Авторизация требуется"})
			return
		}

		if claims, valid := validateJWT(refreshToken); valid {
			userID := int(claims["user_id"].(float64))
			login := claims["login"].(string)

			// 3. Генерируем новый access_token
			newAccessToken, err := utils.GenerateToken(userID, login, 1*time.Hour)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении токена"})
				return
			}

			// Устанавливаем новый access_token
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     "access_token",
				Value:    newAccessToken,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})

			// Продолжаем сессию
			c.Set("user_id", userID)
			c.Set("login", login)
			c.Next()
			return
		}

		// 4. Оба токена невалидны
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Сессия истекла, выполните вход снова"})
	}
}

// validateJWT — универсальная проверка токена
func validateJWT(tokenStr string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return nil, false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, true
	}
	return nil, false
}
