package models

import "github.com/shopspring/decimal"

type Purchase struct {
	ID       int             `json:"id"`
	ImageURL string          `json:"image_url"`
	Cost     decimal.Decimal `json:"cost"`
	ItemURL  string          `json:"item_url"`
	ItemName string          `json:"item_name"`
	UserID   int             `json:"user_id"`
}
