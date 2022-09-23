package shopping_cart

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

// nolint
//go:generate mockgen -destination=./mocks/repository.go -package=shopping_cart_mocks github.com/mottetm/Shopping-Cart/internal/services/shopping_cart ShoppingCartRepository

type ShoppingCartRepository interface {
	CreateItem(context.Context, ItemName, ItemQuantity) (ItemId, error)
	ConfirmItem(context.Context, ItemId, ReservationId) error
	GetItems(context.Context) ([]Item, error)
}

// nolint
//go:generate mockgen -destination=./mocks/client.go -package=shopping_cart_mocks github.com/mottetm/Shopping-Cart/internal/services/shopping_cart ReservationClient

type ReservationClient interface {
	ReserveItem(context.Context, ItemName) (ReservationId, error)
}

type ShoppingCartService struct {
	repository        ShoppingCartRepository
	reservationClient ReservationClient
}

func New(
	repository ShoppingCartRepository,
	reservationClient ReservationClient,
) *ShoppingCartService {
	return &ShoppingCartService{repository, reservationClient}
}

func (s *ShoppingCartService) PostItem(
	ctx context.Context,
	item Item,
) (ItemId, error) {
	id, err := s.repository.CreateItem(ctx, item.Name, item.Quantity)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	go func() {
		reservationId, err := s.reservationClient.ReserveItem(ctx, item.Name)
		if err != nil {
			fmt.Println(err)
			// TODO: handle async errors
		}

		err = s.repository.ConfirmItem(ctx, id, reservationId)
		if err != nil {
			fmt.Println(err)
			// TODO: handle async errors
		}
	}()

	return id, nil
}

func (s *ShoppingCartService) GetItems(
	ctx context.Context,
) ([]Item, error) {
	return s.repository.GetItems(ctx)
}
