package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"serverGo/data"
	"serverGo/db" // Путь импорта для вашей базы данных

	"github.com/gin-gonic/gin"
)

func AddUserToGroupHandler(c *gin.Context) {
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
	var req data.AddUserToGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec("CALL public.add_user_to_group($1, $2)", req.GroupID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении пользователя в группу", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь добавлен в группу (или уже состоял в ней)"})
}

// GetGroupInfoHandler — обработчик для получения информации о группе
func GetGroupInfoHandler(c *gin.Context) {
	// Получаем group_id из параметра запроса
	groupIDStr := c.Param("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат group_id"})
		return
	}

	// Выполняем SQL-запрос к функции get_group_info
	var groupInfoJSON string
	err = db.DB.QueryRow("SELECT public.get_group_info($1)", groupID).Scan(&groupInfoJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Группа с ID %d не найдена", groupID)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении информации о группе", "details": err.Error()})
		}
		return
	}

	// Преобразуем JSON в структуру
	var groupInfo data.GroupInfoResponse
	err = json.Unmarshal([]byte(groupInfoJSON), &groupInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке данных", "details": err.Error()})
		return
	}

	// Отправляем данные клиенту
	c.JSON(http.StatusOK, groupInfo)
}

// IsUserInGroupHandler проверяет, подписан ли пользователь на группу

// IsUserInGroupHandler проверяет, подписан ли пользователь на группу
func IsUserInGroupHandler(c *gin.Context) {
	// Получаем user_id из middleware (установлен AuthMiddleware)
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

	// Получаем group_id из URL
	groupIDStr := c.Param("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат group_id"})
		return
	}

	// Проверяем наличие связи в user_group
	var isMember bool
	err = db.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM public.user_group
			WHERE user_id = $1 AND group_id = $2
		)
	`, userID, groupID).Scan(&isMember)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке подписки", "details": err.Error()})
		return
	}

	// Возвращаем результат
	c.JSON(http.StatusOK, gin.H{
		"group_id":  groupID,
		"user_id":   userID,
		"is_member": isMember,
	})
}
func UnsubscribeFromGroupHandler(c *gin.Context) {
	// Получаем user_id из middleware
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

	// Получаем group_id из параметра URL
	groupIDStr := c.Param("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат group_id"})
		return
	}

	// Вызываем процедуру в PostgreSQL
	_, err = db.DB.Exec("CALL public.unsubscribe_user_from_group($1, $2)", userID, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка при отписке от группы",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Вы успешно отписались от группы",
	})
}

func RemoveFeatureFromRole(c *gin.Context) {
	// Получаем caller_id из контекста (например, после авторизации)
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не авторизован"})
		return
	}
	callerID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Некорректный ID пользователя"})
		return
	}

	// Читаем данные запроса
	var req data.RemoveFeatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Вызываем PostgreSQL процедуру
	_, err := db.DB.Exec(`CALL remove_feature_from_role($1, $2, $3)`, req.RoleGroupID, req.FeatureID, callerID)
	if err != nil {
		log.Println("Ошибка при вызове процедуры:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка при удалении возможности",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Возможность успешно удалена"})
}
