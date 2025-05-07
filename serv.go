// package main

// import (
// 	"fmt"
// 	"log"
// 	"serverGo/db"
// 	"serverGo/handlers"

// 	"time"

// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	// Инициализация базы данных
// 	db.InitDB()

// 	// Создание роутера Gin
// 	r := gin.Default()

// 	// Настройка CORS
// 	r.Use(cors.New(cors.Config{
// 		AllowOrigins:     []string{"http://localhost:3000"}, // Разрешите фронту доступ
// 		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
// 		ExposeHeaders:    []string{"Content-Length"},
// 		AllowCredentials: true,
// 		MaxAge:           12 * time.Hour,
// 	}))

// 	// GET ////////////////////////////////////////////////////////////////////////////

// 	// Получение информации о пользователе по логину
// 	r.GET("/user_info/:login", handlers.GetUserInfo)
// 	// Получение постов пользователя по логину
// 	r.GET("/user_info_post/:login", handlers.GetUserPosts)
// 	// // Получение папок чатов
// 	r.GET("/user/:login/chat-folders", handlers.GetUserFolders)
// 	// Получение чатов
// 	r.GET("/user/:login/chats", handlers.GetUserChat)
// 	// Получение чатов папки
// 	r.GET("/user/:login/chats/:folder", handlers.GetUserChatFolsers)
// 	//Получение информации о чате(пользователи)
// 	r.GET("/user/:login/chats/info/:chet_id", handlers.GetUserChatInfo)
// 	// Получение сообщений чата
// 	r.GET("/user/:login/chats/messenge/:chet_id", handlers.GetUserChatMessenges)

// 	// история ников
// 	r.GET("/user/:login/name-history", handlers.GetUserNameHistory)

// 	// POST/////////////////////////////////////////////////////////////////////
// 	//Отправка сообющения 1in1
// 	r.POST("/send-message-one", handlers.SendMessageToUser)
// 	// Отправка сообющения chat/1на1
// 	r.POST("/send-message", handlers.SendMessageToChat)
// 	// добавление новой папки для чатов
// 	r.POST("/add-chat-folder", handlers.AddChatFolder)

// 	// добавление чата в папку пользователя
// 	r.POST("/add-chat-to-folder", handlers.AddChatToFolder)

// 	// PUT/////////////////////////////////////////////////////////////////
// 	// изм имя/фамилия/отчество пользователя
// 	r.POST("/update-info-user", handlers.UpdateUserInfo)

// 	// изм почты
// 	r.POST("/update-user-email", handlers.UpdateUserEmail)

// 	// изм даты рождения
// 	r.POST("/update-user-birthdate", handlers.UpdateUserBirthDate)

// 	// логина(ника)
// 	r.POST("/update-user-login", handlers.UpdateUserLogin)

// 	// пороля
// 	r.POST("/update-user-password", handlers.UpdateUserPassword)
// 	// изменение имени/фамилии/отчества пользователя
// 	r.POST("/update-info-user", handlers.UpdateUserInfo)

// 	// обновление информации о групповом чате (имя, фото)
// 	r.POST("/update-group-chat", handlers.UpdateGroupChat)

// 	// удаление пользователя из чата (или выход из чата)
// 	r.POST("/remove-user-from-chat", handlers.RemoveUserFromChat)

// 	// удаление чата (если ты используешь DeleteChatRequest)
// 	r.POST("/delete-chat", handlers.DeleteChatAndUsers)

// 	//DELET //////////////////////////////////////////////////////////////////////////////
// 	//уд сообщения
// 	r.POST("/delete-sms-from-chat", handlers.DeleteSmsFromChat)
// 	//уд чата(польз из него тоже удаляться)
// 	r.POST("/delete-chat-and-users", handlers.DeleteChatAndUsers)
// 	//уд заявки в друзья
// 	r.POST("/delete-friend-request", handlers.DeleteFriendRequest)
// 	//уд поста пользлвателя
// 	r.POST("/delete-post-user", handlers.DeletePostUser)

//		// Запуск сервера
//		port := ":8080"
//		fmt.Println("Сервер запущен на порту", port)
//		if err := r.Run(port); err != nil {
//			log.Fatalf("Ошибка запуска сервера: %v", err)
//		}
//	}
package main

import (
	"log"
	"os"
	"serverGo/db"
	"serverGo/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}
	var jwtKey = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtKey) == 0 {
		log.Println(" JWT_SECRET is EMPTY! Please check your .env file or environment variables.")
	}

	db.InitDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // адрес твоего фронта
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // обязательно, чтобы cookies передавались
		MaxAge:           12 * time.Hour,
	}))
	// r.Use(cors.Default())

	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
