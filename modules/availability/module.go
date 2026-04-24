package availability

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"sync"

	"soliant-mock-api/shared/redis"

	"github.com/labstack/echo/v4"
)

const storeKey = "availability:main"

// ErrNotFound is returned when the availability document has not been
// initialised in the store yet.
var ErrNotFound = errors.New("availability not found")

// ErrTimeOffNotFound is returned when a time-off entry does not exist.
var ErrTimeOffNotFound = errors.New("time off not found")

// Module exposes the availability endpoints. It does not use the generic
// CRUDHandler because availability is a singleton with a nested sub-collection
// (time-off) rather than a flat list of items.
type Module struct {
	cache redis.CacheRepository
	mu    sync.Mutex
}

// New builds the availability module.
func New(cache redis.CacheRepository) *Module {
	return &Module{cache: cache}
}

// Name implements httpmod.Module.
func (m *Module) Name() string { return "availability" }

// SeedIfEmpty stores the default availability document when the key does not
// exist yet.
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
	e.GET("/availability", m.get)
	e.PUT("/availability", m.update)

	e.POST("/availability/time-off", m.createTimeOff)
	e.PUT("/availability/time-off/:id", m.updateTimeOff)
	e.DELETE("/availability/time-off/:id", m.deleteTimeOff)
}

// --- handlers ---------------------------------------------------------------

func (m *Module) get(c echo.Context) error {
	a, err := m.load()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"availability": a})
}

// update accepts a partial JSON body and deep-merges it into the stored
// availability document. upcoming_time_off can only be mutated through its
// own endpoints.
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

	delete(patch, "upcoming_time_off")
	deepMerge(current, patch)

	merged, err := json.Marshal(current)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var updated Availability
	if err := json.Unmarshal(merged, &updated); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := m.save(updated); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"availability": updated})
}

func (m *Module) createTimeOff(c echo.Context) error {
	var body TimeOff
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	a, err := m.load()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if body.ID == "" {
		body.ID = nextTimeOffID(a.UpcomingTimeOff)
	} else {
		for _, t := range a.UpcomingTimeOff {
			if t.ID == body.ID {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "time off " + body.ID + " already exists",
				})
			}
		}
	}

	a.UpcomingTimeOff = append(a.UpcomingTimeOff, body)
	if err := m.save(a); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{"time_off": body})
}

func (m *Module) updateTimeOff(c echo.Context) error {
	id := c.Param("id")

	var patch map[string]any
	if err := c.Bind(&patch); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	a, err := m.load()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	idx := -1
	for i, t := range a.UpcomingTimeOff {
		if t.ID == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": ErrTimeOffNotFound.Error()})
	}

	raw, err := json.Marshal(a.UpcomingTimeOff[idx])
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var current map[string]any
	if err := json.Unmarshal(raw, &current); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	delete(patch, "id")
	deepMerge(current, patch)

	merged, err := json.Marshal(current)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var updated TimeOff
	if err := json.Unmarshal(merged, &updated); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	updated.ID = id

	a.UpcomingTimeOff[idx] = updated
	if err := m.save(a); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"time_off": updated})
}

func (m *Module) deleteTimeOff(c echo.Context) error {
	id := c.Param("id")

	m.mu.Lock()
	defer m.mu.Unlock()

	a, err := m.load()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	idx := -1
	for i, t := range a.UpcomingTimeOff {
		if t.ID == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": ErrTimeOffNotFound.Error()})
	}

	a.UpcomingTimeOff = append(a.UpcomingTimeOff[:idx], a.UpcomingTimeOff[idx+1:]...)
	if err := m.save(a); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// --- helpers ----------------------------------------------------------------

func (m *Module) load() (Availability, error) {
	raw, err := m.cache.Get(storeKey)
	if err != nil || len(raw) == 0 {
		return Availability{}, ErrNotFound
	}

	var a Availability
	if err := json.Unmarshal(raw, &a); err != nil {
		return Availability{}, err
	}
	if a.WeeklyDays == nil {
		a.WeeklyDays = []int{}
	}
	if a.UpcomingTimeOff == nil {
		a.UpcomingTimeOff = []TimeOff{}
	}
	return a, nil
}

func (m *Module) save(a Availability) error {
	if a.WeeklyDays == nil {
		a.WeeklyDays = []int{}
	}
	if a.UpcomingTimeOff == nil {
		a.UpcomingTimeOff = []TimeOff{}
	}

	payload, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return m.cache.Set(storeKey, payload, 0)
}

func nextTimeOffID(items []TimeOff) string {
	maxNum := 0
	for _, t := range items {
		if n, err := strconv.Atoi(t.ID); err == nil && n > maxNum {
			maxNum = n
		}
	}
	return strconv.Itoa(maxNum + 1)
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

