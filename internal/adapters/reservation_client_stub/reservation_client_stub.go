package reservation_client_stub

import (
	"context"

	"github.com/google/uuid"
	"github.com/mottetm/Shopping-Cart/internal/services/shopping_cart"
)

type ReservationClientStub struct{}

var _ shopping_cart.ReservationClient = (*ReservationClientStub)(nil)

func New() *ReservationClientStub {
	return &ReservationClientStub{}
}

func (c *ReservationClientStub) ReserveItem(
	ctx context.Context,
	name shopping_cart.ItemName,
) (shopping_cart.ReservationId, error) {
	return shopping_cart.ReservationId(uuid.NewString()), nil
}
