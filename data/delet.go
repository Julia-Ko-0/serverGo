package data

type DeleteSmsRequest struct {
	SmsID int `json:"sms_id"`
}
type DeleteChatRequest struct {
	ChatID int `json:"chat_id"`
}
type DeleteFriendRequest struct {
	SenderID   int `json:"sender_id"`
	ReceiverID int `json:"receiver_id"`
}
type DeletePostRequest struct {
	PostID int `json:"post_id"`
}
type RemoveUserFromChatRequest struct {
	ExecutorID int `json:"executor_id"`
	ChatID     int `json:"chat_id"`
}
type FriendRequest struct {
	FriendID int `json:"friend_id" binding:"required"`
}

type RepostRemoveRequest struct {
	RepostID int `json:"repost_id" binding:"required"`
}
type BlockRequest struct {
	BlockedID int `json:"blocked_id" binding:"required"`
}
