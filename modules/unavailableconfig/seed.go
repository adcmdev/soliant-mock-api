package unavailableconfig

// Seed returns the default unavailable-config document used to populate an
// empty store.
func Seed() UnavailableConfig {
	return UnavailableConfig{
		RecurringWeekdays: []int{1},
		CustomRanges: []CustomRange{
			{
				StartDate: "2026-05-04T00:00:00Z",
				EndDate:   "2026-05-06T00:00:00Z",
			},
		},
	}
}

