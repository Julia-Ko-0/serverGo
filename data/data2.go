package data

type FriendActionRequest struct {
	TargetID int `json:"target_id" binding:"required"` // ID пользователя, к которому относится действие (запрос/подтверждение/отклонение)
}

type Subscriber struct {
	ID             int    `json:"id_user"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}

type RemoveSubscriberRequest struct {
	SubscriberID int `json:"subscriber_id"`
}

type UnsubscribeRequest struct {
	TargetID int `json:"target_id"`
}

type UpdatePostRequest struct {
	PostID     int     `json:"post_id"`
	Header     *string `json:"header"`
	Text       *string `json:"text"`
	FileBase64 *string `json:"file_base64"`
}
type AddUserPostCommentRequest struct {
	PostID          int    `json:"post_id"`
	UserID          int    `json:"user_id"`
	CommentText     string `json:"comment_text"`
	ParentCommentID *int   `json:"parent_comment_id"` // может быть nil
}
type ModifyBlacklistRequest struct {
	UserID int `json:"user_id"` // ID пользователя, которого нужно заблокировать/разблокировать
}

// Запрос на создание стикер-пака
type CreateStickerPackRequest struct {
	Name string `json:"name"`
}

// Запрос на добавление стикера в пак
type AddStickerToPackRequest struct {
	StickerPackID int    `json:"sticker_pack_id"`
	FileBase64    string `json:"file_base64"`
}

// Структура для запроса на создание группового чата
type CreateGroupChatRequest struct {
	CreatorID   int     `json:"creator_id"`
	ChatName    string  `json:"chat_name"`
	UserIDs     []int   `json:"user_ids,omitempty"`
	PhotoBase64 *string `json:"photo_base64,omitempty"`
}

// Структура для запроса на отправку сообщения пользователю
type SendMessageToUserRequest struct {
	SenderID    int     `json:"sender_id"`
	ReceiverID  int     `json:"receiver_id"`
	MessageText string  `json:"message_text"`
	MessageFile *string `json:"message_file,omitempty"`
}

// Структура для запроса на отправку стикера пользователю
type SendStickerToUserRequest struct {
	SenderID   int `json:"sender_id"`
	ReceiverID int `json:"receiver_id"`
	StickerID  int `json:"sticker_id"`
}

// Структура для запроса на получение чатов пользователя
type GetUserChatsRequest struct {
	UserID int `json:"user_id"`
}

// Запрос на отправку стикера
type SendStickerRequest struct {
	ChatID    int `json:"chat_id"`    // ID чата
	StickerID int `json:"sticker_id"` // ID стикера
}

// Запрос на удаление сообщения
type DeleteSMSRequest struct {
	SmsID int `json:"sms_id"` // ID сообщения
}

// Запрос на получение всех папок пользователя
type GetUserChatFoldersRequest struct {
	UserID int `json:"user_id"` // ID пользователя
}

// Запрос на получение чатов по папке
type GetChatsByFolderRequest struct {
	UserID   int `json:"user_id"`   // ID пользователя
	FolderID int `json:"folder_id"` // ID папки
}
type ChatFolderActionRequest struct {
	ChatFolderID int `json:"chat_folder_id"`
	ChatID       int `json:"chat_id"`
}
type RemoveChatFolderRequest struct {
	ChatFolderID int `json:"chat_folder_id"`
}

// Тип данных для запроса на добавление группы
type AddGroupRequest struct {
	Name        string  `json:"name"`         // Название группы
	TypeGroupID int     `json:"type_gr_id"`   // ID типа группы
	Description *string `json:"description"`  // Описание группы (опционально)
	PhotoBase64 *string `json:"photo_base64"` // Фото группы (опционально)
}

// Тип данных для запроса информации о группе
type GetGroupInfoRequest struct {
	GroupID int `json:"group_id"` // ID группы
}
type Sticker struct {
	IDSticker  int    `json:"id_sticker"`
	ImgSticker string `json:"img_sticker"` // base64-кодированное изображение
}

type StickerPack struct {
	IDStickerPack   int       `json:"id_sticker_pack"`
	NameStickerPack string    `json:"name_sticker_pack"`
	Stickers        []Sticker `json:"stickers"`
}

type StickerPackResponse []StickerPack

type Chat struct {
	IDChat        int         `json:"id_chat"`
	NameChat      string      `json:"name_chat"`
	CountChatPepl int         `json:"countChatPepl"`
	Pfoto         string      `json:"pfoto"` // base64
	DateTimeChat  string      `json:"dateTime_chat"`
	LastMessage   LastMessage `json:"last_message"`
}

type ChatResponse []Chat

type ChatFolderResponse []ChatFolder

// Последнее сообщение
type FolderLastMessage struct {
	TextSMS        string `json:"text_sms"`
	DateTimeSMS    string `json:"dateTime_sms"`
	UserID         int    `json:"user_id"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}

