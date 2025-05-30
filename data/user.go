package data

import "time"

// Структура для хранения данных пользователя
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Структура для информации о пользователе
type UserInfoResponse struct {
	ID               int     `json:"id_user"`
	Username         string  `json:"username"`
	Email            string  `json:"email"`
	Lastname         *string `json:"lastname,omitempty"`
	Firstname        *string `json:"firstname,omitempty"`
	Patronymic       *string `json:"patronymic,omitempty"`
	DateBirth        *string `json:"date_birth,omitempty"`
	RegistrationDate string  `json:"registration_date"`
	CountUs          *int    `json:"count_us,omitempty"`
	FriendsCount     int     `json:"friends_count"`
	SubscribersCount int     `json:"subscribers_count"`
	BlacklistCount   int     `json:"blacklist_count"`
	Description      *string `json:"description,omitempty"`
	ProfilePicture   *string `json:"profile_picture,omitempty"`
}

// Структура для постов пользователя
type UserPostResponse struct {
	UserInfo struct {
		UserID           int    `json:"user_id"`
		Username         string `json:"login_us"`
		FirstName        string `json:"firstname"`
		LastName         string `json:"lastname"`
		Patronymic       string `json:"patronymic"`
		Registration     string `json:"dateRegistr"`
		BirthDate        string `json:"date_of_birth"`
		CountUs          int    `json:"count_us"`
		Description      string `json:"description"`
		ProfilePicture   string `json:"profile_picture"`
		FriendsCount     int    `json:"friends_count"`
		SubscribersCount int    `json:"subscribers_count"`
	} `json:"user_info"`

	Posts []struct {
		PostID      int    `json:"id_post_us"`
		Header      string `json:"header"`
		Text        string `json:"text_post"`
		DateTime    string `json:"dateTime_post_us"`
		Comments    bool   `json:"comments_enabled"`
		Views       int    `json:"views_post"`
		Repost      int    `json:"repost"`
		ImageBase64 string `json:"fale_post"`
	} `json:"posts"`
}
type ChatFolder struct {
	ID     int    `json:"id_chatFolders"`
	Name   string `json:"name_chatFolders"`
	UserID int    `json:"user_id"`
}

// UserChat содержит информацию о чате, в котором участвует пользователь
type UserChats struct {
	ID           int          `json:"id_chat"`
	Name         string       `json:"name_chat"`
	CountPeople  string       `json:"countChatPepl"` // <-- тип string
	Photo        string       `json:"pfoto"`
	DateTimeChat string       `json:"dateTime_chat"`
	LastMessage  *LastMessage `json:"last_message"` // указатель, т.к. может быть null
}

// LastMessage содержит информацию о последнем сообщении в чате
type LastMessage struct {
	Text           string `json:"text_sms"`
	DateTime       string `json:"dateTime_sms"`
	UserID         int    `json:"user_id"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"` // base64
}

type MessageSender struct {
	UserID         int    `json:"user_id"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"` // base64
}

type StickerInfo struct {
	StickerID     int    `json:"sticker_id"`
	StickerPackID int    `json:"sticker_pack_id"`
	Image         string `json:"img_sticker"` // base64
}

type UserNameHistory struct {
	UserID    int    `db:"user_is_name" json:"user_id"`
	Name      string `db:"name" json:"name"`
	CountName int    `db:"count_name" json:"count"`
}

/////////////////////////////////////
type RepostResponse struct {
	RepostID   int            `json:"repost_id"`
	RepostDate time.Time      `json:"repost_date"`
	PostInfo   map[string]any `json:"post_info"`
}

// Временная структура для разбора JSON + []byte → base64
type RawUserPostResponse struct {
	UserInfo struct {
		UserID           int    `json:"user_id"`
		Username         string `json:"login_us"`
		FirstName        string `json:"firstname"`
		LastName         string `json:"lastname"`
		Patronymic       string `json:"patronymic"`
		Registration     string `json:"dateRegistr"`
		BirthDate        string `json:"date_of_birth"`
		CountUs          int    `json:"count_us"`
		Description      string `json:"description"`
		ProfilePicture   []byte `json:"profile_picture"` // bytea
		FriendsCount     int    `json:"friends_count"`
		SubscribersCount int    `json:"subscribers_count"`
	} `json:"user_info"`

	Posts []struct {
		PostID    int    `json:"id_post_us"`
		Header    string `json:"header"`
		Text      string `json:"text_post"`
		DateTime  string `json:"dateTime_post_us"`
		Comments  bool   `json:"comments_enabled"`
		Views     int    `json:"views_post"`
		Repost    int    `json:"repost"`
		ImageData []byte `json:"fale_post"` // bytea
	} `json:"posts"`
}
type PostResponse struct {
	PostID      int    `json:"id_post_us"`
	Header      string `json:"header"`
	Text        string `json:"text_post"`
	DateTime    string `json:"dateTime_post_us"`
	Comments    bool   `json:"comments_enabled"`
	Views       int    `json:"views_post"`
	Repost      int    `json:"repost"`
	ImageBase64 string `json:"fale_post"`
}

//
// //////////////////////////////////////////////////////////////////////////
// /////////////////////////
// //////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////
type ChatMessage struct {
	ID       int           `json:"id_sms"`
	Text     string        `json:"text_sms"`
	DateTime string        `json:"dateTime_sms"`
	File     string        `json:"file_sms"` // base64
	Sticker  *StickerInfo  `json:"sticker"`  // может быть null
	User     MessageSender `json:"user"`     // вложенный объект
}
type RawChatMessage struct {
	ID       int             `json:"id_sms"`
	Text     string          `json:"text_sms"`
	DateTime string          `json:"dateTime_sms"`
	File     []byte          `json:"file_sms"`
	Sticker  *StickerInfoRaw `json:"sticker"`
	User     MessageSender   `json:"user"`
}

type StickerInfoRaw struct {
	StickerID     int    `json:"sticker_id"`
	StickerPackID int    `json:"sticker_pack_id"`
	Image         []byte `json:"img_sticker"`
}
type RawChatInfo struct {
	ChatID     int    `json:"chat_id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	LastUpdate string `json:"last_update"`
	Photo      []byte `json:"photo"`
}
type ChatInfo struct {
	ChatID     int    `json:"chat_id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	LastUpdate string `json:"last_update"`
	Photo      string `json:"photo"`
}
type RawChatUser struct {
	UserID         int    `json:"user_id"`
	Login          string `json:"login"`
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	ProfilePicture []byte `json:"profile_picture"`
}
type ChatUser struct {
	UserID         int    `json:"user_id"`
	Username       string `json:"username"`
	Role           string `json:"role"`
	ProfilePicture string `json:"profile_picture"` // base64 (может быть пустым)
}
