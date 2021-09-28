package repository


import (
"github.com/jinzhu/gorm"
"trail1/allEntityActions/propertypage"
"trail1/entity"
)

type CategoryGormRepo struct {
	conn *gorm.DB
}

func NewCategoryGormRepo(db *gorm.DB) propertypage.CategoryRepository  {
	return &CategoryGormRepo{conn: db}
}


func (cRepo CategoryGormRepo) Categories() ([]entity.Category, []error) {
	ctgs := []entity.Category{}
	errs := cRepo.conn.Find(&ctgs).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return ctgs, errs
}

//func (cRepo *CategoryGormRepo) Category(id uint) (*entity.Category, []error) {
//	ctg := entity.Category{}
//	errs := cRepo.conn.First(&ctg, id).GetErrors()
//	if len(errs) > 0 {
//		return nil, errs
//	}
//	return &ctg, errs
//}

func (cRepo *CategoryGormRepo) Category(id uint) (*entity.Category, []error) {
	ctg := entity.Category{}
	errs := cRepo.conn.First(&ctg, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &ctg, errs
}

func (cRepo *CategoryGormRepo) UpdateCategory(category *entity.Category) (*entity.Category, []error) {
	cat := category
	errs := cRepo.conn.Save(cat).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return cat, errs
}

// DeleteCategory deletes a given category from the database
func (cRepo *CategoryGormRepo) DeleteCategory(id uint) (*entity.Category, []error) {
	cat, errs := cRepo.Category(id)
	if len(errs) > 0 {
		return nil, errs
	}
	errs = cRepo.conn.Delete(cat, cat.ID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return cat, errs
}

// StoreCategory stores a given category in the database
func (cRepo *CategoryGormRepo) StoreCategory(category *entity.Category) (*entity.Category, []error) {
	cat := category
	errs := cRepo.conn.Create(cat).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return cat, errs
}

//func (cRepo *CategoryGormRepo) ItemsInCategory(category *entity.Category) ([]entity.Product, []error) {
//	items := []entity.Product{}
//	cat, errs := cRepo.Category(category.ID)
//
//	if len(errs) > 0 {
//		return nil, errs
//	}
//
//	errs = cRepo.conn.Model(cat).Related(&items, "Items").GetErrors()
//	if len(errs) > 0 {
//		return nil, errs
//	}
//	return items, errs
//}

func (cRepo *CategoryGormRepo) PropertiesInCategory(category *entity.Category) ([]entity.Property, []error)  {
	properties := []entity.Property{}
	cat, errs := cRepo.Category(category.ID)
	if len(errs) > 0 {
		return nil, errs
	}

	errs = cRepo.conn.Model(cat).Related(&properties, "Properties").GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return properties, errs
}





