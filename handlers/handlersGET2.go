package handlers

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"serverGo/data"
	"serverGo/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Отправляет запрос в друзья
func SendFriendRequest(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	var req data.FriendActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec(`CALL public.send_friend_request($1, $2)`, userID, req.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при отправке запроса", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Запрос в друзья отправлен"})
}

// Принимает входящий запрос в друзья
func AcceptFriendRequest(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	var req data.FriendActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec(`CALL public.accept_friend_request($1, $2)`, userID, req.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при принятии запроса", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Запрос принят, вы теперь друзья"})
}

// Отклоняет запрос и добавляет отправителя в подписчики
func RejectFriendRequest(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	var req data.FriendActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec(`CALL public.reject_friend_request($1, $2)`, userID, req.TargetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при отклонении запроса", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Запрос отклонен, пользователь теперь подписчик"})
}

// Получение всех стикер-паков пользователя
func GetUserStickerPacks(c *gin.Context) {
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
	err := db.DB.Get(&result, "SELECT * FROM public.get_user_sticker_packs($1)", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения стикерпаков", "details": err.Error()})
		return
	}

	var stickerPacks data.StickerPackResponse
	if err := json.Unmarshal([]byte(result), &stickerPacks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stickerPacks)
}

// Получение чатов пользователя
func GetUserChats(c *gin.Context) {
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
	err := db.DB.Get(&result, "SELECT * FROM public.get_user_chats($1)", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения чатов", "details": err.Error()})
		return
	}

	var chats data.ChatResponse
	if err := json.Unmarshal([]byte(result), &chats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chats)
}

func GetChatsByFolder(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден"})
		return
	}
	userID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неверный тип user_id"})
		return
	}

	var req data.GetChatsByFolderRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный JSON", "details": err.Error()})
		return
	}

	var result string
	err := db.DB.QueryRow("SELECT public.get_chats_by_folder_and_user($1, $2)", userID, req.FolderID).Scan(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении чатов", "details": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", []byte(result))
}

func GetChatsByFolderAndUser(c *gin.Context) {
	// Получаем user_id
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден"})
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
		return
	}

	// Получаем folder_id из параметров
	folderIDStr := c.Query("folder_id")
	folderID, err := strconv.Atoi(folderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный folder_id"})
		return
	}

	var result string
	// Вызов функции из БД
	err = db.DB.Get(&result, "SELECT * FROM public.get_chats_by_folder_and_user($1, $2)", userID, folderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения чатов", "details": err.Error()})
		return
	}

	var chats data.FolderChatsResponse
	// Декодируем JSON
	if err := json.Unmarshal([]byte(result), &chats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, chats)
}
func GetGroupInfo(c *gin.Context) {
	// Получаем group_id из параметров
	groupIDStr := c.Query("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный group_id"})
		return
	}

	var result string
	// Вызов функции из БД
	err = db.DB.Get(&result, "SELECT * FROM public.get_group_info($1)", groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных о группе", "details": err.Error()})
		return
	}

	var groupInfo data.GroupInfoResponse
	// Парсинг JSON
	if err := json.Unmarshal([]byte(result), &groupInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groupInfo)
}

// GetGroupSubscribers — получение подписчиков группы

func GetGroupSubscribers(c *gin.Context) {
	groupID := c.Param("group_id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group_id обязателен"})
		return
	}

	var result string
	// Выполнение SQL-функции, возвращающей JSON
	err := db.DB.Get(&result, "SELECT * FROM public.get_group_subscribers($1)", groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения подписчиков", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
	var subscribers []data.GroupSubscriber
	// Преобразуем JSON-строку в слайс структур
	if err := json.Unmarshal([]byte(result), &subscribers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscribers)
}

// GetGroupPosts — получение информации о группе и её постах

func GetGroupPostsInfo(c *gin.Context) {
	groupIDStr := c.Param("group_id")
	if groupIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group_id обязателен"})
		return
	}
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group_id должен быть числом"})
		return
	}

	// Считываем limit и offset с параметрами по умолчанию
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var result string
	err = db.DB.Get(&result, "SELECT * FROM public.get_group_posts($1, $2, $3)", groupID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
		return
	}

	var response data.GroupPostsResponse
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка парсинга JSON", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetFeaturesInfo — получает информацию обо всех возможностях

