package shopping_cart_repository

import (
	"context"

	"github.com/mottetm/Shopping-Cart/internal/services/shopping_cart"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type item struct {
	gorm.Model
	ID            int    `gorm:"autoIncrement"`
	Name          string `gorm:"not null"`
	Quantity      int    `gorm:"not null"`
	ReservationId string
}

type ShoppingCartRepository struct {
	db *gorm.DB
}

var _ shopping_cart.ShoppingCartRepository = (*ShoppingCartRepository)(nil)

func New(db *gorm.DB) *ShoppingCartRepository {
	db.AutoMigrate(&item{})
	return &ShoppingCartRepository{db}
}

func (r *ShoppingCartRepository) CreateItem(
	ctx context.Context,
	name shopping_cart.ItemName,
	quantity shopping_cart.ItemQuantity,
) (shopping_cart.ItemId, error) {
	item := &item{
		Name:     string(name),
		Quantity: int(quantity),
	}
	result := r.db.Create(item)
	if result.Error != nil {
		return 0, errors.WithStack(result.Error)
	}
	return shopping_cart.ItemId(item.ID), nil
}

func (r *ShoppingCartRepository) ConfirmItem(
	ctx context.Context,
	id shopping_cart.ItemId,
	reservationId shopping_cart.ReservationId,
) error {
	result := r.db.Model(&item{}).
		Where("id = ?", id).
		Update("reservation_id", reservationId)
	return result.Error
}

func (r *ShoppingCartRepository) GetItems(
	ctx context.Context,
) ([]shopping_cart.Item, error) {
	data := make([]item, 0)
	result := r.db.Where("reservation_id <> ?", "").Find(&data)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}
	items := make([]shopping_cart.Item, 0, len(data))
	for _, item := range data {
		items = append(items, shopping_cart.Item{
			Id:            shopping_cart.ItemId(item.ID),
			Name:          shopping_cart.ItemName(item.Name),
			Quantity:      shopping_cart.ItemQuantity(item.Quantity),
			ReservationId: shopping_cart.ReservationId(item.ReservationId),
		})
	}
	return items, nil
}
