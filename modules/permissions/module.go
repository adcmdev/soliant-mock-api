package permissions

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"soliant-mock-api/shared/redis"

	"github.com/labstack/echo/v4"
)

const (
	storeKey = "permissions:main"
	basePath = "/permissions"
)

// ErrNotFound is returned when the permissions document has not been
// initialised in the store yet.
var ErrNotFound = errors.New("permissions not found")

// Module exposes the permissions endpoints. Singleton, hence a dedicated
// Module implementation instead of the generic CRUDHandler.
type Module struct {
	cache redis.CacheRepository
	mu    sync.Mutex
}

// New builds the permissions module.
func New(cache redis.CacheRepository) *Module {
	return &Module{cache: cache}
}

// Name implements httpmod.Module.
func (m *Module) Name() string { return "permissions" }

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
	p, err := m.load()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"permissions": p})
}

// update accepts a partial JSON body and deep-merges it into the stored
// permissions document.
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

	var updated Permissions
	if err := json.Unmarshal(merged, &updated); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := m.save(updated); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"permissions": updated})
}

// --- helpers ----------------------------------------------------------------

func (m *Module) load() (Permissions, error) {
	raw, err := m.cache.Get(storeKey)
	if err != nil || len(raw) == 0 {
		return Permissions{}, ErrNotFound
	}

	var p Permissions
	if err := json.Unmarshal(raw, &p); err != nil {
		return Permissions{}, err
	}
	return p, nil
}

func (m *Module) save(p Permissions) error {
	payload, err := json.Marshal(p)
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

