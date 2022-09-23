package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mottetm/Shopping-Cart/internal/adapters/reservation_client_stub"
	"github.com/mottetm/Shopping-Cart/internal/adapters/shopping_cart_repository"
	"github.com/mottetm/Shopping-Cart/internal/app"
	"github.com/mottetm/Shopping-Cart/internal/servers"
	"github.com/mottetm/Shopping-Cart/internal/services/shopping_cart"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	dbFile := os.Getenv("SHOPPING_CART_DB_FILE")
	if dbFile == "" {
		dbFile = "shopping-cart.db"
	}

	port := os.Getenv("SHOPPING_CART_PORT")
	if port == "" {
		port = "8080"
	}

	host := os.Getenv("SHOPPING_CART_HOST")
	if host == "" {
		host = "localhost"
	}

	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	shoppingCartRepository := shopping_cart_repository.New(db)
	reservationClient := reservation_client_stub.New()
	shoppingCartService := shopping_cart.New(shoppingCartRepository, reservationClient)
	shoppingCartServer := servers.NewShoppingCartServer(shoppingCartService)

	app := app.New(
		app.Config{
			Server: shoppingCartServer,
			Host:   host,
			Port:   port,
		},
	)

	app.Start()

	if err != nil {
		panic(err)
	}

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	defer close(stopCh)
	<-stopCh

	if err := app.Stop(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
