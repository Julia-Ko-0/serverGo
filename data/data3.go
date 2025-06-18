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
	IDGroup      int        `json:"id_group"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Access       bool       `json:"access"`
	MembersCount int        `json:"members_count"`
	PostsCount   int        `json:"posts_count"`
	Owner        GroupOwner `json:"owner"`
	Tags         []Tag      `json:"tags"`
}
type RemoveFeatureRequest struct {
	RoleGroupID int `json:"roleGroupId" binding:"required"`
	FeatureID   int `json:"featureId" binding:"required"`
}
