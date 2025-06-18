package handlers

import (
	"encoding/base64"
	"log"
	"net/http"
	"serverGo/data"
	"serverGo/db"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// Получает список друзей
func GetFriendsList(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	rows, err := db.DB.Query(`SELECT * FROM public.get_friends_list($1)`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении списка друзей", "details": err.Error()})
		return
	}
	defer rows.Close()

	var friends []map[string]interface{}
	for rows.Next() {
		var id int
		var login, avatar string
		rows.Scan(&id, &login, &avatar)
		friends = append(friends, gin.H{
			"id":     id,
			"login":  login,
			"avatar": avatar,
		})
	}

	c.JSON(http.StatusOK, gin.H{"friends": friends})
}

// Получает входящие запросы в друзья
func GetFriendRequests(c *gin.Context) {
	userID := c.MustGet("user_id").(int)

	rows, err := db.DB.Query(`SELECT * FROM public.get_friend_requests($1)`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении запросов", "details": err.Error()})
		return
	}
	defer rows.Close()

	var requests []map[string]interface{}
	for rows.Next() {
		var id int
		var login, avatar string
		rows.Scan(&id, &login, &avatar)
		requests = append(requests, gin.H{
			"id":     id,
			"login":  login,
			"avatar": avatar,
		})
	}

	c.JSON(http.StatusOK, gin.H{"requests": requests})
}

// GetSubscribersList возвращает список подписчиков пользователя
func GetSubscribersList(c *gin.Context) {
	// Получаем user_id из контекста (предполагается, что он устанавливается в middleware аутентификации)
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

	var subscribers []data.Subscriber

	err := db.DB.Select(&subscribers, "SELECT * FROM public.get_subscribers_list($1)", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Ошибка получения списка подписчиков",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, subscribers)
}
func RemoveSubscriber(c *gin.Context) {
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

	var req data.RemoveSubscriberRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec("CALL public.remove_subscriber($1, $2)", userID, req.SubscriberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления подписчика", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Подписчик удалён"})
}

func UnsubscribeFromUser(c *gin.Context) {
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

	var req data.UnsubscribeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec("CALL public.unsubscribe_from_user($1, $2)", userID, req.TargetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка отписки", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Вы отписались от пользователя"})
}

func DeletePost(c *gin.Context) {
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

	var req data.DeletePostRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec("CALL public.delete_post_user($1, $2)", req.PostID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления поста", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Пост удален"})
}

func UpdatePost(c *gin.Context) {
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

	var req data.UpdatePostRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	var fileBytes []byte
	if req.FileBase64 != nil {
		var err error
		fileBytes, err = base64.StdEncoding.DecodeString(*req.FileBase64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка декодирования файла", "details": err.Error()})
			return
		}
	}

	_, err := db.DB.Exec("CALL public.update_post_user($1, $2, $3, $4, $5)",
		req.PostID, userID, req.Header, req.Text, fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления поста", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Пост обновлен"})
}

// Создание нового стикер-пака
func CreateStickerPack(c *gin.Context) {
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

	var req data.CreateStickerPackRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec("CALL public.create_sticker_pack($1, $2)", userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания стикер-пака", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Стикер-пак создан"})
}

// Добавление нового стикера в существующий стикер-пак
func AddStickerToPack(c *gin.Context) {
	var req data.AddStickerToPackRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	fileBytes, err := base64.StdEncoding.DecodeString(req.FileBase64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка декодирования файла", "details": err.Error()})
		return
	}

	_, err = db.DB.Exec("CALL public.add_sticker_to_pack($1, $2)", req.StickerPackID, fileBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления стикера", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Стикер добавлен"})
}

// Создание группового чата
func CreateGroupChat(c *gin.Context) {
	var req data.CreateGroupChatRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	var photoBytes []byte
	if req.PhotoBase64 != nil {
		var err error
		photoBytes, err = base64.StdEncoding.DecodeString(*req.PhotoBase64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка декодирования фото", "details": err.Error()})
			return
		}
	}

	// Вызов процедуры для создания группового чата
	_, err := db.DB.Exec("CALL public.create_group_chat($1, $2, $3, $4)",
		req.CreatorID, req.ChatName, req.UserIDs, photoBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания чата", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Групповой чат создан"})
}

// Отправка стикера пользователю
func SendStickerToUser(c *gin.Context) {
	var req data.SendStickerToUserRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	// Вызов процедуры для отправки стикера пользователю
	_, err := db.DB.Exec("CALL public.send_sticker_to_user($1, $2, $3)",
		req.SenderID, req.ReceiverID, req.StickerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка отправки стикера", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Стикер отправлен"})
}

// Отправка стикера в чат
func SendStickerToChat(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден в контексте"})
		return
	}
	senderID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
		return
	}

	var req data.SendStickerRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	// Вызов процедуры в базе
	_, err := db.DB.Exec("CALL public.send_sticker_to_chat($1, $2, $3)", senderID, req.ChatID, req.StickerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка отправки стикера", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Стикер отправлен"})
}

