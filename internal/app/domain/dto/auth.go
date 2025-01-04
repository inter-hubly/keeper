package dto

type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
}
