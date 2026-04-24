package opportunities

import "soliant-mock-api/modules/shifts"

// Opportunity is an open assignment a user can take. It bundles high-level
// metadata (title, description, rate) with the list of concrete shifts that
// make up the opportunity. Shifts reuse the shared shifts.Shift model.
type Opportunity struct {
	ID                   string              `json:"id"`
	Title                string              `json:"title"`
	Description          string              `json:"description"`
	HourlyRate           float64             `json:"hourly_rate"`
	StartAt              string              `json:"start_at"`
	EndAt                string              `json:"end_at"`
	CanTakePartialShifts bool                `json:"can_take_partial_shifts"`
	ShowUrgencyNote      bool                `json:"show_urgency_note"`
	LocationInfo         shifts.LocationInfo `json:"location_info"`
	Shifts               []shifts.Shift      `json:"shifts"`
}

// GetID implements data.Entity.
func (o *Opportunity) GetID() string { return o.ID }

// SetID implements data.Entity.
func (o *Opportunity) SetID(id string) { o.ID = id }