// ⛓ Возвращает список всех возможностей
func GetFeaturesInfo(c *gin.Context) {
	var featuresInfo []data.FeatureInfo

	// Запрос к БД для получения всех возможностей
	err := db.DB.Select(&featuresInfo, "SELECT * FROM public.get_features_info()")
	if err != nil {
		log.Println("Ошибка выполнения запроса get_features_info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить информацию о возможностях", "details": err.Error()})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, featuresInfo)
}

// GetRolesInfoForGroup — получает информацию о ролях в группе
func GetRolesInfoForGroup(c *gin.Context) {
	// Извлекаем group_id из параметров пути
	groupIDStr := c.Param("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат group_id"})
		return
	}

	// Выполнение функции в БД для получения информации о ролях в группе
	var rolesInfo []data.RoleInfo
	err = db.DB.Select(&rolesInfo, "SELECT * FROM public.get_roles_info_for_group($1)", groupID)
	if err != nil {
		log.Println("Ошибка выполнения запроса get_roles_info_for_group:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось получить информацию о ролях", "details": err.Error()})
		return
	}

	// Успешный ответ с информацией о ролях
	c.JSON(http.StatusOK, gin.H{"roles": rolesInfo})
}

// Получение постов группы
func GetGroupPosts(c *gin.Context) {
	// Получение ID группы из параметра пути /group_posts/:group_id
	groupIDStr := c.Param("group_id")
	if groupIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан group_id"})
		return
	}

	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат group_id"})
		return
	}

	// Чтение limit и offset с дефолтами
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	var result string
	err = db.DB.Get(&result, "SELECT * FROM public.get_group_posts($1, $2, $3)", groupID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
		return
	}

	var response data.GroupPostResponse
	if err := json.Unmarshal([]byte(result), &response); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
func GetGroupBlacklist(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден в контексте"})
		return
	}

	// Преобразуем user_id в int
	callerID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
		return
	}

	// Получаем ID группы из URL параметра
	groupIDStr := c.Param("group_id")
	if groupIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан group_id"})
		return
	}

	// Преобразуем group_id в int
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат group_id"})
		return
	}

	// Выполняем запрос к базе данных
	var result string
	err = db.DB.Get(&result, "SELECT public.get_blacklist_for_group($1, $2)", groupID, callerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения ЧС", "details": err.Error()})
		return
	}

	// Парсим ответ в JSON
	var jsonResponse any
	if err := json.Unmarshal([]byte(result), &jsonResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	// Возвращаем ответ клиенту
	c.JSON(http.StatusOK, gin.H{"blacklist": jsonResponse})
}
func GetTags(c *gin.Context) {
	// Получаем ID группы из URL параметра
	groupIDStr := c.Param("group_id")
	if groupIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан group_id"})
		return
	}
	// Преобразуем group_id в int
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат group_id"})
		return
	}
	// Выполним запрос к функции для получения тегов
	var tags []data.TagsGroupRequest // data.Tag - структура, которая представляет тег в Go
	err = db.DB.Select(&tags, `SELECT * FROM public.get_tags($1)`, groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения тегов", "details": err.Error()})
		return
	}

	// Ответ с полученными тегами
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}
func GetComments(c *gin.Context) {
	// Получаем параметры из запроса
	postIDStr := c.DefaultQuery("post_id", "")
	typeStr := c.DefaultQuery("type", "user")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	// Преобразуем параметры в нужные типы
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат post_id"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат limit"})
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат offset"})
		return
	}

	// Вызываем функцию для получения комментариев
	var commentsJSON string
	err = db.DB.QueryRow(`
		SELECT public.get_comments($1, $2, $3, $4)`,
		postID, typeStr, limit, offset).Scan(&commentsJSON)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении комментариев", "details": err.Error()})
		return
	}

	// Возвращаем результат пользователю
	c.JSON(http.StatusOK, gin.H{"comments": commentsJSON})
}
func GetUserFavouritePosts(c *gin.Context) {
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

	var jsonResult string
	err := db.DB.QueryRow(`SELECT * FROM public.get_user_favourite_posts($1)`, userID).Scan(&jsonResult)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
		return
	}

	var favourites data.FavouritePost
	if err := json.Unmarshal([]byte(jsonResult), &favourites); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, favourites)
}
func GetUserFavouriteSMS(c *gin.Context) {
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

	var results []data.FavouriteSMS
	query := `SELECT * FROM public.get_user_favourite_sms($1)`
	err := db.DB.Select(&results, query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка получения избранных сообщений",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, results)
}

// ..........................................................................
//
// .......................................................................................

func SearchGroups(c *gin.Context) {
	var req data.SearchRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON", "details": err.Error()})
		return
	}

	if req.Search == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан параметр поиска"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный параметр limit"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный параметр offset"})
		return
	}

	var rawResults []string
	query := `SELECT * FROM public.search_groups($1, $2, $3)`
	err = db.DB.Select(&rawResults, query, req.Search, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка выполнения поиска",
			"details": err.Error(),
		})
		return
	}

	var parsedResults []map[string]interface{}

	for _, raw := range rawResults {
		var item map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &item); err != nil {
			continue // или логировать
		}

		// Обработка base64 картинок, если есть
		if val, ok := item["profile_picture"].(string); ok && val != "" {
			item["profile_picture"] = "data:image/png;base64," + val
		}

		parsedResults = append(parsedResults, item)
	}

	c.JSON(http.StatusOK, parsedResults)
}

