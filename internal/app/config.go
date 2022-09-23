package app

import "github.com/mottetm/Shopping-Cart/internal/servers"

type Config struct {
	Server *servers.ShoppingCartServer
	Host   string
	Port   string
}
