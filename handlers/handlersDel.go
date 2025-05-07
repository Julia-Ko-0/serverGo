package handlers

import (
	"net/http"
	delet "serverGo/data"
	"serverGo/db"

	"github.com/gin-gonic/gin"
)

func DeleteSmsFromChat(c *gin.Context) {
	var req delet.DeleteSmsRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	if _, err := db.DB.Exec(`CALL public.delete_sms_from_chat($1)`, req.SmsID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении сообщения", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Сообщение удалено"})
}
func DeleteChatAndUsers(c *gin.Context) {
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
	var req delet.DeleteChatRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	if _, err := db.DB.Exec(`CALL public.delete_chat_and_users($1, $2)`, req.ChatID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении чата", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Чат и все связанные данные удалены"})
}
func DeleteFriendRequest(c *gin.Context) {
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
	var req delet.DeleteFriendRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	if _, err := db.DB.Exec(`CALL public.delete_friend_request($1, $2)`, userID, req.ReceiverID); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении запроса в друзья", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Запрос в друзья удалён"})
}
func DeletePostUser(c *gin.Context) {
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
	var req delet.DeletePostRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	if _, err := db.DB.Exec(`CALL public.delete_post_user($1, $2)`, req.PostID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении поста", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Пост успешно удалён"})
}

// уд пользователя из чата
func RemoveUserFromChat(c *gin.Context) {
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
	var req delet.RemoveUserFromChatRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec(`CALL public.remove_user_from_chat($1, $2, $3)`, req.ExecutorID, req.ChatID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении пользователя из чата", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Пользователь удалён из чата"})
}

// уд пользов из чс
func RemoveUserFromBlacklist(c *gin.Context) {
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

	var req delet.BlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec(`CALL public.remove_user_from_blacklist($1, $2)`, userID, req.BlockedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении из черного списка", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Пользователь удален из черного списка"})
}

// уд репост
func RemoveRepost(c *gin.Context) {
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

	var req delet.RepostRemoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec(`CALL public.remove_repost($1, $2)`, userID, req.RepostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении репоста", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Репост успешно удален"})
}
func RemoveFriend(c *gin.Context) {
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

	var req delet.FriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec(`CALL public.remove_friend($1, $2)`, userID, req.FriendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении друга", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Пользователь удален из друзей"})
}
