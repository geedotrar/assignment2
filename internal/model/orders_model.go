package model

import (
	"time"
)

type Order struct {
	OrderId      uint64    `json:"order_id" gorm:"primaryKey"`
	CustomerName string    `json:"customer_name"`
	OrderedAt    time.Time `json:"ordered_at"`
	Item         []Item    `json:"item" gorm:"foreignKey:order_id; references:order_id"`
}

type Item struct {
	ItemId      uint   `json:"item_id"`
	ItemCode    int    `json:"item_code"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	OrderID     uint64 `json:"order_id"`
}
