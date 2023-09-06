package model

import "mime/multipart"

type PublishActionRequest struct {
	Title string                `json:"title" form:"title" xml:"title" binding:"required,min=1,max=64"`
	Data  *multipart.FileHeader `json:"data" form:"data" xml:"data" binding:"required"`
}

type PublishListRequest struct {
	UserId int64 `json:"user_id" form:"user_id" xml:"user_id" binding:"required,gte=1"`
}

type PublishListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
