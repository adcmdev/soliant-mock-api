package shifts

import (
	"fmt"
	"time"
)

// Seed returns the default list of shifts used to populate an empty store.
func Seed() []Shift {
	today := time.Now().Local()

	countryLane := LocationInfo{
		Location:          "Country Lane Elementary School",
		CityState:         "Chicago, IL",
		AssignmentSummary: "Substitute teaching for Grade 4.",
		ClassroomDetails:  []string{"Grade 4", "20 Students"},
		Latitude:          41.8781,
		Longitude:         -87.6298,
	}

	riverdale := LocationInfo{
		Location:          "Riverdale Highschool",
		CityState:         "Chicago, IL",
		AssignmentSummary: "Art class substitute.",
		ClassroomDetails:  []string{"Art Room"},
		Latitude:          41.8781,
		Longitude:         -87.6298,
	}

	scienceLab := LocationInfo{
		Location:          "Riverdale Highschool",
		CityState:         "Chicago, IL",
		AssignmentSummary: "Urgent need for a science teacher.",
		ClassroomDetails:  []string{"Science Lab", "Grade 10"},
		Latitude:          41.8781,
		Longitude:         -87.6298,
	}

	return []Shift{
		{
			ID:           "1",
			TimeInfo:     buildTimeInfo(today, 0, "08:00:00", "16:00:00"),
			LocationInfo: countryLane,
			Log:          fullLog(today, 0),
			HourlyRate:   35.0,
		},
		{
			ID:           "2",
			TimeInfo:     buildTimeInfo(today, 1, "08:00:00", "16:00:00"),
			LocationInfo: countryLane,
			Log:          fullLog(today, 1),
			HourlyRate:   35.0,
		},
		{
			ID:           "3",
			TimeInfo:     buildTimeInfo(today, 2, "08:00:00", "16:00:00"),
			LocationInfo: countryLane,
			Log:          ShiftLog{},
			HourlyRate:   35.0,
		},
		{
			ID:           "4",
			TimeInfo:     buildTimeInfo(today, 5, "08:00:00", "16:00:00"),
			LocationInfo: riverdale,
			Log:          ShiftLog{},
			HourlyRate:   35.0,
		},
		{
			ID:           "5",
			State:        "opportunity",
			TimeInfo:     buildTimeInfo(today, 7, "09:00:00", "17:00:00"),
			LocationInfo: scienceLab,
			Log:          ShiftLog{},
			HourlyRate:   40.0,
		},
	}
}

func buildTimeInfo(base time.Time, days int, startClock, endClock string) TimeInfo {
	day := base.AddDate(0, 0, days)
	return TimeInfo{
		Date:      day.Format("Monday, January 2"),
		TimeText:  fmt.Sprintf("%s - %s", humanClock(startClock), humanClock(endClock)),
		StartTime: shiftTime(base, days, startClock),
		EndTime:   shiftTime(base, days, endClock),
	}
}

func fullLog(base time.Time, days int) ShiftLog {
	return ShiftLog{
		CheckInTime:    shiftTime(base, days, "08:00:00"),
		CheckOutTime:   shiftTime(base, days, "16:00:00"),
		BreakStartTime: shiftTime(base, days, "12:00:00"),
		BreakEndTime:   shiftTime(base, days, "12:30:00"),
	}
}

func shiftTime(base time.Time, days int, clock string) string {
	return base.AddDate(0, 0, days).Format("2006-01-02") + "T" + clock
}

func humanClock(clock string) string {
	t, err := time.Parse("15:04:05", clock)
	if err != nil {
		return clock
	}
	return t.Format("3:04 PM")
}