func SearchPosts(c *gin.Context) {
	var req data.SearchRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON", "details": err.Error()})
		return
	}

	if req.Search == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан параметр поиска"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный параметр limit"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный параметр offset"})
		return
	}

	var rawResults []string
	query := `SELECT * FROM public.search_posts($1, $2, $3)`
	err = db.DB.Select(&rawResults, query, req.Search, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка выполнения поиска",
			"details": err.Error(),
		})
		return
	}

	var parsedResults []map[string]interface{}

	for _, raw := range rawResults {
		var item map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &item); err != nil {
			continue // или логировать
		}

		// Обработка base64 картинок, если есть
		if val, ok := item["file_post"].(string); ok && val != "" {
			item["file_post"] = "data:image/png;base64," + val
		}

		parsedResults = append(parsedResults, item)
	}

	c.JSON(http.StatusOK, parsedResults)
}

func SearchUsers(c *gin.Context) {
	var req data.SearchRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON", "details": err.Error()})
		return
	}

	if req.Search == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан параметр поиска"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный параметр limit"})
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный параметр offset"})
		return
	}

	var rawResults []string
	query := `SELECT * FROM public.search_users($1, $2, $3)`
	err = db.DB.Select(&rawResults, query, req.Search, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка выполнения поиска",
			"details": err.Error(),
		})
		return
	}

	var parsedResults []map[string]interface{}

	for _, raw := range rawResults {
		var item map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &item); err != nil {
			continue // или логировать
		}

		// Обработка base64 картинок, если есть
		if val, ok := item["profile_picture"].(string); ok && val != "" {
			item["profile_picture"] = "data:image/png;base64," + val
		}

		parsedResults = append(parsedResults, item)
	}

	c.JSON(http.StatusOK, parsedResults)
}

func SearchAll(c *gin.Context) {
	var req data.SearchRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON", "details": err.Error()})
		return
	}

	if req.Search == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан параметр поиска"})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный параметр limit"})
		return
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный параметр offset"})
		return
	}

	var result string
	query := `SELECT jsonb_agg(result_item)::text FROM public.search_all($1, $2, $3)`
	err = db.DB.Get(&result, query, req.Search, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка выполнения поиска",
			"details": err.Error(),
		})
		return
	}

	var results []data.SearchResult
	if err := json.Unmarshal([]byte(result), &results); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	// Обработка base64 картинок
	for i := range results {
		if results[i].ProfilePicture != "" {
			results[i].ProfilePicture = "data:image/png;base64," + results[i].ProfilePicture
		}
		if results[i].FilePost != "" {
			results[i].FilePost = "data:image/png;base64," + results[i].FilePost
		}
	}

	c.JSON(http.StatusOK, results)
}

