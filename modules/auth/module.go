package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"soliant-mock-api/shared/redis"

	"github.com/labstack/echo/v4"
)

const storeKey = "auth:user"

// ErrNotFound is returned when the auth user document has not been
// initialised in the store yet.
var ErrNotFound = errors.New("auth user not found")

// Module exposes the OTP / step verification endpoints. All three verify
// endpoints are mocks: regardless of the request body they respond with the
// stored auth user.
type Module struct {
	cache redis.CacheRepository
	mu    sync.Mutex
}

// New builds the auth module.
func New(cache redis.CacheRepository) *Module {
	return &Module{cache: cache}
}

// Name implements httpmod.Module.
func (m *Module) Name() string { return "auth" }

// SeedIfEmpty stores the default auth user when the key does not exist yet.
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
	// All three verification flows share the same mocked response: the stored
	// auth user wrapped in `{ "user": {...} }`.
	e.POST("/auth/code/verify", m.verify)
	e.POST("/auth/steps/verify", m.verify)
	e.POST("/shifts/detail/verify", m.verify)

	// Convenience endpoint to override the mocked auth user at runtime. Accepts
	// a partial User JSON and deep-merges it into the stored document.
	e.PUT("/auth/user", m.updateUser)
}

// --- handlers ---------------------------------------------------------------

func (m *Module) verify(c echo.Context) error {
	// The body is intentionally ignored; these are stubbed verification
	// endpoints that always succeed.
	u, err := m.load()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"user": u})
}

func (m *Module) updateUser(c echo.Context) error {
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

	delete(patch, "uuid")
	deepMerge(current, patch)

	merged, err := json.Marshal(current)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var updated User
	if err := json.Unmarshal(merged, &updated); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	updated.UUID = existing.UUID

	if err := m.save(updated); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"user": updated})
}

// --- helpers ----------------------------------------------------------------

func (m *Module) load() (User, error) {
	raw, err := m.cache.Get(storeKey)
	if err != nil || len(raw) == 0 {
		return User{}, ErrNotFound
	}

	var u User
	if err := json.Unmarshal(raw, &u); err != nil {
		return User{}, err
	}
	return u, nil
}

func (m *Module) save(u User) error {
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

