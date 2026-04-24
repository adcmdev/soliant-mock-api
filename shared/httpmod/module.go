// Package httpmod provides a generic CRUD handler that turns any
// data.Repository into a standard set of REST endpoints. Each "module" in
// this project is just a thin wrapper around CRUDHandler that knows its own
// model, seed data and base path.
package httpmod

import "github.com/labstack/echo/v4"

// Module is the contract every feature module must satisfy so the main
// application can bootstrap it uniformly.
type Module interface {
	// Name is a short human-readable identifier, useful for logs.
	Name() string
	// SeedIfEmpty populates the data store with default values when empty.
	SeedIfEmpty() error
	// Register attaches the module routes to the given Echo instance.
	Register(e *echo.Echo)
}

