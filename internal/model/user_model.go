package model

type PostUserV1Request struct {
	Username string `json:"username" binding:"required,min=5,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type PostUserV1Response struct {
	ID int64 `json:"id"`
}
