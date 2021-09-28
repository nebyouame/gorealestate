package propertypage

import "trail1/entity"



type CategoryRepository interface {
	Categories() ([]entity.Category, []error)
	Category(id uint) (*entity.Category, []error)
	UpdateCategory(category *entity.Category) (*entity.Category, []error)
	DeleteCategory(id uint) (*entity.Category, []error)
	StoreCategory(category *entity.Category) (*entity.Category, []error)
	PropertiesInCategory(category *entity.Category) ([]entity.Property, []error)
}


type PropertyRepository interface {
	Properties() ([]entity.Property, []error)
	Property(id uint) (*entity.Property, []error)
	UpdateProperty(property *entity.Property) (*entity.Property, []error)
	DeleteProperty(id uint) (*entity.Property, []error)
	StoreProperty(property *entity.Property) (*entity.Property, []error)
	RateProperty(property *entity.Property) (*entity.Property, []error)
	SearchProperty(index string) ([]entity.Property, error)
	StorePropertyCateg(property *entity.Property) []error
}


