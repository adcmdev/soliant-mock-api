package shifts

// Shift is the unified model that merges the original flat shift with the
// richer nested structure (time_info, location_info, log, state).
type Shift struct {
	ID           string       `json:"id"`
	State        string       `json:"state,omitempty"`
	TimeInfo     TimeInfo     `json:"time_info"`
	LocationInfo LocationInfo `json:"location_info"`
	Log          ShiftLog     `json:"log"`
	HourlyRate   float64      `json:"hourly_rate"`
}

type TimeInfo struct {
	Date      string `json:"date"`
	TimeText  string `json:"time_text"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type LocationInfo struct {
	Location          string   `json:"location"`
	CityState         string   `json:"city_state,omitempty"`
	AssignmentSummary string   `json:"assignment_summary"`
	ClassroomDetails  []string `json:"classroom_details"`
	Latitude          float64  `json:"latitude"`
	Longitude         float64  `json:"longitude"`
}

type ShiftLog struct {
	CheckInTime    string `json:"check_in_time,omitempty"`
	CheckOutTime   string `json:"check_out_time,omitempty"`
	BreakStartTime string `json:"break_start_time,omitempty"`
	BreakEndTime   string `json:"break_end_time,omitempty"`
}

// GetID implements data.Entity.
func (s *Shift) GetID() string { return s.ID }

// SetID implements data.Entity.
func (s *Shift) SetID(id string) { s.ID = id }