// Удаление сообщения из чата
func DeleteSMSFromChat(c *gin.Context) {
	var req data.DeleteSMSRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	// Удаляем сообщение
	_, err := db.DB.Exec("CALL public.delete_sms_from_chat($1)", req.SmsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления сообщения", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Сообщение удалено"})
}

// RemoveChatFromFolder удаляет чат из указанной папки пользователя
func RemoveChatFromFolder(c *gin.Context) {
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

	var req data.ChatFolderActionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec("CALL public.remove_chat_from_folder($1, $2, $3)", req.ChatFolderID, req.ChatID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления чата из папки", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Чат успешно удален из папки"})
}

// RemoveChatFolder удаляет папку чатов пользователя вместе со всеми её чатами
func RemoveChatFolder(c *gin.Context) {
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

	var req data.RemoveChatFolderRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec("CALL public.remove_chat_folder($1, $2)", req.ChatFolderID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления папки", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Папка успешно удалена"})
}

// Обработчик для добавления группы
func AddGroup(c *gin.Context) {
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

	var req data.AddGroupRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	// Декодируем фото, если оно передано
	var photoBytes []byte
	if req.PhotoBase64 != nil {
		var err error
		photoBytes, err = base64.StdEncoding.DecodeString(*req.PhotoBase64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка декодирования фото", "details": err.Error()})
			return
		}
	}

	// Вызываем хранимую процедуру add_group
	_, err := db.DB.Exec("CALL public.add_group($1, $2, $3, $4, $5)", userID, req.Name, req.TypeGroupID, req.Description, photoBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления группы", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Группа успешно добавлена"})
}

// AddFeature — добавление новой возможности для роли

func AddFeature(c *gin.Context) {
	var featureRequest data.FeatureRequest

	// Чтение данных из тела запроса
	if err := c.ShouldBindJSON(&featureRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Выполнение SQL-процедуры добавления новой возможности
	_, err := db.DB.Exec("CALL public.add_feature($1, $2, $3)", featureRequest.AdminID, featureRequest.NameFeatureRole, featureRequest.DescriptionFeature)
	if err != nil {
		log.Println("Ошибка выполнения процедуры add_feature:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось добавить возможность", "details": err.Error()})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Возможность успешно добавлена"})
}

// CreateRoleWithFeatures — создает роль с возможностями для группы

// ⛓ Принимает тело запроса с данными роли, проверяет права пользователя
func CreateRoleWithFeatures(c *gin.Context) {
	// Извлекаем group_id из параметров пути
	groupIDStr := c.Param("group_id")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат group_id"})
		return
	}

	// Извлекаем данные из тела запроса
	var request data.CreateRoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Получаем ID пользователя, который вызывает процедуру
	callerIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден в контексте"})
		return
	}
	callerID, ok := callerIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
		return
	}

	// Выполнение процедуры в БД
	_, err = db.DB.Exec(
		"CALL public.create_role_group_with_features($1, $2, $3, $4)",
		groupID,
		request.RoleName,
		pq.Array(request.FeatureIDs), // Передача массива с использованием pq.Array
		callerID,
	)
	if err != nil {
		log.Println("Ошибка выполнения процедуры create_role_group_with_features:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать роль", "details": err.Error()})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Роль успешно создана"})
}

// AddFeatureToRole — добавляет возможность в роль группы
func AddFeatureToRole(c *gin.Context) {
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

	// Извлекаем role_id из параметров пути
	roleIDStr := c.Param("role_id")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат role_id"})
		return
	}

	// Извлекаем данные из тела запроса
	var request data.AddFeatureRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Выполнение процедуры в БД
	_, err = db.DB.Exec("CALL public.add_feature_to_role($1, $2, $3)", roleID, request.FeatureID, userID)
	if err != nil {
		log.Println("Ошибка выполнения процедуры add_feature_to_role:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось добавить возможность в роль", "details": err.Error()})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Возможность успешно добавлена в роль"})
}

