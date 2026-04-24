package profile

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"soliant-mock-api/shared/redis"

	"github.com/labstack/echo/v4"
)

const (
	storeKey = "profile:main"
	basePath = "/profile"
)

// ErrNotFound is returned when the profile document has not been initialised
// in the store yet.
var ErrNotFound = errors.New("profile not found")

// Module exposes the profile endpoints. Profile is a singleton so it uses a
// dedicated Module implementation instead of the generic CRUDHandler.
type Module struct {
	cache redis.CacheRepository
	mu    sync.Mutex
}

// New builds the profile module.
func New(cache redis.CacheRepository) *Module {
	return &Module{cache: cache}
}

// Name implements httpmod.Module.
func (m *Module) Name() string { return "profile" }

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
	return c.JSON(http.StatusOK, map[string]interface{}{"profile": p})
}

// update accepts a partial JSON body and deep-merges it into the stored
// profile document. The user.uuid is preserved regardless of the patch.
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

	// Strip any attempt to mutate the immutable user uuid.
	if userPatch, ok := patch["user"].(map[string]any); ok {
		delete(userPatch, "uuid")
	}
	deepMerge(current, patch)

	merged, err := json.Marshal(current)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var updated Profile
	if err := json.Unmarshal(merged, &updated); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	updated.User.UUID = existing.User.UUID

	if err := m.save(updated); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"profile": updated})
}

// --- helpers ----------------------------------------------------------------

func (m *Module) load() (Profile, error) {
	raw, err := m.cache.Get(storeKey)
	if err != nil || len(raw) == 0 {
		return Profile{}, ErrNotFound
	}

	var p Profile
	if err := json.Unmarshal(raw, &p); err != nil {
		return Profile{}, err
	}
	return p, nil
}

func (m *Module) save(p Profile) error {
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

