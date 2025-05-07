package data

type UpdateUserInfoRequest struct {
	Lastname   string `json:"Lastname"`   // Может быть пустым
	Firstname  string `json:"Firstname"`  // Может быть пустым
	Patronymic string `json:"Patronymic"` // Может быть пустым
}
type UpdateUserEmailRequest struct {
	NewEmail string `json:"new_email"`
}
type UpdateUserBirthDateRequest struct {
	BirthDate string `json:"birth_date"` // формат: "2006-01-02"
}
type UpdateUserLoginRequest struct {
	NewLogin string `json:"new_login"`
}
type UpdateUserPasswordRequest struct {
	NewPassword string `json:"new_password"`
}
type UpdateGroupChatRequest struct {
	ChatID int `json:"chat_id"`

	ChatName *string `json:"chat_name,omitempty"` // nullable
	Pfoto    []byte  `json:"pfoto,omitempty"`     // optional
}
