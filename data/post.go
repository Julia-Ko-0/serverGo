package data

// Запрос для отправки сообщения в чат
type SendMessageToChatRequest struct {
	ChatID      int    `json:"chat_id"`
	MessageText string `json:"message_text"`
	MessageFile string `json:"message_file"` // base64, может быть пустым
}
type SendMessageToChat struct {
	PChatId      int    `json:"p_chat_id"`
	PMessageText string `json:"p_message_text"`
	PMessageFile string `json:"p_message_file"` // base64, может быть пустым
}
type AddChatToFolderRequest struct {
	FolderID int `json:"folder_id"`
	ChatID   int `json:"chat_id"`
}
type AddChatFolderRequest struct {
	FolderName string `json:"folder_name"`
}

//////////////////////////////////////////////////////////
type AddPostUserRequest struct {
	TextPost string `json:"text_post"`
	Header   string `json:"header"`
	FalePost string `json:"fale_post"` // <-- тип string, т.к. это base64
}

type AddRepostRequest struct {
	PostID     int    `json:"post_id" binding:"required"`
	TypeRepost string `json:"type_repost" binding:"required"` // 'us' или 'gr'
}
