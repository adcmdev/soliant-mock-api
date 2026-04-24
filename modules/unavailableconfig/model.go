package unavailableconfig

// UnavailableConfig is the singleton document that describes a user's
// recurring weekly unavailability and one-off unavailable date ranges.
type UnavailableConfig struct {
	RecurringWeekdays []int          `json:"recurring_weekdays"`
	CustomRanges      []CustomRange `json:"custom_ranges"`
}

// CustomRange represents a single one-off date range during which the user
// is unavailable.
type CustomRange struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

