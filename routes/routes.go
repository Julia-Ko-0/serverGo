package routes

import (
	"serverGo/handlers"
	"serverGo/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Маршруты для авторизации
	r.POST("/login", handlers.LoginUser)
	r.POST("/logout", handlers.LogoutUser)
	r.POST("/register", handlers.RegisterUser)
	r.POST("/check-auth", handlers.CheckAuth)

	// Защищённые маршруты
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{

		// Получение информации о пользователе
		auth.POST("/user_info", handlers.GetUserInfo)
		auth.POST("/users_info/:login", handlers.GetUsersInfo)

		// Получение постов пользователя
		auth.POST("/user_info_post", handlers.GetUserPosts)
		auth.POST("/user_info_post_/:login", handlers.GetUsersPosts)

		// Получение чатов пользователя
		auth.POST("/user/chats", handlers.GetUserChat)
		auth.POST("/user/chats/:folder", handlers.GetUserChatFolsers)
		auth.POST("/user/chats/info/:chet_id", handlers.GetUserChatInfo)
		auth.POST("/user/chats/info_details/:chat_id", handlers.GetChatDetails)
		auth.POST("/user/chats/messages/:chet_id", handlers.GetUserChatMessenges)
		auth.POST("/user/chats/boolen/:folder", handlers.GetUserChatFoldersBoolen)
		// История изменений имени пользователя
		auth.POST("/user/name-history", handlers.GetUserNameHistory)

		// Получение репостов
		auth.POST("/reposts", handlers.GetUserReposts)
		auth.POST("/reposts/:login", handlers.GetUsersReposts)

		// Чёрный список
		auth.POST("/blacklist/add", handlers.AddUserToBlacklist)
		auth.POST("/blacklist/remove", handlers.RemoveUserFromBlacklist)

		// Добавление поста
		auth.POST("/add-post-user", handlers.AddPostUser)

		// Добавление репоста
		auth.POST("/add-repost", handlers.AddRepost)

		// Отправка сообщений пользователю
		auth.POST("/send-message-one", handlers.SendMessageToUser)
		auth.POST("/send-message", handlers.SendMessageToChat)

		// Создание и добавление чатов в папки
		auth.POST("/add-chat-folder", handlers.AddChatFolder)
		auth.POST("/add-chat-to-folder", handlers.AddChatToFolder)

		// Обновление информации пользователя
		auth.POST("/update-info-user", handlers.UpdateUserInfo)
		auth.POST("/update-user-email", handlers.UpdateUserEmail)
		auth.POST("/update-user-birthdate", handlers.UpdateUserBirthDate)
		auth.POST("/update-user-login", handlers.UpdateUserLogin)
		auth.POST("/update-user-password", handlers.UpdateUserPassword)

		// Обновление информации о групповых чатах
		auth.POST("/update-group-chat", handlers.UpdateGroupChat)
		auth.POST("/remove-user-from-chat", handlers.RemoveUserFromChat)
		auth.POST("/delete-chat", handlers.DeleteChatAndUsers)

		// Удаление различных объектов
		auth.POST("/delete-sms-from-chat", handlers.DeleteSmsFromChat)
		auth.POST("/delete-chat-and-users", handlers.DeleteChatAndUsers)
		auth.POST("/delete-friend-request", handlers.DeleteFriendRequest)
		auth.POST("/delete-post-user", handlers.DeletePostUser)

		// Получение списка друзей
		auth.POST("/friends/list", handlers.GetFriendsList)

		// Отправка и принятие заявок в друзья
		auth.POST("/friend-request/send", handlers.SendFriendRequest)
		auth.POST("/friend-request/accept", handlers.AcceptFriendRequest)
		auth.POST("/friend-request/reject", handlers.RejectFriendRequest)
		auth.POST("/friend-request/incoming", handlers.GetFriendRequests)

		// Управление подписчиками
		auth.POST("/user/subscribers", handlers.GetSubscribersList)
		auth.POST("/user/remove-subscriber", handlers.RemoveSubscriber)
		auth.POST("/user/unsubscribe", handlers.UnsubscribeFromUser)

		// Управление стикер-паками
		auth.POST("/user/sticker-pack/create", handlers.CreateStickerPack)
		auth.POST("/user/sticker-pack/add-sticker", handlers.AddStickerToPack)
		auth.POST("/user/sticker-pack/list", handlers.GetUserStickerPacks)
		auth.POST("/sticker/send", handlers.SendStickerToUser)

		// Управление группами
		auth.POST("/user/group/add", handlers.AddGroup)
		auth.POST("/group/:group_id/subscribers", handlers.GetGroupSubscribers)
		auth.POST("/group/:group_id/posts", handlers.GetGroupPosts)
		auth.POST("/group/:group_id/add-post-comment", handlers.AddCommentToGroupPost)
		auth.POST("/group/:group_id/blacklist/add", handlers.AddUserToGroupBlacklist)
		auth.POST("/group/:group_id/blacklist/remove", handlers.RemoveUserFromGroupBlacklist)
		auth.POST("/group/info", handlers.GetGroupInfo)
		// Управление ролями в группе
		auth.POST("/groups/:group_id/roles", handlers.CreateRoleWithFeatures)
		auth.POST("/roles/:role_id/features", handlers.AddFeatureToRole)
		auth.POST("/roles/:role_id/users", handlers.AddUserToRoleGroup)
		auth.POST("/groups/:group_id/roles-list", handlers.GetRolesInfoForGroup)

		// Теги для групп
		auth.POST("/tags/create", handlers.CreateTag)
		auth.POST("/tags/:group_id", handlers.GetTags)
		auth.POST("/group/update-tags", handlers.UpdateTagsInGroup)

		// Управление постами пользователя
		auth.POST("/user/post/delete", handlers.DeletePost)
		auth.POST("/user/post/update", handlers.UpdatePost)

		// Поиск
		auth.POST("/search/posts", handlers.SearchPosts)
		auth.POST("/search/users", handlers.SearchUsers)
		auth.POST("/search/group", handlers.SearchGroups)
		auth.POST("/search/all", handlers.SearchAll)

		// Комментарии к постам
		auth.POST("/comments", handlers.GetComments)
		auth.POST("/toggle_like_comment", handlers.ToggleLikeComment)

		// Добавление в избранное
		auth.POST("/add_to_favourites", handlers.AddToFavourites)
		auth.POST("/user/favourite-sms", handlers.GetUserFavouriteSMS)

		// Управление чатами и папками чатов
		auth.POST("/user/chats-in-folder", handlers.GetChatsByFolder)
		auth.POST("/user/folder/chat/remove", handlers.RemoveChatFromFolder)
		auth.POST("/user/folder/remove", handlers.RemoveChatFolder)

		// Новый путь для получения рекомендованных постов
		auth.POST("/recommended-posts", handlers.GetRecommendedPosts)
		auth.POST("/filtered-posts", handlers.GetFilteredPosts)
		auth.POST("/get-user-chat-folders", handlers.GetUserFolders)

		//обновить имя и фото чата
		auth.POST("/chat/update", handlers.UpdateChatInfo)
		// проверка есть ли в друзьях
		auth.POST("/friends/check/:id", handlers.AreFriends)
		//есть ли заявка в друзья
		auth.POST("/friends/checkR/:other_user_id", handlers.CheckFriendRequest)
		//подписки
		auth.POST("/subscription/check/:other_user_id", handlers.CheckSubscription)
		//группы
		auth.POST("/groupsUs/me", handlers.GetGroupsForAuthorizedUser)
		auth.POST("/groupUs/:id", handlers.GetGroupsByUserID)
	}
}
