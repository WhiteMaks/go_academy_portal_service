package model

type PostUserV1Request struct {
	Username string `json:"username" binding:"required,min=5,max=64,alphanum"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type PostUserV1Response struct {
	ID int64 `json:"id"`
}

type PostUserTokenV1Request struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PostUserTokenV1Response struct {
	Token string `json:"token"`
}
