package service

import (
	"context"

	"assignment2/internal/model"
	"assignment2/internal/repository"
)

type OrderService interface {
	GetOrders(ctx context.Context) ([]model.Order, error)
	GetOrderById(ctx context.Context, id uint64) (model.Order, error)
	CreateOrder(ctx context.Context, order model.Order) (model.Order, error)
	DeleteOrder(ctx context.Context, id uint64) error
	UpdateOrder(ctx context.Context, id uint64, order model.Order) error
}

type orderServiceImpl struct {
	repo repository.OrderQuery
}

func NewOrderService(repo repository.OrderQuery) OrderService {
	return &orderServiceImpl{repo: repo}
}

func (u *orderServiceImpl) GetOrders(ctx context.Context) ([]model.Order, error) {
	orders, err := u.repo.GetOrder(ctx)
	if err != nil {
		return nil, err
	}
	return orders, err
}
func (u *orderServiceImpl) GetOrderById(ctx context.Context, id uint64) (model.Order, error) {
	order, err := u.repo.GetOrderById(ctx, id)
	if err != nil {
		return model.Order{}, err
	}
	return order, err
}

func (u *orderServiceImpl) CreateOrder(ctx context.Context, order model.Order) (model.Order, error) {
	createdOrder, err := u.repo.CreateOrder(ctx, order)
	if err != nil {
		return model.Order{}, err
	}
	order.OrderId = uint64(createdOrder)
	return order, nil
}

func (u *orderServiceImpl) DeleteOrder(ctx context.Context, id uint64) error {
	err := u.repo.DeleteOrder(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (u *orderServiceImpl) UpdateOrder(ctx context.Context, id uint64, order model.Order) error {
	_, err := u.GetOrderById(ctx, id)
	if err != nil {
		return err
	}

	order.OrderId = id

	err = u.repo.UpdateOrder(ctx, id, order)
	if err != nil {
		return err
	}

	return nil
}
