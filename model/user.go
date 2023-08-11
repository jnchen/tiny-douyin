package model

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

type UserRegisterRequest struct {
	UserName string `json:"username" form:"username" xml:"username" binding:"required,min=1,max=32"`
	Password string `json:"password" form:"password" xml:"password" binding:"required,min=1,max=32"`
}

type UserLoginRequest struct {
	UserName string `json:"username" form:"username" xml:"username" binding:"required"`
	Password string `json:"password" form:"password" xml:"password" binding:"required"`
}
