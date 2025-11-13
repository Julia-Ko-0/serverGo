package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Глобальная переменная для базы данных
var DB *sqlx.DB

// Инициализация БД
func InitDB() {
	connStr := "host=localhost port=5433 user=postgres password=1234567890 dbname=mes sslmode=disable"
	var err error
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	fmt.Println(" Подключение к БД успешно!")
}

// Функция получения информации о пользователе
func GetUserInfo(userID int) (string, error) {
	var userInfo string

	// Проверяем, что соединение с БД установлено
	if DB == nil {
		return "", fmt.Errorf(" БД не инициализирована, вызови InitDB()")
	}

	// Запрос к БД
	err := DB.Get(&userInfo, "SELECT * from public.get_user_info($1)", userID)
	if err != nil {
		log.Println(" Ошибка запроса к БД:", err)
		return "", err
	}

	return userInfo, nil
}
