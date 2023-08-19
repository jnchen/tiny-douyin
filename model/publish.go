package model

import "mime/multipart"

type PublishActionRequest struct {
	Token string                `json:"token" form:"token" xml:"token" binding:"required"`
	Title string                `json:"title" form:"title" xml:"title" binding:"required,min=1,max=32"`
	Data  *multipart.FileHeader `json:"data" form:"data" xml:"data" binding:"required"`
}

type PublishListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