// Чат в папке
type FolderChat struct {
	IDChat         int               `json:"id_chat"`
	NameChat       string            `json:"name_chat"`
	CountChatPepl  int               `json:"countChatPepl"`
	UserIDAdmin    int               `json:"user_id_admin"`
	DateTimeChat   string            `json:"dateTime_chat"`
	ProfilePicture string            `json:"profile_picture"`
	LastMessage    FolderLastMessage `json:"last_message"`
	ChatFolderID   int               `json:"chat_folder_id"`
}

type FolderChatsResponse []FolderChat

// Владелец группы
type GroupOwner struct {
	UserID         int    `json:"user_id"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"`
}

// Теги группы
type GroupTag struct {
	Name        string `json:"name_tag"`
	ID          int    `json:"id_tag"`
	Description string `json:"description_tag"`
}

// Информация о группе
type GroupInfoResponse struct {
	ID           int        `json:"id_group"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Access       string     `json:"access"`
	MembersCount int        `json:"members_count"`
	PostsCount   int        `json:"posts_count"`
	Owner        GroupOwner `json:"owner"`
	Tags         []GroupTag `json:"tags"`
}

// GroupSubscriber — структура одного подписчика в группе
type GroupSubscriber struct {
	ID             int    `json:"id_user"`         // ID пользователя
	Username       string `json:"username"`        // Логин
	Firstname      string `json:"firstname"`       // Имя
	Lastname       string `json:"lastname"`        // Фамилия
	Patronymic     string `json:"patronymic"`      // Отчество
	ProfilePicture string `json:"profile_picture"` // Фото профиля (base64)
}

// GroupInfo — информация о группе
type GroupInfo struct {
	ID          int        `json:"id_group"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Access      string     `json:"access"`
	Owner       GroupOwner `json:"owner"`
}

// PostAuthor — автор поста
type PostAuthor struct {
	UserID         int    `json:"user_id"`
	Username       string `json:"username"`
	ProfilePicture string `json:"profile_picture"` // base64
}

// GroupPost — структура поста в группе
type GroupPost struct {
	ID                 int        `json:"id_post_gr"`
	Header             string     `json:"header"`
	Text               string     `json:"text"`
	CommentsPermission bool       `json:"comments_permission"`
	DateTime           string     `json:"dateTime_post_gr"`
	Views              int        `json:"views_post"`
	Repost             bool       `json:"repost"`
	File               string     `json:"fale_post_gr"` // base64
	Author             PostAuthor `json:"author"`
}

// GroupPostsResponse — ответ функции get_group_posts
type GroupPostsResponse struct {
	GroupInfo GroupInfo   `json:"group_info"`
	Posts     []GroupPost `json:"posts"`
}

// FeatureRequest — структура запроса для добавления новой возможности
type FeatureRequest struct {
	AdminID            int    `json:"admin_id"`            // ID администратора
	NameFeatureRole    string `json:"name_featuresrole"`   // Название новой возможности
	DescriptionFeature string `json:"description_feature"` // Описание новой возможности
}

// FeatureInfo — структура для представления информации о возможности
type FeatureInfo struct {
	IDFeature          int    `json:"id_feature"`          // ID возможности
	NameFeatureRole    string `json:"name_featuresrole"`   // Название возможности
	DescriptionFeature string `json:"description_feature"` // Описание возможности
}

// CreateRoleRequest — структура для отправки данных на создание роли
type CreateRoleRequest struct {
	RoleName   string `json:"role_name"`    // Название роли
	FeatureIDs []int  `json:"features_ids"` // Список ID возможностей (features)
	CallerID   int    `json:"caller_id"`    // ID пользователя, вызывающего процедуру
}

// AddFeatureRequest — структура для отправки данных на добавление возможности в роль
type AddFeatureRequest struct {
	FeatureID int `json:"feature_id"` // ID возможности, которую нужно добавить
	CallerID  int `json:"caller_id"`  // ID пользователя, который вызывает процедуру
}

// AddUserToRoleRequest — структура для отправки данных на добавление пользователя в роль группы
type AddUserToRoleRequest struct {
	UserID   int `json:"user_id"`   // ID пользователя, которого нужно добавить
	CallerID int `json:"caller_id"` // ID пользователя, который вызывает процедуру
}

// RoleFeature — структура для описания возможности роли
type RoleFeature struct {
	FeatureID          int    `json:"feature_id"`          // ID возможности
	FeatureName        string `json:"feature_name"`        // Название возможности
	FeatureDescription string `json:"feature_description"` // Описание возможности
}

// UserRole — структура для описания пользователя, связанного с ролью
type UserRole struct {
	UserID         int    `json:"user_id"`         // ID пользователя
	Username       string `json:"username"`        // Имя пользователя
	ProfilePicture string `json:"profile_picture"` // Фото профиля в base64
}

