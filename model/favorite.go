package model

type FavoriteActionRequest struct {
	// Token      string `json:"token" form:"token" xml:"token" binding:"required"`
	VideoId    int64 `json:"video_id" form:"video_id" xml:"video_id" binding:"required"`
	ActionType int32 `json:"action_type" form:"action_type" xml:"action_type" binding:"required,oneof=1 2"`
	// ActionType: 1 - 点赞, 2 - 取消点赞
}

type FavoriteListRequest struct {
	// Token  string `json:"token" form:"token" xml:"token" binding:"required"`
	UserId int64 `json:"user_id" form:"user_id" xml:"user_id" binding:"required"`
}

type FavoriteListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
