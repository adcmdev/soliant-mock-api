package httpmod

import (
	"errors"
	"net/http"

	"soliant-mock-api/shared/data"

	"github.com/labstack/echo/v4"
)

// CRUDHandler wraps a generic data.Repository and exposes it as a set of
// REST endpoints under BasePath.
//
// Response JSON payloads keep the same shape as the original shifts API:
//
//	GET    /<base>        -> { "<pluralKey>": [...] }
//	GET    /<base>/:id    -> { "<key>": {...} }
//	POST   /<base>        -> { "<key>": {...} }
//	PATCH  /<base>/:id    -> { "<key>": {...} }
//	DELETE /<base>/:id    -> 204
type CRUDHandler[T any, PT interface {
	*T
	data.Entity
}] struct {
	// ModuleName is shown in logs and returned by Name().
	ModuleName string
	// BasePath is the URL prefix, e.g. "/shifts".
	BasePath string
	// ResourceKey is the singular JSON key, e.g. "shift".
	ResourceKey string
	// ResourceKeyPlural is the plural JSON key, e.g. "shifts".
	ResourceKeyPlural string
	// Repo is the underlying generic repository.
	Repo *data.Repository[T, PT]
}

// Name implements Module.
func (h *CRUDHandler[T, PT]) Name() string { return h.ModuleName }

// SeedIfEmpty implements Module.
func (h *CRUDHandler[T, PT]) SeedIfEmpty() error { return h.Repo.SeedIfEmpty() }

// Register implements Module and attaches all CRUD routes to e.
func (h *CRUDHandler[T, PT]) Register(e *echo.Echo) {
	notFound := h.Repo.NotFound()

	e.GET(h.BasePath, func(c echo.Context) error {
		items, err := h.Repo.GetAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{h.ResourceKeyPlural: items})
	})

	e.GET(h.BasePath+"/:id", func(c echo.Context) error {
		item, err := h.Repo.GetByID(c.Param("id"))
		if err != nil {
			if errors.Is(err, notFound) {
				return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{h.ResourceKey: item})
	})

	e.POST(h.BasePath, func(c echo.Context) error {
		var body T
		if err := c.Bind(&body); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		created, err := h.Repo.Create(body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusCreated, map[string]interface{}{h.ResourceKey: created})
	})

	e.PATCH(h.BasePath+"/:id", func(c echo.Context) error {
		var body map[string]any
		if err := c.Bind(&body); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		updated, err := h.Repo.Update(c.Param("id"), body)
		if err != nil {
			if errors.Is(err, notFound) {
				return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]interface{}{h.ResourceKey: updated})
	})

	e.DELETE(h.BasePath+"/:id", func(c echo.Context) error {
		if err := h.Repo.Delete(c.Param("id")); err != nil {
			if errors.Is(err, notFound) {
				return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return c.NoContent(http.StatusNoContent)
	})
}

