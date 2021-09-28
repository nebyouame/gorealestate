package entity

// type Category struct {
// 	ID          uint
// 	Name        string `gorm:"type:varchar(255);not null"`
// 	Description string
// 	Image       string    `gorm:"type:varchar(255)"`
// 	Products    []Product `gorm:"many2many:product_categories"`
// }

type Category struct {
	ID uint
	Name string `gorm:"type:varchar(255);not null"`
	Description string
	Image string `gorm:"type:varchar(255)"`
	Properties []Property `gorm:"many2many:property_categories"`
}
