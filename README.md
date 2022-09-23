# Shopping-Cart
A shopping cart

## Run the server
```shell
go run ./cmd/server
```

## Execute the tests
```shell
go test -v ./...
```

## Build the docker image
```shell
docker build -f Dockerfile .
```

## Project architecture

- The core logic of the application is located in the `ShoppingCartService` implemented in [internal/services/shopping_cart/shopping_cart.go](./internal/services/shopping_cart/shopping_cart.go#L26).
- The HTTP layer is handled in the `ShoppingCartServer` implemented in [internal/servers/shopping_cart.go](./internal/servers/shopping_cart.go#L21).
- The interaction with the persistence layer (sqlite3) and the remote reservation services (stubbed) are handled in [the two adapters](./internal/adapters). They are both abstracted in the service by an interface and can easily be substituted with more mature implementations.
- The application binds the HTTP layer to an Echo server.

## Area of Improvement

- Input validation
- Async error handling
- Global error handling (all errors are currently returned as 500)
- E2E testing
- Authentication and authorization
- Logging
- Item status tracking (pending, confirmed, errored)

## Framework use

### Echo
Echo was used for handling the http layer. The choice of Echo was motivated by my knowledge of it and is simplicity.
The framework was however contained as much as possible in the App and Server layer.

### Gorm
Gorm was used for the data persistency in a sqlite3 database. Its use was motivated by ease of use for the quick project
setup. Again, the framework is wrapped as to prevent it from appearing in the business logic.