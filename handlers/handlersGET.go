package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	user "serverGo/data"
	"serverGo/db" // Путь импорта для вашей базы данных

	"github.com/gin-gonic/gin"
)

// Получение ID пользователя по логину
func getUserIDByLogin(login string) (int, error) {
	var userID int
	err := db.DB.Get(&userID, "SELECT id_user FROM public.user WHERE login_us = $1", login)
	if err != nil {
		return 0, fmt.Errorf("user not found: %v", err)
	}
	return userID, nil
}

// login := c.Param("login")
// if login == "" {
// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Login обязателен"})
// 	return
// }

// userID, err := getUserIDByLogin(login)
//
//	if err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден", "details": err.Error()})
//		return
//	}
//
// Получение информации о пользователе по логину
func GetUsersInfo(c *gin.Context) {

	login := c.Param("login")
	if login == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login обязателен"})
		return
	}

	userID, err := getUserIDByLogin(login)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден", "details": err.Error()})
		return
	}

	userInfoStr, err := db.GetUserInfo(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных о пользователе", "details": err.Error()})
		return
	}

	var userResponse user.UserInfoResponse
	if err := json.Unmarshal([]byte(userInfoStr), &userResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_info": userResponse})
}

// пол инф о пользователе
func GetUserInfo(c *gin.Context) {

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
	userInfoStr, err := db.GetUserInfo(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных о пользователе", "details": err.Error()})
		return
	}

	var userResponse user.UserInfoResponse
	if err := json.Unmarshal([]byte(userInfoStr), &userResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_info": userResponse})
}

// // Получение папок чатов
func GetUserFolders(c *gin.Context) {
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

	var result string
	// Запрос к БД
	err := db.DB.Get(&result, "SELECT * from public.get_user_chat_folders_only($1)", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
		return
	}
	var responseData []user.ChatFolder
	if err := json.Unmarshal([]byte(result), &responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, responseData)
}
func GetUserPosts(c *gin.Context) {
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

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Получаем JSON как строку из функции БД
	var result string
	err := db.DB.Get(&result, "SELECT * FROM public.get_user_posts($1, $2, $3)", userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
		return
	}

	var raw user.RawUserPostResponse
	if err := json.Unmarshal([]byte(result), &raw); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	// Преобразуем в финальную структуру, где ImageBase64 уже строка
	var response user.UserPostResponse

	response.UserInfo.ProfilePicture = "data:image/png;base64," + base64.StdEncoding.EncodeToString(raw.UserInfo.ProfilePicture)
	// Копируем посты
	for _, p := range raw.Posts {
		response.Posts = append(response.Posts, user.PostResponse{
			PostID:      p.PostID,
			Header:      p.Header,
			Text:        p.Text,
			DateTime:    p.DateTime,
			Comments:    p.Comments,
			Views:       p.Views,
			Repost:      p.Repost,
			ImageBase64: "data:image/png;base64," + base64.StdEncoding.EncodeToString(p.ImageData),
		})
	}

	c.JSON(http.StatusOK, response)
}

// Получение постов пользователя по логину
func GetUsersPosts(c *gin.Context) {
	login := c.Param("login")
	if login == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login обязателен"})
		return
	}

	userID, err := getUserIDByLogin(login)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден", "details": err.Error()})
		return
	}

	// Читаем лимит и оффсет
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var result string
	err = db.DB.Get(&result, "SELECT * FROM public.get_user_posts($1, $2, $3)", userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
		return
	}

	var responseData user.UserPostResponse
	if err := json.Unmarshal([]byte(result), &responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responseData)
}

// Получение чатов
func GetUserChat(c *gin.Context) {
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

	var result string
	// Запрос к БД
	err := db.DB.Get(&result, "SELECT * from public.get_user_chats($1)", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
		return
	}
	// c.JSON(http.StatusOK, result)
	var responseData []user.UserChats
	if err := json.Unmarshal([]byte(result), &responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responseData)
}

// Получение чатов папки
func GetUserChatFolsers(c *gin.Context) {

	id_chatFolders := c.Param("folder")
	if id_chatFolders == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Folders обязателен"})
		return
	}

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

	var result string
	// Запрос к БД
	err := db.DB.Get(&result, "SELECT * from public.get_user_chats($1, $2)", userID, id_chatFolders)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
		return
	}
	// c.JSON(http.StatusOK, result)
	var responseData []user.UserChats
	if err := json.Unmarshal([]byte(result), &responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responseData)
}

//Получение информации о чате(пользователи)

func GetUserChatInfo(c *gin.Context) {
	chet_id := c.Param("chet_id")
	if chet_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chet_id обязателен"})
		return
	}
	var result string
	// Запрос к БД
	err := db.DB.Get(&result, "SELECT * from public.get_chat_users($1)", chet_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
		return
	}
	// c.JSON(http.StatusOK, result)
	var responseData []user.ChatUser
	if err := json.Unmarshal([]byte(result), &responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responseData)
}

// Получение сообщений чата
func GetUserChatMessenges(c *gin.Context) {

	chet_id := c.Param("chet_id")
	if chet_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chet_id обязателен"})
		return
	}

	// Читаем лимит и оффсет
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var result string
	err := db.DB.Get(&result, "SELECT * FROM public.get_chat_messages($1, $2, $3)", chet_id, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
		return
	}

	var responseData []user.ChatMessage
	if err := json.Unmarshal([]byte(result), &responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responseData)
}

// история ников
func GetUserNameHistory(c *gin.Context) {
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
	var history []user.UserNameHistory

	// Выполнение SQL-запроса к вашей функции
	err := db.DB.Select(&history, "SELECT * FROM public.get_user_name_history($1)", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения истории ников", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

// //////////////////////////////////////////////////////////////////////////////////////////////

// получить репосты my
func GetUserReposts(c *gin.Context) {
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

	rows, err := db.DB.Query(`SELECT * FROM public.get_user_reposts_with_full_info($1)`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении репостов", "details": err.Error()})
		return
	}
	defer rows.Close()

	var reposts []user.RepostResponse
	for rows.Next() {
		var r user.RepostResponse
		var postData []byte
		if err := rows.Scan(&r.RepostID, &r.RepostDate, &postData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сканирования", "details": err.Error()})
			return
		}
		if err := json.Unmarshal(postData, &r.PostInfo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка разбора JSON", "details": err.Error()})
			return
		}
		reposts = append(reposts, r)
	}

	c.JSON(http.StatusOK, reposts)
}

// получить репосты пользователя
func GetUsersReposts(c *gin.Context) {

	login := c.Param("login")
	if login == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Login обязателен"})
		return
	}

	userID, err := getUserIDByLogin(login)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден", "details": err.Error()})
		return
	}

	rows, err := db.DB.Query(`SELECT * FROM public.get_user_reposts_with_full_info($1)`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении репостов", "details": err.Error()})
		return
	}
	defer rows.Close()

	var reposts []user.RepostResponse
	for rows.Next() {
		var r user.RepostResponse
		var postData []byte
		if err := rows.Scan(&r.RepostID, &r.RepostDate, &postData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сканирования", "details": err.Error()})
			return
		}
		if err := json.Unmarshal(postData, &r.PostInfo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка разбора JSON", "details": err.Error()})
			return
		}
		reposts = append(reposts, r)
	}

	c.JSON(http.StatusOK, reposts)
}
