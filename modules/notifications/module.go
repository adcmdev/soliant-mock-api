package notifications

import (
	"errors"

	"soliant-mock-api/shared/data"
	"soliant-mock-api/shared/httpmod"
	"soliant-mock-api/shared/redis"
)

// ErrNotFound is returned when a notification does not exist in the store.
var ErrNotFound = errors.New("notification not found")

// New builds the notifications module (repository + CRUD HTTP handler).
func New(cache redis.CacheRepository) *httpmod.CRUDHandler[Notification, *Notification] {
	repo := data.New[Notification, *Notification](data.Config[Notification]{
		Cache:    cache,
		Prefix:   "notifications:",
		NotFound: ErrNotFound,
		Seed:     Seed,
	})

	return &httpmod.CRUDHandler[Notification, *Notification]{
		ModuleName:        "notifications",
		BasePath:          "/notifications",
		ResourceKey:       "notification",
		ResourceKeyPlural: "notifications",
		Repo:              repo,
	}
}

