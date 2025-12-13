# A simple REST-API in Go

## Features

- automatically generated Open-API-Spec
- Swagger-UI
- Dependency Injection for database
- Mocking database
- API-Tests
- Input Validation

## Dependencies

[Validator](https://github.com/go-playground/validator) for JSON validation

- install: `go get github.com/go-playground/validator/v10`
- usage: `https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme`

[Swaggo](https://github.com/swaggo/swag) for automatically creating Open-API Specs and display them with Swagger

- install: `go install github.com/swaggo/swag/cmd/swag@latest`
- update docs: `swag init -g ./cmd/api/main.go`
- view docs: `localhost:8090/swagger`

[Air](https://github.com/air-verse/air) for almost hot reloading the app

- install: `go install github.com/air-verse/air@latest`
- run: `air`

## Structure

```
cmd/
  api/
    main.go => entry point of the app
docs/ => api-documentation, automatically generated
internal/
  models/ => description of all structs
  server/
    server.go => server setup
    routes.go => definition of all routes
    handlers.go => definition of handlers with comments for swaggo
  store/
    store.go => interface for the storage
    memeory_store.go => for mocking the interface
    mongo_store => example for implementing the interface with a MongoDB backend
```

## Developing

- create new entities in `internal/models`
  - add tags to the fields for validation
- add new routes in `internal/server/routes.go`
  - define their behavior in `internal/server/handlers.go`
  - add comments for swaggo
  - run `swag init -g ./cmd/api/main.go` to generate API-documentation
- expand the API with first adding methods to the interface in `internal/store/store.go`
  - implement the interface in the actual stores, like `internal/store/mongo_store`

## Testing

Testing the API with `go test ./internal/server/ -v`

## TODO

- different build environments or build tags (dev vs. prod)
