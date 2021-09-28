package entity

import (
	"time"
)

// type Order struct {
// 	ID        uint
// 	UserID    uint      `json:"userId" gorm:"not null"`
// 	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null;"`
// 	ItemsID   string    `json:"ProductId" gorm:"not null"`
// 	Total     float64   `json:"total" gorm:"type:float;not null;"`
// }

type Order struct {
	ID        uint
	UserID    uint      `json:"userId" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null;"`
	ItemsID   string    `json:"ProductId" gorm:"not null"`
	Total     float64   `json:"total" gorm:"type:float;not null;"`
	PropertyID uint		`json:"propertyId" gorm:"not null"`
	Name     string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"type:varchar(255);not null"`
	Phone    string `gorm:"type:varchar(100);not null"`
}



//type Order struct {
//	ID uint
//	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null;"`
//	ItemsID string `json:"PropertyId" gorm:"not null"`
//	Total float64 `json:"total" gorm:"type:float;not null;"`
//}