package order

import "trail1/entity"

type OrderService interface {
	Orders() ([]entity.Order, []error)
	Order(id uint) (*entity.Order, []error)
	CustomerOrders(customer *entity.User) (entity.Order, []error)
	UpdateOrder(order *entity.Order) (*entity.Order, []error)
	DeleteOrder(id uint) (*entity.Order, []error)
	StoreOrder(order *entity.Order) (*entity.Order, []error)
}


