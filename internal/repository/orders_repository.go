package repository

import (
	"context"

	"assignment2/internal/infrastructure"
	"assignment2/internal/model"
)

type OrderQuery interface {
	GetOrder(ctx context.Context) ([]model.Order, error)
	GetOrderById(ctx context.Context, id uint64) (model.Order, error)
	CreateOrder(ctx context.Context, data model.Order) (int64, error)
	DeleteOrder(ctx context.Context, id uint64) error
	UpdateOrder(ctx context.Context, id uint64, data model.Order) error
}

type OrderCommand interface {
	CreateOrder(ctx context.Context, order model.Order) (model.Order, error)
}

type orderQueryImpl struct {
	db infrastructure.GormPostgres
}

type orderQueryMongoImpl struct{}

func NewOrderQuery(db infrastructure.GormPostgres) OrderQuery {
	return &orderQueryImpl{db: db}
}

func (u *orderQueryImpl) GetOrder(ctx context.Context) ([]model.Order, error) {
	db := u.db.GetConnection()
	var orders []model.Order
	err := db.Model(&model.Order{}).Preload("Item").Find(&orders).Error
	return orders, err
}

func (u *orderQueryImpl) GetOrderById(ctx context.Context, id uint64) (model.Order, error) {
	db := u.db.GetConnection()
	orders := model.Order{}
	err := db.Model(&model.Order{}).Preload("Item").Find(&orders).Error
	return orders, err
}

func (u *orderQueryImpl) CreateOrder(ctx context.Context, data model.Order) (int64, error) {
	db := u.db.GetConnection()

	orderQuery := `
        INSERT INTO orders (customer_name, ordered_at)
        VALUES (?, ?)
        ON CONFLICT (order_id) DO UPDATE
        SET customer_name = EXCLUDED.customer_name, ordered_at = EXCLUDED.ordered_at
        RETURNING order_id`

	var orderID int64
	err := db.Raw(orderQuery, data.CustomerName, data.OrderedAt).Row().Scan(&orderID)
	if err != nil {
		return 0, err
	}

	for i := range data.Item {
		itemQuery := `
            INSERT INTO items (order_id, item_code, description, quantity)
            VALUES (?, ?, ?, ?)
            RETURNING item_id`

		var itemID int64
		err := db.Raw(itemQuery, orderID, data.Item[i].ItemCode, data.Item[i].Description, data.Item[i].Quantity).Row().Scan(&itemID)
		if err != nil {
			return 0, err
		}
		data.Item[i].ItemId = uint(itemID)
	}

	return orderID, nil
}

func (u *orderQueryImpl) DeleteOrder(ctx context.Context, id uint64) error {
	db := u.db.GetConnection()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("order_id = ?", id).Delete(&model.Item{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("order_id = ?", id).Delete(&model.Order{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (u *orderQueryImpl) UpdateOrder(ctx context.Context, id uint64, data model.Order) error {
	db := u.db.GetConnection()

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&model.Order{}).Where("order_id = ?", id).Updates(data).Error; err != nil {
		tx.Rollback()
		return err
	}

	for i := range data.Item {
		item := data.Item[i]
		if err := tx.Model(&model.Item{}).Where("item_id = ?", item.ItemId).Updates(map[string]interface{}{
			"description": item.Description,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
