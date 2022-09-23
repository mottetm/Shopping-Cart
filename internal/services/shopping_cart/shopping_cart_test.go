package shopping_cart_test

import (
	"context"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mottetm/Shopping-Cart/internal/services/shopping_cart"
	shopping_cart_mocks "github.com/mottetm/Shopping-Cart/internal/services/shopping_cart/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type shoppingCartTestSetup struct {
	repository *shopping_cart_mocks.MockShoppingCartRepository
	client     *shopping_cart_mocks.MockReservationClient
	service    *shopping_cart.ShoppingCartService
}

func setupShoppingCartTest(t *testing.T) *shoppingCartTestSetup {
	ctrl := gomock.NewController(t)

	repository := shopping_cart_mocks.NewMockShoppingCartRepository(ctrl)
	client := shopping_cart_mocks.NewMockReservationClient(ctrl)

	service := shopping_cart.New(repository, client)

	return &shoppingCartTestSetup{repository, client, service}
}

func TestShoppingCartService_PostItems_ShouldCreateAndReserveItem(t *testing.T) {
	setup := setupShoppingCartTest(t)

	ctx := context.Background()
	want := shopping_cart.ItemId(1)
	reservationId := shopping_cart.ReservationId("reservation-id")
	item := shopping_cart.Item{
		Name:     "Item Name",
		Quantity: 123,
	}

	wg := new(sync.WaitGroup)
	wg.Add(1)

	setup.repository.EXPECT().
		CreateItem(ctx, item.Name, item.Quantity).
		Return(want, nil)

	setup.client.EXPECT().
		ReserveItem(ctx, item.Name).
		Return(reservationId, nil)

	setup.repository.EXPECT().
		ConfirmItem(ctx, want, reservationId).
		Do(func(context.Context, shopping_cart.ItemId, shopping_cart.ReservationId) {
			wg.Done()
		}).
		Return(nil)

	got, err := setup.service.PostItem(ctx, item)

	wg.Wait()

	require.Nil(t, err, "PostItem failed")

	require.Equal(t, got, want)
}

func TestShoppingCartService_PostItems_ShouldFailAndReturnErrors(t *testing.T) {
	setup := setupShoppingCartTest(t)

	ctx := context.Background()
	want := errors.New("something happened")
	item := shopping_cart.Item{
		Name:     "Item Name",
		Quantity: 123,
	}

	setup.repository.EXPECT().
		CreateItem(ctx, item.Name, item.Quantity).
		Return(shopping_cart.ItemId(0), want)

	_, got := setup.service.PostItem(ctx, item)

	require.ErrorIs(t, got, want)
}

func TestShoppingCartService_GetItems_ShouldGetAndReturnItems(t *testing.T) {
	setup := setupShoppingCartTest(t)

	ctx := context.Background()
	want := []shopping_cart.Item{
		{
			Id:            1,
			Name:          "Item Name 1",
			Quantity:      123,
			ReservationId: "reservation-id-1",
		},
		{
			Id:            2,
			Name:          "Item Name 2",
			Quantity:      234,
			ReservationId: "reservation-id-2",
		},
	}

	setup.repository.EXPECT().
		GetItems(ctx).
		Return(want, nil)

	got, err := setup.service.GetItems(ctx)

	require.Nil(t, err, "GetItems failed")

	require.Equal(t, got, want)
}

func TestShoppingCartService_GetItems_ShouldGetAndReturnErrors(t *testing.T) {
	setup := setupShoppingCartTest(t)

	ctx := context.Background()
	want := errors.New("something happened")

	setup.repository.EXPECT().
		GetItems(ctx).
		Return(nil, want)

	_, got := setup.service.GetItems(ctx)

	require.ErrorIs(t, got, want)
}
