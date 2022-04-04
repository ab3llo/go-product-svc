package models

type StockDecreaseLog struct {
	Id           string `json:"id" gorm:"primaryKey"`
	OrderId      string `json:"orderId"`
	ProductRefer string `json:"productId"`
}
