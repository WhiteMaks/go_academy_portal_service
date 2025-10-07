package model

type PostUserV1Request struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type PostUserV1Response struct {
	ID int64 `json:"id"`
}