// AddUserToRoleGroup — добавляет пользователя в роль группы
func AddUserToRoleGroup(c *gin.Context) {
	// Извлекаем role_id из параметров пути
	roleIDStr := c.Param("role_id")
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат role_id"})
		return
	}

	// Извлекаем данные из тела запроса
	var request data.AddUserToRoleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Получаем ID пользователя, который вызывает процедуру
	callerIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден в контексте"})
		return
	}
	callerID, ok := callerIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
		return
	}

	// Выполнение процедуры в БД
	_, err = db.DB.Exec("CALL public.add_user_to_role_group($1, $2, $3)", roleID, request.UserID, callerID)
	if err != nil {
		log.Println("Ошибка выполнения процедуры add_user_to_role_group:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось добавить пользователя в роль", "details": err.Error()})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно добавлен в роль группы"})
}
func AddCommentToGroupPost(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден в контексте"})
		return
	}

	// Преобразуем user_id в int
	userID, ok := userIDRaw.(int)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Неверный тип user_id"})
		return
	}

	// Получаем ID группы из параметра пути
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

	// Получаем данные комментария из тела запроса
	var req data.AddCommentToGroupPostRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	// Если parent_comment_id не передан, передаем nil
	var parentID interface{} = nil
	if req.ParentCommentID != nil {
		parentID = *req.ParentCommentID
	}

	// Вставляем комментарий
	_, err = db.DB.Exec(`CALL public.add_comment_to_group_post($1, $2, $3, $4)`,
		groupID, userID, req.CommentText, parentID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления комментария", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Комментарий добавлен"})
}
func AddUserToGroupBlacklist(c *gin.Context) {
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

	// Получаем ID группы из параметра пути
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

	// Получаем данные для добавления пользователя в черный список
	var req data.AddUserToGroupBlacklistRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	// Выполняем запрос для добавления пользователя в черный список
	_, err = db.DB.Exec(`CALL public.add_user_to_group_blacklist($1, $2, $3)`,
		groupID, req.UserID, callerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления в ЧС", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Пользователь добавлен в ЧС группы"})
}
func RemoveUserFromGroupBlacklist(c *gin.Context) {
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

	// Получаем ID группы из параметра пути
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

	// Получаем ID пользователя, которого нужно удалить из черного списка
	var req data.RemoveUserFromGroupBlacklistRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	// Выполняем запрос для удаления пользователя из черного списка
	_, err = db.DB.Exec(`CALL public.remove_user_from_group_blacklist($1, $2, $3)`,
		req.UserID, groupID, callerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления из ЧС", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Пользователь удален из ЧС группы"})
}
func CreateTag(c *gin.Context) {
	// Получаем user_id из контекста (он должен быть передан авторизованным пользователем)
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

	// Прочитаем тело запроса для данных о теге
	var req data.TagsGroupRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	// Вызовем процедуру для создания тега
	_, err := db.DB.Exec(`CALL public.create_tag($1, $2, $3)`, req.TagName, req.TagDescription, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания тега", "details": err.Error()})
		return
	}

	// Ответ успешного создания тега
	c.JSON(http.StatusOK, gin.H{"status": "Тег успешно создан"})
}
func UpdateTagsInGroup(c *gin.Context) {
	// Получаем user_id из контекста (он должен быть передан авторизованным пользователем)
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

	// Прочитаем данные из тела запроса
	var req data.UpdateTagGroup
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	// Вызовем процедуру для обновления тегов
	_, err := db.DB.Exec(`CALL public.update_tags_in_group($1, $2, $3)`,
		req.GroupID, pq.Array(req.Tags), userID) // Используем pq.Array для передачи массива
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления тегов", "details": err.Error()})
		return
	}

	// Ответ успешного обновления
	c.JSON(http.StatusOK, gin.H{"status": "Теги успешно обновлены"})
}
func ToggleLikeComment(c *gin.Context) {
	// Получаем параметры из запроса
	commentIDStr := c.DefaultQuery("comment_id", "")
	userIDStr := c.DefaultQuery("user_id", "")
	isGroupStr := c.DefaultQuery("is_group", "false")

	// Преобразуем параметры в нужные типы
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат comment_id"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат user_id"})
		return
	}

	// Проверяем is_group, если не передано, принимаем как false
	isGroup := false
	if isGroupStr == "true" {
		isGroup = true
	}

	// Выполняем запрос к БД для вызова процедуры
	_, err = db.DB.Exec(`CALL public.toggle_like_comment($1, $2, $3)`,
		commentID, userID, isGroup)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке лайка", "details": err.Error()})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusOK, gin.H{"status": "Лайк успешно обработан"})
}
func AddToFavourites(c *gin.Context) {
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

	var req data.AddToFavouritesRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	_, err := db.DB.Exec(`CALL public.add_to_favourites($1, $2, $3)`,
		userID, req.ItemID, req.FavType)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка добавления в избранное", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Добавлено в избранное"})
}

func UpdateChatInfo(c *gin.Context) {
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден в контексте"})
		return
	}

	var req data.UpdateChatInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON", "details": err.Error()})
		return
	}

	// Обработка изображения: декод base64
	var imageBytes []byte
	if req.Pfoto != "" {
		if commaIdx := strings.Index(req.Pfoto, ","); commaIdx != -1 {
			req.Pfoto = req.Pfoto[commaIdx+1:]
		}

		decoded, err := base64.StdEncoding.DecodeString(req.Pfoto)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка декодирования изображения", "details": err.Error()})
			return
		}
		imageBytes = decoded
	} else {
		imageBytes = nil // null передаём в SQL
	}

	_, err := db.DB.Exec("SELECT update_chat_info($1, $2, $3)", req.ChatID, req.NameChat, imageBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления чата", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Чат успешно обновлён"})
}
