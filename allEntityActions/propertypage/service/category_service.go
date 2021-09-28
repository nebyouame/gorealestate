package service

import (
	"trail1/allEntityActions/propertypage"
	"trail1/entity"
)

//type CategoryService struct {
//	categoryRepo productpage.CategoryRepository
//}
//
//// NewCategoryService will create new CategoryService object
//func NewCategoryService(CatRepo productpage.CategoryRepository) productpage.CategoryService {
//	return &CategoryService{categoryRepo: CatRepo}
//}

type CategoryService struct {
	categoryRepo propertypage.CategoryRepository
}

func NewCategoryService(catRepo propertypage.CategoryRepository) propertypage.CategoryService {
	return &CategoryService{categoryRepo: catRepo}
}

func (cs *CategoryService) Categories() ([]entity.Category, []error) {
	categories, errs := cs.categoryRepo.Categories()
	if len(errs) > 0 {
		return nil, errs
	}

	return categories, nil
}


func (cs *CategoryService) Category(id uint) (*entity.Category, []error) {
	c, err := cs.categoryRepo.Category(id)
	if err != nil {
		return c, err
	}

	return c, nil
}

func (cs *CategoryService) StoreCategory(category *entity.Category) (*entity.Category, []error) {
	cat, errs := cs.categoryRepo.StoreCategory(category)
	if len(errs) > 0 {
		return nil, errs
	}

	return cat, nil
}

func (cs *CategoryService) UpdateCategory(category *entity.Category) (*entity.Category, []error) {
	cat, errs := cs.categoryRepo.UpdateCategory(category)
	if len(errs) > 0 {
		return nil, errs
	}

	return cat, nil
}

func (cs *CategoryService) DeleteCategory(id uint) (*entity.Category, []error) {
	cat, errs := cs.categoryRepo.DeleteCategory(id)
	if len(errs) > 0 {
		return nil, errs
	}

	return cat, nil
}

func (cs *CategoryService) PropertiesInCategory(category *entity.Category) ([]entity.Property, []error) {
	cts, errs := cs.categoryRepo.PropertiesInCategory(category)
	if len(errs) > 0 {
		return nil, errs
	}

	return cts, nil

}