// func GetFilteredPosts(c *gin.Context) {
// 	userIDRaw, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден в контексте"})
// 		return
// 	}

// 	userID, ok := userIDRaw.(int)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
// 		return
// 	}

// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
// 	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

// 	var result string
// 	err := db.DB.Get(&result, "SELECT * FROM public.get_filtered_posts($1, $2, $3)", userID, limit, offset)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
// 		return
// 	}

// 	var responseData []data.FilteredPost
// 	if err := json.Unmarshal([]byte(result), &responseData); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, responseData)
// }

// func GetRecommendedPosts(c *gin.Context) {
// 	userIDRaw, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден в контексте"})
// 		return
// 	}

// 	userID, ok := userIDRaw.(int)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
// 		return
// 	}

// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
// 	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

// 	var result string
// 	err := db.DB.Get(&result, "SELECT * FROM public.get_recommended_posts($1, $2, $3)", userID, limit, offset)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
// 		return
// 	}

// 	var responseData []data.RecommendedPost
// 	if err := json.Unmarshal([]byte(result), &responseData); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
// 		return
// 	}

//		c.JSON(http.StatusOK, responseData)
//	}

// func GetFilteredPosts(c *gin.Context) {
// 	// Получаем user_id из контекста
// 	userIDRaw, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден в контексте"})
// 		return
// 	}

// 	userID, ok := userIDRaw.(int)
// 	if !ok {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
// 		return
// 	}

// 	// Получаем параметры limit и offset из запроса
// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
// 	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

// 	// Получаем JSON как строку из функции БД
// 	// Получаем данные постов из базы данных
// 	var result string
// 	err := db.DB.Get(&result, "SELECT * FROM public.get_filtered_posts($1, $2, $3)", userID, limit, offset)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
// 		return
// 	}

// 	// Разбираем результат на структуру
// 	var responseData []data.FilteredPost
// 	if err := json.Unmarshal([]byte(result), &responseData); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
// 		return
// 	}

// 	// Преобразуем поле FalePost (если оно содержит изображение) в Base64
// 	for i := range responseData {
// 		if len(responseData[i].FalePost) > 0 {
// 			// Преобразуем изображение в строку Base64
// 			responseData[i].FalePost = "data:image/png;base64," + base64.StdEncoding.EncodeToString([]byte(responseData[i].FalePost))
// 		}
// 	}

// 	c.JSON(http.StatusOK, responseData)

