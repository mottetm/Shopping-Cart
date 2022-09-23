package servers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/mottetm/Shopping-Cart/internal/servers"
	servers_mocks "github.com/mottetm/Shopping-Cart/internal/servers/mocks"
	"github.com/mottetm/Shopping-Cart/internal/services/shopping_cart"
	"github.com/stretchr/testify/require"
)

type shoppingCartServerTestSetup struct {
	service *servers_mocks.MockShoppingCartService
	server  *servers.ShoppingCartServer
}

func setupShoppingCartServerTest(t *testing.T) *shoppingCartServerTestSetup {
	ctrl := gomock.NewController(t)

	service := servers_mocks.NewMockShoppingCartService(ctrl)

	server := servers.NewShoppingCartServer(service)

	return &shoppingCartServerTestSetup{service, server}
}

func TestShoppingCartServer_PostItems_ShouldCreateItemsAndReturnLocation(t *testing.T) {
	setup := setupShoppingCartServerTest(t)

	name := "Item Name"
	quantity := 1
	id := shopping_cart.ItemId(1)
	body, err := json.Marshal(
		struct {
			Name     string `json:"name"`
			Quantity int    `json:"quantity"`
		}{
			Name:     name,
			Quantity: quantity,
		},
	)
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://localhost/items", bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("content-type", "application/json")

	resp := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, resp)

	setup.service.EXPECT().
		PostItem(
			req.Context(),
			shopping_cart.Item{
				Name:     shopping_cart.ItemName(name),
				Quantity: shopping_cart.ItemQuantity(quantity),
			},
		).
		Return(id, nil)

	err = setup.server.PostItem(ctx)
	require.NoError(t, err)

	require.Equal(t, resp.Code, http.StatusCreated)
	require.Equal(t, resp.Header().Get("location"), fmt.Sprintf("/items/%s", id))
}

func TestShoppingCartServer_GetItems_ShouldSendBackItems(t *testing.T) {
	setup := setupShoppingCartServerTest(t)

	items := []shopping_cart.Item{
		{
			Id:            1,
			Name:          "Item Name 1",
			Quantity:      1,
			ReservationId: "reservation-id-1",
		},
		{
			Id:            2,
			Name:          "Item Name 2",
			Quantity:      1,
			ReservationId: "reservation-id-2",
		},
		{
			Id:            3,
			Name:          "Item Name 3",
			Quantity:      1,
			ReservationId: "reservation-id-3",
		},
	}
	want, _ := json.Marshal(items)

	req, err := http.NewRequest(http.MethodGet, "http://localhost/items", nil)
	require.NoError(t, err)

	resp := httptest.NewRecorder()
	ctx := echo.New().NewContext(req, resp)

	setup.service.EXPECT().
		GetItems(req.Context()).
		Return(items, nil)

	err = setup.server.GetItems(ctx)
	require.NoError(t, err)

	require.Equal(t, resp.Code, http.StatusOK)

	got, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Equal(t, got, append(want, '\n'))
}
