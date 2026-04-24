package unavailableconfig

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"soliant-mock-api/shared/redis"

	"github.com/labstack/echo/v4"
)

const (
	storeKey = "unavailable-config:main"
	basePath = "/schedule/unavailable-config"
)

// ErrNotFound is returned when the unavailable-config document has not been
// initialised in the store yet.
var ErrNotFound = errors.New("unavailable config not found")

// Module exposes the unavailable-config endpoints. It is a singleton, hence
// a dedicated Module implementation instead of the generic CRUDHandler.
type Module struct {
	cache redis.CacheRepository
	mu    sync.Mutex
}

// New builds the unavailable-config module.
func New(cache redis.CacheRepository) *Module {
	return &Module{cache: cache}
}

// Name implements httpmod.Module.
func (m *Module) Name() string { return "unavailable-config" }

// SeedIfEmpty stores the default document when the key does not exist yet.
func (m *Module) SeedIfEmpty() error {
	exists, err := m.cache.Exists(storeKey)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return m.save(Seed())
}

// Register implements httpmod.Module.
func (m *Module) Register(e *echo.Echo) {
	e.GET(basePath, m.get)
	e.PUT(basePath, m.update)
}

// --- handlers ---------------------------------------------------------------

func (m *Module) get(c echo.Context) error {
	u, err := m.load()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"unavailable_config": u})
}

// update accepts a partial JSON body and deep-merges it into the stored
// document. Arrays (recurring_weekdays, custom_ranges) are replaced wholesale
// when present in the patch.
func (m *Module) update(c echo.Context) error {
	var patch map[string]any
	if err := c.Bind(&patch); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	existing, err := m.load()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	raw, err := json.Marshal(existing)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var current map[string]any
	if err := json.Unmarshal(raw, &current); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	deepMerge(current, patch)

	merged, err := json.Marshal(current)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var updated UnavailableConfig
	if err := json.Unmarshal(merged, &updated); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := m.save(updated); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"unavailable_config": updated})
}

// --- helpers ----------------------------------------------------------------

func (m *Module) load() (UnavailableConfig, error) {
	raw, err := m.cache.Get(storeKey)
	if err != nil || len(raw) == 0 {
		return UnavailableConfig{}, ErrNotFound
	}

	var u UnavailableConfig
	if err := json.Unmarshal(raw, &u); err != nil {
		return UnavailableConfig{}, err
	}
	if u.RecurringWeekdays == nil {
		u.RecurringWeekdays = []int{}
	}
	if u.CustomRanges == nil {
		u.CustomRanges = []CustomRange{}
	}
	return u, nil
}

func (m *Module) save(u UnavailableConfig) error {
	if u.RecurringWeekdays == nil {
		u.RecurringWeekdays = []int{}
	}
	if u.CustomRanges == nil {
		u.CustomRanges = []CustomRange{}
	}

	payload, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return m.cache.Set(storeKey, payload, 0)
}

// deepMerge recursively merges src into dst. Nested objects are merged;
// primitive values and arrays from src replace the corresponding value in dst.
func deepMerge(dst, src map[string]any) {
	for key, value := range src {
		if srcMap, ok := value.(map[string]any); ok {
			if dstMap, ok := dst[key].(map[string]any); ok {
				deepMerge(dstMap, srcMap)
				continue
			}
		}
		dst[key] = value
	}
}

