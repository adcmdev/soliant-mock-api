// Package pastshifts exposes a read-only endpoint that returns the stored
// shifts grouped as a map keyed by date (YYYY-MM-DD), matching the mobile
// client's Map<String, List<Shift>> contract.
package pastshifts

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"soliant-mock-api/modules/shifts"
	"soliant-mock-api/shared/data"

	"github.com/labstack/echo/v4"
)

const basePath = "/shifts/past"

// Module exposes GET /shifts/past. It has no storage of its own; the shifts
// repository is the source of truth.
type Module struct {
	repo *data.Repository[shifts.Shift, *shifts.Shift]
}

// New builds the past-shifts module. The repo argument must be the shifts
// module's underlying repository so both endpoints stay in sync.
func New(repo *data.Repository[shifts.Shift, *shifts.Shift]) *Module {
	return &Module{repo: repo}
}

// Name implements httpmod.Module.
func (m *Module) Name() string { return "past-shifts" }

// SeedIfEmpty is a no-op: past shifts are derived from the shifts collection.
func (m *Module) SeedIfEmpty() error { return nil }

// Register implements httpmod.Module.
func (m *Module) Register(e *echo.Echo) {
	e.GET(basePath, m.get)
}

// --- handlers ---------------------------------------------------------------

func (m *Module) get(c echo.Context) error {
	all, err := m.repo.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	now := time.Now().UTC()
	grouped := map[string][]shifts.Shift{}

	for _, s := range all {
		if !isPast(s, now) {
			continue
		}

		key := dateKey(s.TimeInfo.StartTime)
		if key == "" {
			continue
		}
		grouped[key] = append(grouped[key], s)
	}

	// Sort each day's shifts by start time for deterministic output.
	for _, list := range grouped {
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].TimeInfo.StartTime < list[j].TimeInfo.StartTime
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"past_shifts": grouped})
}

// --- helpers ----------------------------------------------------------------

// dateKey extracts the YYYY-MM-DD portion from a shift timestamp. The stored
// values may be either "2026-04-24T08:00:00" or "2026-04-24T08:00:00Z"; we
// simply take whatever precedes the 'T'.
func dateKey(ts string) string {
	if ts == "" {
		return ""
	}
	if idx := strings.Index(ts, "T"); idx > 0 {
		return ts[:idx]
	}
	return ts
}

// isPast returns true when the shift has already ended. If the end time can't
// be parsed we include the shift by default — it's a mock, not a filter we
// want to be picky about.
func isPast(s shifts.Shift, now time.Time) bool {
	end := s.TimeInfo.EndTime
	if end == "" {
		return false
	}

	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05"} {
		if t, err := time.Parse(layout, end); err == nil {
			return t.Before(now)
		}
	}
	return true
}

