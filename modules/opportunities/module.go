package opportunities

import (
	"errors"

	"soliant-mock-api/shared/data"
	"soliant-mock-api/shared/httpmod"
	"soliant-mock-api/shared/redis"
)

// ErrNotFound is returned when an opportunity does not exist in the store.
var ErrNotFound = errors.New("opportunity not found")

// New builds the opportunities module (repository + CRUD HTTP handler).
func New(cache redis.CacheRepository) *httpmod.CRUDHandler[Opportunity, *Opportunity] {
	repo := data.New[Opportunity, *Opportunity](data.Config[Opportunity]{
		Cache:    cache,
		Prefix:   "opportunities:",
		NotFound: ErrNotFound,
		Seed:     Seed,
	})

	return &httpmod.CRUDHandler[Opportunity, *Opportunity]{
		ModuleName:        "opportunities",
		BasePath:          "/opportunities",
		ResourceKey:       "opportunity",
		ResourceKeyPlural: "opportunities",
		Repo:              repo,
	}
}

