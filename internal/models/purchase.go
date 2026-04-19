package models

import "github.com/shopspring/decimal"

type Purchase struct {
	ID       int             `json:"id"`
	ImageURL string          `json:"image_url"`
	Cost     decimal.Decimal `json:"cost" validate:"required,gt=0"`
	ItemURL  string          `json:"item_url"`
	ItemName string          `json:"item_name" validate:"required,min=2,max=100"`
	UserID   int             `json:"user_id"`
}
