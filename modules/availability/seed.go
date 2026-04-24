package availability

// Seed returns the default availability document used to populate an empty
// store.
func Seed() Availability {
	return Availability{
		IsAvailable: true,
		WeeklyDays:  []int{1, 2, 3, 4, 5},
		UpcomingTimeOff: []TimeOff{
			{
				ID:        "1",
				StartDate: "2026-04-20T00:00:00Z",
				EndDate:   "2026-04-25T00:00:00Z",
			},
		},
	}
}

