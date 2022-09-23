package servers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mottetm/Shopping-Cart/internal/services/shopping_cart"
	"github.com/pkg/errors"
)

type ShoppingCartService interface {
	PostItem(context.Context, shopping_cart.Item) (shopping_cart.ItemId, error)
	GetItems(ctx context.Context) ([]shopping_cart.Item, error)
}

// nolint
//go:generate mockgen -destination=./mocks/service.go -package=servers_mocks github.com/mottetm/Shopping-Cart/internal/servers ShoppingCartService

type ShoppingCartServer struct {
	service ShoppingCartService
}

func NewShoppingCartServer(
	service ShoppingCartService,
) *ShoppingCartServer {
	return &ShoppingCartServer{service}
}

func (s *ShoppingCartServer) PostItem(c echo.Context) error {
	item := shopping_cart.Item{}
	c.Bind(&item)

	id, err := s.service.PostItem(c.Request().Context(), item)
	if err != nil {
		return errors.WithStack(err)
	}

	c.Response().Header().Set("location", fmt.Sprintf("/items/%s", id))
	c.Response().WriteHeader(http.StatusCreated)

	return nil
}

func (s *ShoppingCartServer) GetItems(c echo.Context) error {
	items, err := s.service.GetItems(c.Request().Context())
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(c.JSON(http.StatusOK, items))
}
