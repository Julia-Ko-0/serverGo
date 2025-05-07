package handlers

import (
	"encoding/base64"
	"net/http"
	post "serverGo/data"
	"serverGo/db"

	"github.com/gin-gonic/gin"
)

// Отправка сообющения 1in1
func SendMessageToUser(c *gin.Context) {
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
	var req struct {
		ReceiverID  int    `json:"receiver_id"`
		MessageText string `json:"message_text"`
		MessageFile string `json:"message_file"` // base64, можно пустым
	}

	// Читаем JSON из запроса
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON", "details": err.Error()})
		return
	}

	// Декодируем файл, если он есть
	var fileData []byte
	if req.MessageFile != "" {
		decoded, err := base64.StdEncoding.DecodeString(req.MessageFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка декодирования файла", "details": err.Error()})
			return
		}
		fileData = decoded
	}

	// Вызов процедуры
	_, err := db.DB.Exec(`CALL public.send_message_to_user($1, $2, $3, $4)`,
		userID, req.ReceiverID, req.MessageText, fileData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось отправить сообщение", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Сообщение успешно отправлено"})
}

// Отправка сообющения chat/1на1
func SendMessageToChat(c *gin.Context) {
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
	var req post.SendMessageToChat // или просто SendMessageToChat, в зависимости от структуры

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	var fileData interface{} = nil
	if req.PMessageFile != "" {
		decoded, err := base64.StdEncoding.DecodeString(req.PMessageFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка декодирования файла", "details": err.Error()})
			return
		}
		fileData = decoded
	}

	_, err := db.DB.Exec(`CALL public.send_message_to_chat($1, $2, $3, $4)`,
		userID, req.PChatId, req.PMessageText, fileData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при отправке сообщения", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Сообщение успешно отправлено в чат"})
}

// доб папку чатов
func AddChatFolder(c *gin.Context) {
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
	var req post.AddChatFolderRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec(`CALL public.add_chat_folder($1, $2)`, req.FolderName, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении папки", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Папка добавлена"})
}

// доп чат в папку
func AddChatToFolder(c *gin.Context) {
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
	var req post.AddChatToFolderRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec(`CALL public.add_chat_to_folder($1, $2, $3)`, req.FolderID, req.ChatID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении чата в папку", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Чат добавлен в папку"})
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////////////
// доп пост на свою страницу
func AddPostUser(c *gin.Context) {
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

	var req post.AddPostUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec(`CALL public.add_post_user($1, $2, $3, $4)`,
		userID, req.TextPost, req.Header, req.FalePost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении поста", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Пост успешно добавлен"})
}

// репост
func AddRepost(c *gin.Context) {
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

	var req post.AddRepostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec(`CALL public.add_repost($1, $2, $3)`,
		userID, req.PostID, req.TypeRepost)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении репоста", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Репост выполнен"})
}

// доб польз в чс
func AddUserToBlacklist(c *gin.Context) {
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

	var req post.BlockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec(`CALL public.add_user_to_blacklist($1, $2)`, userID, req.BlockedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении в черный список", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Пользователь добавлен в черный список"})
}
