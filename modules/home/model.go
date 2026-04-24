package home

import (
	"soliant-mock-api/modules/shifts"
	"soliant-mock-api/modules/unavailableconfig"
)

// Home is the singleton aggregate document returned by GET /home.
// today_shift and upcoming_shifts reuse the shared shifts.Shift model so the
// shape stays consistent with the /shifts endpoints.
type Home struct {
	UserName              string                          `json:"user_name"`
	HasNewMessage         bool                            `json:"has_new_message"`
	OpportunityBanner     *OpportunityBanner              `json:"opportunity_banner,omitempty"`
	TodayShift            *shifts.Shift                   `json:"today_shift,omitempty"`
	UpcomingShifts        []shifts.Shift                  `json:"upcoming_shifts"`
	UnavailableWeekdays   []int                           `json:"unavailable_weekdays"`
	UnavailableDateRanges []unavailableconfig.CustomRange `json:"unavailable_date_ranges"`
}

// OpportunityBanner is a compact teaser shown on the home screen linking to
// an opportunity.
type OpportunityBanner struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

