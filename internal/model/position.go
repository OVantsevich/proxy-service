// Package model position model
package model

// Position model
type Position struct {
	ID            string  `json:"id"`
	User          string  `json:"user"`
	Name          string  `json:"name" validate:"required,alpha,gte=2,lte=30"`
	Amount        float64 `json:"amount" validate:"required,gte=0"`
	SellingPrice  float64 `json:"selling_price"`
	PurchasePrice float64 `json:"purchase_price"`
	StopLoss      float64 `json:"stop_loss"`
	TakeProfit    float64 `json:"take_profit"`
	ShortPosition bool    `json:"short_position" validate:"required"`
	Closed        int64   `json:"closed"`
}
