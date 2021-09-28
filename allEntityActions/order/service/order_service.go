package service

import (
	"trail1/allEntityActions/order"
	"trail1/entity"
)


type OrderService struct {
	orderRepo order.OrderRepository
}

func NewOrderService(orderRepository order.OrderRepository) order.OrderService {
	return &OrderService{orderRepo: orderRepository}
}

func (os *OrderService) Orders() ([]entity.Order, []error) {
	ords, errs := os.orderRepo.Orders()
	if len(errs) > 0 {
		return nil, errs
	}
	return ords, errs
}

// Order retrieves an order by its id
func (os *OrderService) Order(id uint) (*entity.Order, []error) {
	ord, errs := os.orderRepo.Order(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return ord, errs
}

// CustomerOrders returns all orders of a given customer
func (os *OrderService) CustomerOrders(customer *entity.User) (entity.Order, []error) {
	ords, errs := os.orderRepo.CustomerOrders(customer)
	if len(errs) > 0 {
		return ords, errs
	}
	return ords, errs
}

// UpdateOrder updates a given order
func (os *OrderService) UpdateOrder(order *entity.Order) (*entity.Order, []error) {
	ord, errs := os.orderRepo.UpdateOrder(order)
	if len(errs) > 0 {
		return nil, errs
	}
	return ord, errs
}

// DeleteOrder deletes a given order
func (os *OrderService) DeleteOrder(id uint) (*entity.Order, []error) {
	ord, errs := os.orderRepo.DeleteOrder(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return ord, errs
}

// StoreOrder stores a given order
func (os *OrderService) StoreOrder(order *entity.Order) (*entity.Order, []error) {
	ord, errs := os.orderRepo.StoreOrder(order)
	if len(errs) > 0 {
		return nil, errs
	}
	return ord, errs
}
