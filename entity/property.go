package entity
// type Product struct {
// 	ID          uint
// 	Name        string `gorm:"type:varchar(255);not null"`
// 	CategoryID  uint   `gorm:"many2many:product_categories"`
// 	Quantity    int
// 	Price       float64
// 	Description string
// 	Image       string `gorm:"type:varchar(255)"`
// 	Rating      float64
// 	RatersCount float64
// }

type Property struct {
	ID	uint
	Name string `gorm:"type:varchar(255);not null"`
	CategoryID uint `gorm:"many2many:property_categories"`
	Quantity int
	Price float64
	Description string
	Image string `gorm:"type:varchar(255)"`
	Image2 string `gorm:"type:varchar(255)"`
	Image3 string `gorm:"type:varchar(255)"`
	Image4 string `gorm:"type:varchar(255)"`
	Rating float64
	RatersCount float64
	UserId uint `json:"userId" gorm:"not null"`
}

//type Property struct {
//	ID	uint
//	Name string `gorm:"type:varchar(255);not null"`
//	CategoryID uint `gorm:"many2many:property_categories"`
//	Quantity int
//	Price float64
//	Description string
//	Image string `gorm:"type:varchar(255)"`
//	Rating float64
//	RatersCount float64
// 	UserID	uint	json:"userId" gorm:"not null"`
//}


