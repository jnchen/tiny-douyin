package model

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment"`
}

type CommentActionRequest struct {
	Token      string `json:"token" form:"token" xml:"token" binding:"required"`
	VideoId    int64  `json:"video_id" form:"video_id" xml:"video_id" binding:"required"`
	ActionType int32  `json:"action_type" form:"action_type" xml:"action_type" binding:"required,oneof=1 2"` // 1 - post, 2 - delete
	Content    string `json:"comment_text,omitempty" form:"comment_text" xml:"comment_text"`                 // used when action_type = 1
	CommentId  int64  `json:"comment_id,omitempty" form:"comment_id" xml:"comment_id"`                       // used when action_type = 2
}

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list"`
}
