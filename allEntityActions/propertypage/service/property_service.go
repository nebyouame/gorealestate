package service

import (

	"trail1/allEntityActions/propertypage"
	"trail1/entity"
)

//type ItemService struct {
//itemRepo productpage.ItemRepository
//}
//
//// NewItemService returns new ItemService object
//func NewItemService(itemRepository productpage.ItemRepository) productpage.ItemService {
//return &ItemService{itemRepo: itemRepository}
//}


type PropertyService struct {
	propertyRepo propertypage.PropertyRepository
}

func NewPropertyService(propertyRepository propertypage.PropertyRepository) propertypage.PropertyService {
	return &PropertyService{propertyRepo:propertyRepository}
}


func (ps *PropertyService) Properties() ([]entity.Property, []error) {
	pros, errs := ps.propertyRepo.Properties()
	if len(errs) > 0 {
		return nil, errs
	}
	return pros, errs
}


func (ps *PropertyService) Property(id uint) (*entity.Property, []error) {
	pro, errs := ps.propertyRepo.Property(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return pro, nil
}



func (ps PropertyService) UpdateProperty(property *entity.Property) (*entity.Property, []error) {
	pro, errs := ps.propertyRepo.UpdateProperty(property)
	if len(errs) > 0 {
		return nil, errs
	}
	return pro, nil
}

func (ps *PropertyService) DeleteProperty(id uint) (*entity.Property, []error) {
	pro, errs := ps.propertyRepo.DeleteProperty(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return pro, nil
}



func (ps *PropertyService) StoreProperty(property *entity.Property) (*entity.Property, []error) {
	pro, errs := ps.propertyRepo.StoreProperty(property)
	if len(errs) > 0 {
		return nil, errs
	}
	return pro, nil
}


func (ps *PropertyService) SearchProperty(index string) ([]entity.Property, error) {
	properties, err := ps.propertyRepo.SearchProperty(index)
	if err != nil {
		return nil, err
	}
	return properties, nil
}



func (ps *PropertyService) RateProperty(pro *entity.Property) (*entity.Property, []error) {
	prowithrate, err := ps.propertyRepo.RateProperty(pro)
	if err != nil {
		return prowithrate, err
	}
	return prowithrate, nil
}




func (ps *PropertyService) StorePropertyCateg(property *entity.Property) []error {
	err := ps.propertyRepo.StorePropertyCateg(property)
	if err != nil {
		return err
	}
	return nil
}









