package handlers

import (
	"database/sql"
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

type RequestBody struct {
	GroupID int `json:"group_id"`
}

// GroupInfoResponse - структура для ответа
type GroupInfoResponse struct {
	IDGroup      int    `json:"id_group"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Access       string `json:"access"`
	MembersCount int    `json:"members_count"`
	PostsCount   int    `json:"posts_count"`
	Owner        struct {
		UserID         int    `json:"user_id"`
		Username       string `json:"username"`
		ProfilePicture string `json:"profile_picture"`
	} `json:"owner"`
	Tags []struct {
		NameTag        string `json:"name_tag"`
		IDTag          int    `json:"id_tag"`
		DescriptionTag string `json:"description_tag"`
	} `json:"tags"`
}

// GetGroupInfo - обработчик запроса для получения информации о группе
func GetGroupInfo(c *gin.Context) {
	var requestBody RequestBody

	// Парсим тело запроса в структуру RequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении данных", "details": err.Error()})
		return
	}

	// Извлекаем group_id из тела запроса
	groupID := requestBody.GroupID

	// Запрашиваем информацию о группе из базы данных
	var result string
	err := db.DB.Get(&result, "SELECT * FROM public.get_group_info($1)", groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных о группе", "details": err.Error()})
		return
	}

	// Преобразуем результат в структуру GroupInfoResponse
	var groupInfo GroupInfoResponse
	if err := json.Unmarshal([]byte(result), &groupInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	// Отправляем успешный ответ
	c.JSON(http.StatusOK, groupInfo)
}

// GetGroupSubscribers — получение подписчиков группы

func GetGroupSubscribers(c *gin.Context) {
	groupID := c.Param("group_id")
	if groupID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group_id обязателен"})
		return
	}

	var jsonBytes []byte
	err := db.DB.Get(&jsonBytes, "SELECT public.get_group_subscribers($1)", groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка получения подписчиков",
			"details": err.Error(),
		})
		return
	}

	var subscribers []data.GroupSubscriber
	if err := json.Unmarshal(jsonBytes, &subscribers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка обработки JSON",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, subscribers)
}

// GetGroupPosts — получение информации о группе и её постах

// func GetGroupPostsInfo(c *gin.Context) {
// 	groupIDStr := c.Param("group_id")
// 	if groupIDStr == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "group_id обязателен"})
// 		return
// 	}
// 	groupID, err := strconv.Atoi(groupIDStr)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "group_id должен быть числом"})
// 		return
// 	}

// 	// Считываем limit и offset с параметрами по умолчанию
// 	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
// 	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

// 	var result string
// 	err = db.DB.Get(&result, "SELECT * FROM public.get_group_posts($1, $2, $3)", groupID, limit, offset)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
// 		return
// 	}

// 	var response data.GroupPostsResponse
// 	if err := json.Unmarshal([]byte(result), &response); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка парсинга JSON", "details": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// GetFeaturesInfo — получает информацию обо всех возможностях

// ⛓ Возвращает список всех возможностей
func GetFeaturesInfo(c *gin.Context) {
	var features []data.FeatureInfo

	err := db.DB.Select(&features, "SELECT * FROM public.get_features_info()")
	if err != nil {
		log.Println("Ошибка при запросе features info:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка запроса к базе",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"features": features})
}

// GetRolesInfoForGroup — получает информацию о ролях в группе
func GetRolesInfoForGroup(c *gin.Context) {
	groupIDStr := c.Param("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат group_id"})
		return
	}

	var rawJSON []byte
	err = db.DB.Get(&rawJSON, "SELECT public.get_roles_info_for_group($1)", groupID)
	if err != nil {
		log.Println("Ошибка выполнения функции get_roles_info_for_group:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных", "details": err.Error()})
		return
	}

	var rolesInfo []data.RoleInfo
	err = json.Unmarshal(rawJSON, &rolesInfo)
	if err != nil {
		log.Println("Ошибка парсинга JSON в RoleInfo:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки данных", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": rolesInfo})
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

func GetAllTags(c *gin.Context) {
	var tags []data.Tag_

	err := db.DB.Select(&tags, `SELECT id_tag, name_tag, description_tag FROM public.tags`)
	if err != nil {
		log.Println("Ошибка при получении тегов:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Не удалось получить теги",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tags)
}
func GetTags_admin(c *gin.Context) {
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
		if val, ok := item["group_photo_base64"].(string); ok && val != "" {
			item["group_photo_base64"] = "data:image/png;base64," + val
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
		// Профиль пользователя
		if results[i].ProfilePicture != "" {
			results[i].ProfilePicture = "data:image/png;base64," + results[i].ProfilePicture
		}

		// Файл поста
		if results[i].FilePost != "" {
			results[i].FilePost = "data:image/png;base64," + results[i].FilePost
		}

		// Фото группы (если это группа)
		if results[i].Type == "group" && results[i].Photo != "" {
			results[i].Photo = "data:image/png;base64," + results[i].Photo
		}

		// Фото группы внутри поста (post_group)
		if results[i].Group != nil && results[i].Group.Photo != "" {
			results[i].Group.Photo = "data:image/png;base64," + results[i].Group.Photo
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

func CheckFriendRequest(c *gin.Context) {
	currentUserID := c.GetInt("user_id") // Получи ID из middleware или токена
	otherUserIDParam := c.Param("other_user_id")

	otherUserID, err := strconv.Atoi(otherUserIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID пользователя"})
		return
	}

	var status string
	query := `
		SELECT CASE
			WHEN EXISTS (
				SELECT 1 FROM friendrequest 
				WHERE user_id_r = $1 AND user_id_f_r = $2
			) THEN 'me'
			WHEN EXISTS (
				SELECT 1 FROM friendrequest 
				WHERE user_id_r = $2 AND user_id_f_r = $1
			) THEN 'he'
			ELSE 'false'
		END AS status;
	`

	err = db.DB.Get(&status, query, currentUserID, otherUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка запроса", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}
func CheckSubscription(c *gin.Context) {
	currentUserID := c.GetInt("user_id") // ID текущего пользователя из middleware или токена
	otherUserIDParam := c.Param("other_user_id")

	otherUserID, err := strconv.Atoi(otherUserIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID пользователя"})
		return
	}

	var status string
	query := `
		SELECT CASE
			WHEN EXISTS (
				SELECT 1 FROM subscribers 
				WHERE user_id_s = $1 AND user_id_f_s = $2
			) THEN 'me'
			WHEN EXISTS (
				SELECT 1 FROM subscribers 
				WHERE user_id_s = $2 AND user_id_f_s = $1
			) THEN 'he'
			ELSE 'false'
		END AS status;
	`

	err = db.DB.Get(&status, query, currentUserID, otherUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка проверки подписки", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}

func GetGroupsByUserID(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID пользователя обязателен"})
		return
	}

	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID пользователя должен быть числом"})
		return
	}

	var result string
	err = db.DB.Get(&result, "SELECT json_agg(t) FROM get_user_groups($1) t", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении групп по ID", "details": err.Error()})
		return
	}

	var groups []data.GroupResponse
	if err := json.Unmarshal([]byte(result), &groups); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка парсинга JSON", "details": err.Error()})
		return
	}

	for i := range groups {
		if groups[i].ProfilePictureBase64 != "" {
			groups[i].ProfilePictureBase64 = "data:image/png;base64," + groups[i].ProfilePictureBase64
		}
	}

	c.JSON(http.StatusOK, groups)
}

func GetGroupPosts(c *gin.Context) {
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

	limit, err1 := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, err2 := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные limit или offset"})
		return
	}

	var result string
	err = db.DB.Get(&result, "SELECT * FROM public.get_group_posts($1, $2, $3)", groupID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения постов", "details": err.Error()})
		return
	}

	var raw data.RawGroupPostsResponse
	if err := json.Unmarshal([]byte(result), &raw); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	for i, post := range raw.Posts {
		if post.ImageBase64 != "" {
			raw.Posts[i].ImageBase64 = "data:image/png;base64," + post.ImageBase64
		}
	}

	c.JSON(http.StatusOK, raw)
}

func CheckChatExistence(c *gin.Context) {
	currentUserID := c.GetInt("user_id") // ID текущего пользователя (получен из middleware)
	otherUserIDParam := c.Param("other_user_id")

	// Преобразуем параметр в int
	otherUserID, err := strconv.Atoi(otherUserIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID пользователя"})
		return
	}

	var chatID int
	query := `
		SELECT chat_id
		FROM userchat
		WHERE user_ch_id = $1
		AND chat_id IN (
			SELECT chat_id
			FROM userchat
			WHERE user_ch_id = $2
		)
		LIMIT 1
	`

	err = db.DB.Get(&chatID, query, currentUserID, otherUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusOK, gin.H{"exists": false})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке чата", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exists": true,
		"chatID": chatID,
	})
}

type GroupRequest struct {
	GroupID int `json:"group_id"`
}

func GetUserRolesInGroup(c *gin.Context) {
	var req GroupRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.GroupID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON или отсутствует group_id"})
		return
	}

	currentUserID := c.GetInt("user_id") // получен из middleware

	var result string
	err := db.DB.Get(&result, "SELECT get_user_roles_and_features_json($1, $2)", currentUserID, req.GroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных", "details": err.Error()})
		return
	}

	var responseData data.UserRolesAndFeaturesResponse
	if err := json.Unmarshal([]byte(result), &responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обработки JSON", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, responseData)
}
func GetGroupsForAuthorizedUser(c *gin.Context) {
	// Получаем user_id из контекста
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неавторизованный доступ"})
		return
	}

	// Преобразуем user_id в int
	userID, ok := userIDVal.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка извлечения ID пользователя"})
		return
	}

	// Выполняем SQL-запрос для получения групп пользователя
	var result string
	err := db.DB.Get(&result, "SELECT json_agg(t) FROM get_user_groups($1) t", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении групп", "details": err.Error()})
		return
	}

	// Преобразуем результат в список групп
	var groups []data.GroupResponse
	if err := json.Unmarshal([]byte(result), &groups); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка парсинга JSON", "details": err.Error()})
		return
	}

	// Добавляем префикс base64 для фото группы и профиля
	for i := range groups {
		// Для фото владельца группы (profile_picture)
		if groups[i].ProfilePictureBase64 != "" {
			groups[i].ProfilePictureBase64 = "data:image/png;base64," + groups[i].ProfilePictureBase64
		}
		// Для фото группы (group_photo)
		if groups[i].GroupPhotoBase64 != "" {
			groups[i].GroupPhotoBase64 = "data:image/png;base64," + groups[i].GroupPhotoBase64
		}
	}

	// Отправляем результат
	c.JSON(http.StatusOK, groups)
}
