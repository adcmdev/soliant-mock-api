package home

import (
	"soliant-mock-api/modules/shifts"
	"soliant-mock-api/modules/unavailableconfig"
)

// Seed returns the default home document used to populate an empty store.
// It reuses the shared shifts.Shift struct; fields not present in that model
// (e.g. the original payload's shift_count) are intentionally omitted.
func Seed() Home {
	countryLane := shifts.LocationInfo{
		Location:          "Country Lane Elementary School",
		AssignmentSummary: "",
		ClassroomDetails:  []string{},
		Latitude:          41.8781,
		Longitude:         -87.6298,
	}

	todayShift := shifts.Shift{
		ID:    "today-1",
		State: "accepted",
		TimeInfo: shifts.TimeInfo{
			Date:      "Wednesday, Feb 15",
			TimeText:  "8:00 AM - 4:00 PM",
			StartTime: "2026-04-24T08:00:00Z",
			EndTime:   "2026-04-24T16:00:00Z",
		},
		LocationInfo: countryLane,
		Log: shifts.ShiftLog{
			CheckInTime:  "2026-04-24T08:00:00Z",
			CheckOutTime: "2026-04-24T16:00:00Z",
		},
		HourlyRate: 25.0,
	}

	upcomingShift := shifts.Shift{
		ID:    "upcoming-1",
		State: "accepted",
		TimeInfo: shifts.TimeInfo{
			Date:      "Thursday, Feb 16",
			TimeText:  "8:00 AM - 4:00 PM",
			StartTime: "2026-04-25T08:00:00Z",
			EndTime:   "2026-04-25T16:00:00Z",
		},
		LocationInfo: countryLane,
		Log: shifts.ShiftLog{
			CheckInTime:  "2026-04-25T08:00:00Z",
			CheckOutTime: "2026-04-25T16:00:00Z",
		},
		HourlyRate: 25.0,
	}

	return Home{
		UserName:      "Jamie",
		HasNewMessage: false,
		OpportunityBanner: &OpportunityBanner{
			ID:       "opp-1",
			Title:    "New Opportunity",
			Subtitle: "Take up the opportunity now",
		},
		TodayShift:            &todayShift,
		UpcomingShifts:        []shifts.Shift{upcomingShift},
		UnavailableWeekdays:   []int{},
		UnavailableDateRanges: []unavailableconfig.CustomRange{},
	}
}

