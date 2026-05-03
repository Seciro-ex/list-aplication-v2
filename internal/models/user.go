package models

import "time"

type User struct {
	ID         int       `json:"id"`
	Username   string    `json:"username" validate:"required,min=2,max=100"`
	PictureUrl string    `json:"picture_url"`
	Email      string    `json:"email" validate:"email,required,min=2,max=100"`
	Password   string    `json:"password" validate:"required,min=6,max=50"`
	CreatedAt  time.Time `json:"created_at"`
}