// }
func GetRecommendedPosts(c *gin.Context) {
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

	var result string
	err := db.DB.Get(&result, "SELECT * FROM public.get_recommended_posts($1, $2, $3)", userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
		return
	}

	var rawPosts []data.RawFilteredPost
	if err := json.Unmarshal([]byte(result), &rawPosts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	var finalPosts []data.FilteredPostResponse
	for _, p := range rawPosts {
		finalPosts = append(finalPosts, data.FilteredPostResponse{
			PostType:           p.PostType,
			PostID:             p.PostID,
			UserID:             p.UserID,
			Header:             p.Header,
			Text:               p.Text,
			FalePostBase64:     "data:image/png;base64," + base64.StdEncoding.EncodeToString(p.FalePost),
			DateTimePost:       p.DateTimePost,
			ViewsPost:          p.ViewsPost,
			Repost:             p.Repost,
			CommentsPermission: p.CommentsPermission,
			CommentsCount:      p.CommentsCount,
			GroupInfo:          p.GroupInfo,
			Author:             p.Author,
			LikesCount:         p.LikesCount,
		})
	}

	c.JSON(http.StatusOK, finalPosts)
}
func GetFilteredPosts(c *gin.Context) {
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

	var result string
	err := db.DB.Get(&result, "SELECT * FROM public.get_filtered_posts($1, $2, $3)", userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
		return
	}

	var rawPosts []data.RawFilteredPost
	if err := json.Unmarshal([]byte(result), &rawPosts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	var finalPosts []data.FilteredPostResponse
	for _, p := range rawPosts {
		finalPosts = append(finalPosts, data.FilteredPostResponse{
			PostType:           p.PostType,
			PostID:             p.PostID,
			UserID:             p.UserID,
			Header:             p.Header,
			Text:               p.Text,
			FalePostBase64:     "data:image/png;base64," + base64.StdEncoding.EncodeToString(p.FalePost),
			DateTimePost:       p.DateTimePost,
			ViewsPost:          p.ViewsPost,
			Repost:             p.Repost,
			CommentsPermission: p.CommentsPermission,
			CommentsCount:      p.CommentsCount,
			GroupInfo:          p.GroupInfo,
			Author:             p.Author,
			LikesCount:         p.LikesCount,
		})
	}

	c.JSON(http.StatusOK, finalPosts)
}

type ChatFolderBool struct {
	IDChat        int     `json:"id_chat"`
	NameChat      string  `json:"name_chat"`
	CountChatPepl string  `json:"countChatPepl"`
	Pfoto         *string `json:"pfoto"` // Здесь используем указатель на строку, чтобы поддерживать null
	DateTimeChat  string  `json:"dateTime_chat"`
	IsInFolder    bool    `json:"is_in_folder"`
	LastMessage   struct {
		TextSMS        string  `json:"text_sms"`
		DateTimeSMS    string  `json:"dateTime_sms"`
		UserID         int     `json:"user_id"`
		Username       string  `json:"username"`
		ProfilePicture *string `json:"profile_picture"` // Используем указатель для возможности null
	} `json:"last_message"`
}

// GetUserChatFolders - обработчик запросов для получения чатов пользователя по папке
// GetUserChatFolders - обработчик запросов для получения чатов пользователя по папке
func GetUserChatFoldersBoolen(c *gin.Context) {
	// Получаем параметр папки из URL
	idChatFolderRaw := c.Param("folder")
	if idChatFolderRaw == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Folder is required"})
		return
	}

	// Преобразуем строку в целое число
	idChatFolder, err := strconv.Atoi(idChatFolderRaw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	// Получаем userID из контекста
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user_id type"})
		return
	}

	// Параметры пагинации
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Логируем параметры перед запросом для отладки
	log.Printf("Fetching chats for userID: %d, folderID: %d, limit: %d, offset: %d", userID, idChatFolder, limit, offset)

	// Выполняем запрос к базе данных
	var result string
	err = db.DB.Get(&result, "SELECT * FROM public.get_user_chats_folderAdd($1, $2, $3, $4)", userID, idChatFolder, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching chats", "details": err.Error()})
		return
	}

	// Преобразуем результат в структуру
	var responseData []ChatFolderBool
	if err := json.Unmarshal([]byte(result), &responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing JSON", "details": err.Error()})
		return
	}

	// Возвращаем результат
	c.JSON(http.StatusOK, responseData)
}

func GetChatDetails(c *gin.Context) {
	chatID := c.Param("chat_id")
	if chatID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "chat_id обязателен"})
		return
	}

	var result string
	err := db.DB.Get(&result, "SELECT public.get_chat_details($1)", chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных о чате", "details": err.Error()})
		return
	}

	var chatDetails data.ChatDetailsResponse
	if err := json.Unmarshal([]byte(result), &chatDetails); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	// Обработка фото чата
	if chatDetails.ChatInfo.PFoto != "" {
		chatDetails.ChatInfo.PFoto = "data:image/png;base64," + chatDetails.ChatInfo.PFoto
	}

	// Обработка фото пользователей
	for i := range chatDetails.UsersInfo {
		if chatDetails.UsersInfo[i].ProfilePicture != "" {
			chatDetails.UsersInfo[i].ProfilePicture = "data:image/png;base64," + chatDetails.UsersInfo[i].ProfilePicture
		}
	}

	c.JSON(http.StatusOK, chatDetails)
}
func AreFriends(c *gin.Context) {
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

	otherIDParam := c.Param("id") // URL: /friends/check/:id
	otherID, err := strconv.Atoi(otherIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID второго пользователя"})
		return
	}

	var areFriends bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM friends 
			WHERE (user_id = $1 AND user_id_f = $2) 
			   OR (user_id = $2 AND user_id_f = $1)
		);
	`

	err = db.DB.Get(&areFriends, query, userID, otherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки дружбы", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"are_friends": areFriends})
}