// RoleInfo — структура для описания информации о роли
type RoleInfo struct {
	RoleID   int           `json:"role_id"`   // ID роли в группе
	RoleName string        `json:"role_name"` // Название роли
	Features []RoleFeature `json:"features"`  // Список возможностей роли
	Users    []UserRole    `json:"users"`     // Список пользователей, назначенных на роль
}
type GroupPostResponse struct {
	GroupInfo struct {
		ID          int    `json:"id_group"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Access      string `json:"access"`
		Owner       struct {
			UserID         int    `json:"user_id"`
			Username       string `json:"username"`
			ProfilePicture string `json:"profile_picture"`
		} `json:"owner"`
	} `json:"group_info"`

	Posts []struct {
		PostID          int    `json:"id_post_gr"`
		Header          string `json:"header"`
		Text            string `json:"text"`
		CommentsEnabled bool   `json:"comments_permission"`
		DateTime        string `json:"dateTime_post_gr"`
		Views           int    `json:"views_post"`
		Repost          int    `json:"repost"`
		ImageBase64     string `json:"fale_post_gr"`
		Author          struct {
			UserID         int    `json:"user_id"`
			Username       string `json:"username"`
			ProfilePicture string `json:"profile_picture"`
		} `json:"author"`
	} `json:"posts"`
}
type AddCommentToGroupPostRequest struct {
	PostID          int    `json:"post_id"`
	CommentText     string `json:"comment_text"`
	ParentCommentID *int   `json:"parent_comment_id,omitempty"` // nullable
}
type AddUserToGroupBlacklistRequest struct {
	GroupID int `json:"group_id"`
	UserID  int `json:"user_id"`
}

type RemoveUserFromGroupBlacklistRequest struct {
	UserID int `json:"user_id"` // ID пользователя, которого нужно удалить из черного списка
}
type SearchGroupsRequest struct {
	GroupID     int      `json:"group_id"`
	GroupName   string   `json:"group_name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}
type TagsGroupRequest struct {
	TagName        string `json:"tag_name"`
	TagDescription string `json:"tag_description"`
}
type UpdateTagGroup struct {
	GroupID int   `json:"group_id"`
	Tags    []int `json:"tags"` // массив ID тегов
}
type AddToFavouritesRequest struct {
	ItemID  int    `json:"item_id" binding:"required"`
	FavType string `json:"fav_type" binding:"required"`
}

type FavouritePost struct {
	PostID             int        `json:"post_id"`
	RepostDate         string     `json:"repost_date"`
	PostType           string     `json:"post_type"`
	Header             string     `json:"header"`
	Content            string     `json:"content"`
	FalePost           bool       `json:"fale_post"`
	DatetimePost       string     `json:"datetime_post"`
	ViewsPost          int        `json:"views_post"`
	Repost             bool       `json:"repost"`
	CommentsPermission bool       `json:"comments_permission"`
	CommentsCount      int        `json:"comments_count"`
	GroupInfo          GroupInfo  `json:"group_info"`
	Author             PostAuthor `json:"author"`
	LikesCount         int        `json:"likes_count"`
}

type FavouritePostsResponse struct {
	FavouritePosts []FavouritePost `json:"favourite_posts"`
}
type FavouriteSMS struct {
	SmsID          int    `json:"sms_id"`
	SmsText        string `json:"sms_text"`
	DateTime       string `json:"date_time"` // string, если приходит в виде ISO-формата
	FileSMS        []byte `json:"file_sms"`  // может быть nil
	SenderID       int    `json:"sender_id"`
	Login          string `json:"login_us"`
	Lastname       string `json:"lastname"`
	Firstname      string `json:"firstname"`
	Patronymic     string `json:"patronymic"`
	ProfilePicture []byte `json:"profile_picture"`
}
type GroupSearchResult struct {
	GroupID     int      `json:"group_id"`
	GroupName   string   `json:"group_name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"` // Массив тегов
}
type UserSearchResult struct {
	UserID         int    `json:"user_id"`
	Username       string `json:"username"`
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	Patronymic     string `json:"patronymic"`
	Description    string `json:"description"`
	ProfilePicture string `json:"profile_picture"` // base64 строка изображения
}

type Author struct {
	UserID         int    `json:"user_id"`
	Username       string `json:"username"`
	FirstName      string `json:"firstname"`
	LastName       string `json:"lastname"`
	ProfilePicture string `json:"profile_picture"` // base64 строка изображения
}

type Group struct {
	GroupID     int    `json:"group_id"`
	GroupName   string `json:"group_name"`
	Description string `json:"description"`
}

type Post struct {
	PostType string `json:"post_type"` // "user" или "group"
	PostID   int    `json:"post_id"`
	Header   string `json:"header"`
	Text     string `json:"text"`
	DatePost string `json:"date_post"` // дата поста
	FilePost string `json:"file_post"` // base64 строка файла, если есть
	Author   Author `json:"author"`
	Group    *Group `json:"group,omitempty"` // Если пост от группы
}

type SearchPostsResponse struct {
	FoundPost []Post `json:"found_post"`
}

type SearchAllResponse struct {
	FoundItems []Post `json:"result_item"`
}
