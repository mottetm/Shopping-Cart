package shopping_cart_repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/mottetm/Shopping-Cart/internal/adapters/shopping_cart_repository"
	"github.com/mottetm/Shopping-Cart/internal/services/shopping_cart"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getDb(t *testing.T, dbFile string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	require.NoError(t, err, "failed to connect database")
	return db
}

func dbFile(t *testing.T) string {
	dir := t.TempDir()
	testId := uuid.NewString()
	return fmt.Sprintf("%s/%s.db", dir, testId)
}

func TestShoppingCartRepository_CreateItem_ShouldCreateTheItem(t *testing.T) {
	dbFile := dbFile(t)
	db := getDb(t, dbFile)
	repository := shopping_cart_repository.New(db)

	ctx := context.Background()
	name := shopping_cart.ItemName("Item Name")
	quantity := shopping_cart.ItemQuantity(123)

	id, err := repository.CreateItem(ctx, name, quantity)
	require.NoError(t, err)

	item := shopping_cart.Item{}
	result := db.Find(&item).Where("id = ?", id)
	require.NoError(t, result.Error)

	require.Equal(t, item.Id, id)
	require.Equal(t, item.Name, name)
	require.Equal(t, item.Quantity, quantity)
	require.Equal(t, item.ReservationId, shopping_cart.ReservationId(""))
}

func TestShoppingCartRepository_ConfirmItem_ShouldUpdateTheItem(t *testing.T) {
	dbFile := dbFile(t)
	db := getDb(t, dbFile)
	repository := shopping_cart_repository.New(db)

	ctx := context.Background()
	name := shopping_cart.ItemName("Item Name")
	quantity := shopping_cart.ItemQuantity(123)
	reservationId := shopping_cart.ReservationId("reservation-id")

	id, err := repository.CreateItem(ctx, name, quantity)
	require.NoError(t, err)

	err = repository.ConfirmItem(ctx, id, reservationId)
	require.NoError(t, err)

	item := shopping_cart.Item{}
	result := db.Find(&item).Where("id = ?", id)
	require.NoError(t, result.Error)

	require.Equal(t, item.ReservationId, reservationId)
}

func TestShoppingCartRepository_GetItems_ShouldRetrievedConfirmedItems(t *testing.T) {
	dbFile := dbFile(t)
	db := getDb(t, dbFile)
	repository := shopping_cart_repository.New(db)

	ctx := context.Background()

	items := []shopping_cart.Item{
		{
			Name:     "Item Name 1",
			Quantity: 1,
		},
		{
			Name:          "Item Name 2",
			Quantity:      2,
			ReservationId: "reservation-id-2",
		},
		{
			Name:          "Item Name 3",
			Quantity:      3,
			ReservationId: "reservation-id-3",
		},
	}

	want := make([]shopping_cart.Item, 0)
	for _, item := range items {
		id, err := repository.CreateItem(ctx, item.Name, item.Quantity)
		require.NoError(t, err)

		if item.ReservationId != "" {
			err := repository.ConfirmItem(ctx, id, item.ReservationId)
			require.NoError(t, err)

			want = append(want, shopping_cart.Item{
				Id:            id,
				Name:          item.Name,
				Quantity:      item.Quantity,
				ReservationId: item.ReservationId,
			})
		}
	}

	got, err := repository.GetItems(ctx)
	require.NoError(t, err)

	require.Equal(t, got, want)
}
