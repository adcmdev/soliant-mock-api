package shifts

import (
	"errors"

	"soliant-mock-api/shared/data"
	"soliant-mock-api/shared/httpmod"
	"soliant-mock-api/shared/redis"
)

// ErrNotFound is returned when a shift does not exist in the store.
var ErrNotFound = errors.New("shift not found")

// New builds the shifts module (repository + CRUD HTTP handler).
func New(cache redis.CacheRepository) *httpmod.CRUDHandler[Shift, *Shift] {
	repo := data.New[Shift, *Shift](data.Config[Shift]{
		Cache:    cache,
		Prefix:   "shifts:",
		NotFound: ErrNotFound,
		Seed:     Seed,
	})

	return &httpmod.CRUDHandler[Shift, *Shift]{
		ModuleName:        "shifts",
		BasePath:          "/shifts",
		ResourceKey:       "shift",
		ResourceKeyPlural: "shifts",
		Repo:              repo,
	}
}

