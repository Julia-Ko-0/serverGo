package routes

import (
	"serverGo/handlers"
	"serverGo/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Auth
	r.POST("/login", handlers.LoginUser)
	r.POST("/logout", handlers.LogoutUser)
	// r.POST("/register", handlers.Register)
	// r.POST("/update-info-user", handlers.UpdateUserInfo)
	// Protected routes
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{

		// auth.POST("/update-user-email", handlers.UpdateUserEmail)
		// auth.POST("/update-info-user", handlers.UpdateUserInfo)
		// auth.POST("/update-user-birthdate", handlers.UpdateUserBirthDate)
		// auth.POST("/update-user-password", handlers.UpdateUserPassword)

		//пол свои репосты
		auth.POST("/reposts", handlers.GetUserReposts)
		//пол репосты польз
		auth.POST("/reposts/:login", handlers.GetUsersReposts)
		// Получение информации о пользователе по логину
		auth.POST("/user_info", handlers.GetUserInfo)
		auth.POST("/user_info/:login", handlers.GetUserInfo)
		// Получение постов пользователя по логину
		auth.POST("/user_info_post", handlers.GetUserPosts)
		auth.POST("/user_info_post/:login", handlers.GetUsersPosts)
		// // Получение папок чатов
		auth.POST("/user/chat-folders", handlers.GetUserFolders)
		// Получение чатов
		auth.POST("/user/chats", handlers.GetUserChat)
		// Получение чатов папки
		auth.POST("/user/chats/:folder", handlers.GetUserChatFolsers)
		//Получение информации о чате(пользователи)
		auth.POST("/user/chats/info/:chet_id", handlers.GetUserChatInfo)
		// Получение сообщений чата
		auth.POST("/user/chats/messenge/:chet_id", handlers.GetUserChatMessenges)

		// история ников
		auth.POST("/user/name-history", handlers.GetUserNameHistory)

		// POST/////////////////////////////////////////////////////////////////////

		//доб в чс пользователя
		auth.POST("/blacklist/add", handlers.AddUserToBlacklist)
		//Добававить пост к себе на страницу
		auth.POST("/add-post-user", handlers.AddPostUser)
		//доб репост
		auth.POST("/add-repost", handlers.AddRepost)

		//Отправка сообющения 1in1
		auth.POST("/send-message-one", handlers.SendMessageToUser)
		// Отправка сообющения chat/1на1
		auth.POST("/send-message", handlers.SendMessageToChat)
		// добавление новой папки для чатов
		auth.POST("/add-chat-folder", handlers.AddChatFolder)

		// добавление чата в папку пользователя
		auth.POST("/add-chat-to-folder", handlers.AddChatToFolder)

		// PUT/////////////////////////////////////////////////////////////////
		// изм имя/фамилия/отчество пользователя
		auth.POST("/update-info-user", handlers.UpdateUserInfo)

		// изм почты
		auth.POST("/update-user-email", handlers.UpdateUserEmail)

		// изм даты рождения
		auth.POST("/update-user-birthdate", handlers.UpdateUserBirthDate)

		// логина(ника)
		auth.POST("/update-user-login", handlers.UpdateUserLogin)

		// пороля
		auth.POST("/update-user-password", handlers.UpdateUserPassword)

		// обновление информации о групповом чате (имя, фото)
		auth.POST("/update-group-chat", handlers.UpdateGroupChat)

		// удаление пользователя из чата (или выход из чата)
		auth.POST("/remove-user-from-chat", handlers.RemoveUserFromChat)

		// удаление чата (если ты используешь DeleteChatRequest)
		auth.POST("/delete-chat", handlers.DeleteChatAndUsers)

		//DELET //////////////////////////////////////////////////////////////////////////////
		//уд сообщения
		auth.POST("/delete-sms-from-chat", handlers.DeleteSmsFromChat)
		//уд чата(польз из него тоже удаляться)
		auth.POST("/delete-chat-and-users", handlers.DeleteChatAndUsers)
		//уд заявки в друзья
		auth.POST("/delete-friend-request", handlers.DeleteFriendRequest)
		//уд поста пользлвателя
		auth.POST("/delete-post-user", handlers.DeletePostUser)
		//уд из чс пользователя
		auth.POST("/blacklist/remove", handlers.RemoveUserFromBlacklist)
		//уд из друзей
		auth.POST("/friends/remove", handlers.RemoveFriend)
		//уд репост
		auth.POST("/repost/remove", handlers.RemoveRepost)

	}
}
