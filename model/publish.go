package model

import "mime/multipart"

type PublishActionRequest struct {
	// Token string                `json:"token" form:"token" xml:"token" binding:"required"`
	Title string                `json:"title" form:"title" xml:"title" binding:"required,min=1,max=128"`
	Data  *multipart.FileHeader `json:"data" form:"data" xml:"data" binding:"required"`
}

type PublishListRequest struct {
	// Token string `json:"token" form:"token" xml:"token" binding:"required"`
	UserId int64 `json:"user_id" form:"user_id" xml:"user_id" binding:"required"`
}

type PublishListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
