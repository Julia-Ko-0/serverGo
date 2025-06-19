package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

func RemoveRoleGroup(c *gin.Context) {
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
	var req data.DeleteRoleGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Вызываем PostgreSQL-процедуру
	_, err := db.DB.Exec(`CALL delete_role_group($1, $2)`, req.RoleGroupID, callerID)
	if err != nil {
		log.Println("Ошибка при удалении роли:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка при удалении роли",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Роль успешно удалена"})
}

func RemoveUserFromRole(c *gin.Context) {
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

	// Читаем тело запроса
	var req data.RemoveUserFromRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных запроса"})
		return
	}

	// Вызываем процедуру удаления пользователя из роли
	_, err := db.DB.Exec(`CALL delete_user_from_role($1, $2, $3)`,
		req.RoleGroupID, req.UserIDToRemove, callerID)

	if err != nil {
		log.Printf("Ошибка при удалении пользователя из роли: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка при удалении пользователя из роли",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно удалён из роли"})
}
func UpdateGroupInfo(c *gin.Context) {
	var req struct {
		GroupID        int    `json:"group_id"`
		NewName        string `json:"new_name"`
		NewDescription string `json:"new_description"`
		NewPhoto       string `json:"new_photo"` // base64-строка с префиксом
	}

	// Получаем user_id из контекста
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

	// Привязываем данные из JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверный формат JSON",
			"details": err.Error(),
		})
		return
	}

	// Обрабатываем изображение
	var photoBytes []byte
	if req.NewPhoto != "" {
		// Если строка Base64 начинается с префикса, то удаляем его
		if commaIdx := strings.Index(req.NewPhoto, ","); commaIdx != -1 {
			req.NewPhoto = req.NewPhoto[commaIdx+1:]
		}

		// Декодируем Base64 строку
		var err error
		photoBytes, err = base64.StdEncoding.DecodeString(req.NewPhoto)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Ошибка декодирования изображения",
				"details": err.Error(),
			})
			return
		}
	}

	// Вызов SQL-процедуры для обновления информации о группе
	query := `CALL public.update_group_info($1, $2, $3, $4, $5)`
	_, err := db.DB.Exec(
		query,
		req.GroupID,
		userID,
		nullIfEmptyInline(req.NewName),
		photoBytes, // передаем изображение как []byte
		nullIfEmptyInline(req.NewDescription),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка при обновлении группы",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Группа успешно обновлена",
	})
}

// Функция для проверки пустого значения строки
func nullIfEmptyInline(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
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

	// Если группа имеет фото, добавляем префикс base64
	if groupInfo.Photo != "" {
		groupInfo.Photo = "data:image/png;base64," + groupInfo.Photo
	}

	// Отправляем данные клиенту
	c.JSON(http.StatusOK, groupInfo)
}

func ToggleLikePost(c *gin.Context) {
	type ToggleLikeInput struct {
		PostID   int    `json:"post_id"`
		TypePost string `json:"type_post"` // "us" или "gr"
	}

	var input ToggleLikeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if input.TypePost != "us" && input.TypePost != "gr" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type_post must be 'us' or 'gr'"})
		return
	}

	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Вызов функции и получение результата
	var likeAdded bool
	err := db.DB.Get(&likeAdded, "SELECT public.toggle_like_post($1, $2, $3)", input.PostID, userID, input.TypePost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle like", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"like_added": likeAdded, // true или false
	})
}

func AddUserToChat(c *gin.Context) {
	var req struct {
		ChatID int    `json:"chat_id"` // ID чата
		UserID int    `json:"user_id"` // ID пользователя, которого добавляем
		Role   string `json:"role"`    // Роль (опционально)
	}

	// Получаем executor_id из контекста (текущий пользователь)
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден в контексте"})
		return
	}

	executorID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
		return
	}

	// Привязываем JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверный формат JSON",
			"details": err.Error(),
		})
		return
	}

	// Если роль не указана — по умолчанию user
	if req.Role == "" {
		req.Role = "user"
	}

	// Вызов процедуры
	query := `CALL public.add_user_to_chat($1, $2, $3, $4)`
	_, err := db.DB.Exec(query, executorID, req.ChatID, req.UserID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка при добавлении пользователя в чат",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Пользователь успешно добавлен в чат",
	})
}
