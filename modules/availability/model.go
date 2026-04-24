package availability

// Availability is the singleton document that describes the user's general
// availability plus any upcoming time-off entries.
type Availability struct {
	IsAvailable     bool      `json:"is_available"`
	WeeklyDays      []int     `json:"weekly_days"`
	UpcomingTimeOff []TimeOff `json:"upcoming_time_off"`
}

// TimeOff is a single period the user won't be available.
type TimeOff struct {
	ID        string `json:"id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

