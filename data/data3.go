package data

type AddUserToGroupRequest struct {
	GroupID int `json:"group_id" binding:"required"`
}

type Tag struct {
	NameTag        string `json:"name_tag"`
	IDTag          int    `json:"id_tag"`
	DescriptionTag string `json:"description_tag"`
}

type GroupInfoResponse struct {
	IDGroup      int        `json:"group_id"`
	Name         string     `json:"group_name"`
	Description  string     `json:"description"`
	Access       bool       `json:"access"`
	MembersCount int        `json:"members_count"`
	PostsCount   int        `json:"posts_count"`
	Owner        GroupOwner `json:"owner"`
	Tags         []Tag      `json:"tags"`
	Photo        string     `json:"photo"` // Это поле будет хранить фото в base64
}
type RemoveFeatureRequest struct {
	RoleGroupID int `json:"roleGroupId" binding:"required"`
	FeatureID   int `json:"featureId" binding:"required"`
}

type DeleteRoleGroupRequest struct {
	RoleGroupID int `json:"role_group_id"`
}

type RemoveUserFromRoleRequest struct {
	RoleGroupID    int `json:"role_group_id" binding:"required"`
	UserIDToRemove int `json:"user_id_to_remove" binding:"required"`
}
type UpdateGroupRequest struct {
	GroupID        int    `json:"group_id"`
	NewName        string `json:"new_name"`
	NewDescription string `json:"new_description"`
	NewPhoto       string `json:"new_photo"` // base64 строка
}
