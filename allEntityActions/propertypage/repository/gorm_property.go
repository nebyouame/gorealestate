package repository

import (
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"math"
	"trail1/allEntityActions/propertypage"
	"trail1/entity"
)




type PropertyGormRepo struct {
	conn *gorm.DB
}

func NewPropertyGormRepo(db *gorm.DB) propertypage.PropertyRepository {
	return &PropertyGormRepo{conn:db}
}




func (propertyRepo *PropertyGormRepo) Properties() ([]entity.Property, []error) {
	properties := []entity.Property{}
	errs := propertyRepo.conn.Find(&properties).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return properties, errs
}


func (propertyRepo *PropertyGormRepo) Property(id uint) (*entity.Property, []error) {
	property := entity.Property{}
	errs := propertyRepo.conn.First(&property, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &property, errs
}



func (propertyRepo *PropertyGormRepo) UpdateProperty(property *entity.Property) (*entity.Property, []error)  {
	pro := property
	errs := propertyRepo.conn.Save(pro).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return pro, errs
}




func (propertyRepo *PropertyGormRepo) DeleteProperty(id uint) (*entity.Property, []error) {
	pro, errs := propertyRepo.Property(id)

	if len(errs) > 0 {
		return nil, errs
	}
	errs = propertyRepo.conn.Delete(pro, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return pro, errs
}



func (propertyRepo *PropertyGormRepo) StoreProperty(property *entity.Property) (*entity.Property, []error)  {
	pro := property
	errs := propertyRepo.conn.Create(pro).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return pro, errs
}



func (propertyRepo *PropertyGormRepo) SearchProperty(index string) ([]entity.Property, error) {
	properties := []entity.Property{}


	err := propertyRepo.conn.Where("name ILIKE ?", "%"+index+"&").Find(&properties).GetErrors()
	if len(err) !=0 {
		errors.New("Search Property not working")
	}
	return properties, nil
}



func (propertyRepo *PropertyGormRepo) RateProperty(pro *entity.Property) (*entity.Property, []error) {
	u := entity.Property{}
	property := entity.Property{}
	row := propertyRepo.conn.Select("rating").First(&property).Where("id = ?", pro.ID).Scan(&u)
	log.Println("Old rate", u.Rating)
	if row.RecordNotFound() {
		panic(row.Error)
	}

	row = propertyRepo.conn.Select("raters_count").First(&property).Where("id = ?", pro.ID).Scan(&u)
	log.Println("Old count", u.RatersCount)
	if row.RecordNotFound() {
		panic(row.Error)
	}

	newratings := ((u.Rating * u.RatersCount) + pro.Rating) / (u.RatersCount + 1)
	log.Println(newratings)
	log.Println("Pro ", pro.Rating)

	row = propertyRepo.conn.Model(&pro).Updates(entity.Property{Rating: float64((math.Round((newratings * 2 )))) / 2, RatersCount:u.RatersCount + 1 })
	if row.RowsAffected < 1 {
		return &property, []error{errors.New("Error")}
	}
	return &property, nil
}



func (propertyRepo *PropertyGormRepo) StorePropertyCateg(property *entity.Property) []error {
	pro := property

	err := propertyRepo.conn.Exec("Insert into property_categories (property_id, category_id) values (?, ?)", pro.ID, pro.CategoryID).GetErrors()
	if err != nil {
		return err
	}
	return nil
}














