package models

type Product struct {
	Id                string           `json:"id" gorm:"primaryKey"`
	Name              string           `json:"name"`
	Stock             int64            `json:"stock"`
	Price             float32          `json:"price"`
	StockDecreaseLogs StockDecreaseLog `gorm:"foreignKey:ProductRefer"`
}
