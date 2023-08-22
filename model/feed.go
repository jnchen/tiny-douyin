package model

type FeedRequest struct {
	LatestTime int64  `json:"latest_time" form:"latest_time" xml:"latest_time"`
	Token      string `json:"token" form:"token" xml:"token"`
}

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list"`
	NextTime  int64   `json:"next_time,omitempty"`
}
