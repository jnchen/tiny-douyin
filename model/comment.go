package model

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment"`
}

type CommentActionRequest struct {
	Token      string `json:"token"`
	VideoId    int64  `json:"video_id"`
	ActionType int32  `json:"action_type"`            // 1 - post, 2 - delete
	Content    string `json:"comment_text,omitempty"` // used when action_type = 1
	CommentId  int64  `json:"comment_id,omitempty"`   // used when action_type = 2
}

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list"`
}

type CommentListRequest struct {
	Token   string `json:"token"`
	VideoId int64  `json:"video_id"`
}
