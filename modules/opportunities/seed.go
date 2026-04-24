package opportunities

import (
	"fmt"
	"strconv"
	"time"

	"soliant-mock-api/modules/shifts"
)

// Seed returns 10 hand-crafted opportunities with between 1 and 10 shifts
// each. The shift count per opportunity is deterministic (not random) so the
// seed is reproducible across restarts, but the counts are intentionally
// varied to exercise different list sizes.
func Seed() []Opportunity {
	// Each entry defines the opportunity-level metadata plus how many shifts
	// it should contain and the base hourly rate used for its shifts.
	specs := []struct {
		title         string
		description   string
		hourlyRate    float64
		shiftRate     float64
		startDate     time.Time
		days          int // number of shifts == number of consecutive days
		startClock    string
		endClock      string
		partial       bool
		urgent        bool
		location      shifts.LocationInfo
		withBreak     bool
	}{
		{
			title:       "SAI Program Support — School A",
			description: "You'll be providing support for 3 different students enrolled in our SAI Program.",
			hourlyRate:  50.0,
			shiftRate:   20.0,
			startDate:   date(2026, 5, 4),
			days:        3,
			startClock:  "09:00:00",
			endClock:    "17:00:00",
			partial:     true,
			urgent:      true,
			location: shifts.LocationInfo{
				Location:          "School A",
				CityState:         "San Francisco, CA",
				AssignmentSummary: "You'll be providing support for 3 different students enrolled in our SAI Program...",
				ClassroomDetails:  []string{"Grade 4", "20 Students"},
				Latitude:          37.7749,
				Longitude:         -122.4194,
			},
			withBreak: true,
		},
		{
			title:       "Physics Substitute — Lincoln High",
			description: "Cover AP Physics classes while the regular teacher is on leave.",
			hourlyRate:  45.0,
			shiftRate:   22.0,
			startDate:   date(2026, 5, 6),
			days:        1,
			startClock:  "08:00:00",
			endClock:    "15:30:00",
			partial:     false,
			urgent:      false,
			location: shifts.LocationInfo{
				Location:          "Lincoln High School",
				CityState:         "Oakland, CA",
				AssignmentSummary: "Teach two sections of AP Physics using provided lesson plans.",
				ClassroomDetails:  []string{"Grade 11-12", "Physics Lab"},
				Latitude:          37.8044,
				Longitude:         -122.2712,
			},
			withBreak: true,
		},
		{
			title:       "Middle School Math Tutor",
			description: "Small-group tutoring sessions focused on pre-algebra and algebra 1.",
			hourlyRate:  38.0,
			shiftRate:   25.0,
			startDate:   date(2026, 5, 11),
			days:        5,
			startClock:  "15:00:00",
			endClock:    "18:00:00",
			partial:     true,
			urgent:      false,
			location: shifts.LocationInfo{
				Location:          "Westview Middle School",
				CityState:         "Berkeley, CA",
				AssignmentSummary: "After-school tutoring for 6th-8th graders.",
				ClassroomDetails:  []string{"Grade 6-8", "Small Group"},
				Latitude:          37.8716,
				Longitude:         -122.2727,
			},
			withBreak: false,
		},
		{
			title:       "Music Teacher — Elementary",
			description: "Lead general music classes, K-5, using the school's Orff instruments.",
			hourlyRate:  42.0,
			shiftRate:   21.0,
			startDate:   date(2026, 5, 18),
			days:        2,
			startClock:  "08:30:00",
			endClock:    "14:30:00",
			partial:     false,
			urgent:      false,
			location: shifts.LocationInfo{
				Location:          "Maple Grove Elementary",
				CityState:         "San Jose, CA",
				AssignmentSummary: "General music classes, Grades K-5.",
				ClassroomDetails:  []string{"K-5", "Music Room"},
				Latitude:          37.3382,
				Longitude:         -121.8863,
			},
			withBreak: true,
		},
		{
			title:       "Library Aide — Two Week Coverage",
			description: "Support the librarian during a two-week absence.",
			hourlyRate:  30.0,
			shiftRate:   18.0,
			startDate:   date(2026, 5, 25),
			days:        10,
			startClock:  "09:00:00",
			endClock:    "16:00:00",
			partial:     true,
			urgent:      false,
			location: shifts.LocationInfo{
				Location:          "Riverdale Public Library (School)",
				CityState:         "Chicago, IL",
				AssignmentSummary: "Help students check out books, shelve returns, and run story time.",
				ClassroomDetails:  []string{"K-8", "Main Library"},
				Latitude:          41.8781,
				Longitude:         -87.6298,
			},
			withBreak: true,
		},
		{
			title:       "Lunch Monitor",
			description: "Supervise students during lunch and recess.",
			hourlyRate:  22.0,
			shiftRate:   15.0,
			startDate:   date(2026, 6, 1),
			days:        4,
			startClock:  "11:30:00",
			endClock:    "13:30:00",
			partial:     true,
			urgent:      true,
			location: shifts.LocationInfo{
				Location:          "Sunset Elementary",
				CityState:         "San Francisco, CA",
				AssignmentSummary: "Lunchroom supervision and playground monitoring.",
				ClassroomDetails:  []string{"K-5", "Cafeteria"},
				Latitude:          37.7544,
				Longitude:         -122.4869,
			},
			withBreak: false,
		},
		{
			title:       "Chemistry Teacher — High School",
			description: "Full week coverage for Chemistry 1 and Honors Chemistry.",
			hourlyRate:  55.0,
			shiftRate:   28.0,
			startDate:   date(2026, 6, 8),
			days:        7,
			startClock:  "07:45:00",
			endClock:    "15:15:00",
			partial:     false,
			urgent:      true,
			location: shifts.LocationInfo{
				Location:          "Jefferson High School",
				CityState:         "Portland, OR",
				AssignmentSummary: "Run the existing lab schedule; safety briefings provided.",
				ClassroomDetails:  []string{"Grade 10-11", "Chemistry Lab"},
				Latitude:          45.5152,
				Longitude:         -122.6784,
			},
			withBreak: true,
		},
		{
			title:       "Kindergarten Substitute",
			description: "Follow the lead teacher's plans for kindergarten classroom.",
			hourlyRate:  36.0,
			shiftRate:   19.0,
			startDate:   date(2026, 6, 15),
			days:        6,
			startClock:  "08:15:00",
			endClock:    "14:45:00",
			partial:     true,
			urgent:      false,
			location: shifts.LocationInfo{
				Location:          "Country Lane Elementary School",
				CityState:         "Chicago, IL",
				AssignmentSummary: "Kindergarten substitute, structured lesson plans provided.",
				ClassroomDetails:  []string{"Kindergarten", "18 Students"},
				Latitude:          41.8781,
				Longitude:         -87.6298,
			},
			withBreak: true,
		},
		{
			title:       "Special Education Aide",
			description: "1:1 support for students with IEPs across multiple classrooms.",
			hourlyRate:  48.0,
			shiftRate:   24.0,
			startDate:   date(2026, 6, 22),
			days:        8,
			startClock:  "08:00:00",
			endClock:    "15:00:00",
			partial:     true,
			urgent:      true,
			location: shifts.LocationInfo{
				Location:          "Harborview K-8",
				CityState:         "Seattle, WA",
				AssignmentSummary: "Provide 1:1 instructional support for students on IEPs.",
				ClassroomDetails:  []string{"K-8", "Resource Room"},
				Latitude:          47.6062,
				Longitude:         -122.3321,
			},
			withBreak: true,
		},
		{
			title:       "PE Teacher — Middle School",
			description: "Cover physical education classes, all grade levels.",
			hourlyRate:  40.0,
			shiftRate:   20.0,
			startDate:   date(2026, 6, 29),
			days:        9,
			startClock:  "08:30:00",
			endClock:    "16:00:00",
			partial:     false,
			urgent:      false,
			location: shifts.LocationInfo{
				Location:          "Eastgate Middle School",
				CityState:         "Bellevue, WA",
				AssignmentSummary: "PE classes for grades 6-8, gymnasium and outdoor field.",
				ClassroomDetails:  []string{"Grade 6-8", "Gymnasium"},
				Latitude:          47.6101,
				Longitude:         -122.2015,
			},
			withBreak: true,
		},
	}

	out := make([]Opportunity, 0, len(specs))
	for i, s := range specs {
		id := strconv.Itoa(i + 1)

		oppShifts := make([]shifts.Shift, 0, s.days)
		for d := 0; d < s.days; d++ {
			day := s.startDate.AddDate(0, 0, d)
			oppShifts = append(oppShifts, shifts.Shift{
				ID:           fmt.Sprintf("%s-%d", id, d+1),
				State:        "opportunity",
				TimeInfo:     buildTimeInfo(day, s.startClock, s.endClock),
				LocationInfo: s.location,
				Log:          buildLog(day, s.startClock, s.endClock, s.withBreak),
				HourlyRate:   s.shiftRate,
			})
		}

		last := s.startDate.AddDate(0, 0, s.days-1)
		out = append(out, Opportunity{
			ID:                   id,
			Title:                s.title,
			Description:          s.description,
			HourlyRate:           s.hourlyRate,
			StartAt:              isoAt(s.startDate, s.startClock),
			EndAt:                isoAt(last, s.endClock),
			CanTakePartialShifts: s.partial,
			ShowUrgencyNote:      s.urgent,
			LocationInfo:         s.location,
			Shifts:               oppShifts,
		})
	}

	return out
}

// --- helpers ----------------------------------------------------------------

func date(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func buildTimeInfo(day time.Time, startClock, endClock string) shifts.TimeInfo {
	return shifts.TimeInfo{
		Date:      day.Format("Monday, January 2"),
		TimeText:  fmt.Sprintf("%s - %s", humanClock(startClock), humanClock(endClock)),
		StartTime: isoAt(day, startClock),
		EndTime:   isoAt(day, endClock),
	}
}

func buildLog(day time.Time, startClock, endClock string, withBreak bool) shifts.ShiftLog {
	log := shifts.ShiftLog{
		CheckInTime:  isoAt(day, startClock),
		CheckOutTime: isoAt(day, endClock),
	}
	if withBreak {
		log.BreakStartTime = isoAt(day, "12:00:00")
		log.BreakEndTime = isoAt(day, "12:30:00")
	}
	return log
}

func isoAt(day time.Time, clock string) string {
	t, err := time.Parse("15:04:05", clock)
	if err != nil {
		return day.Format("2006-01-02") + "T" + clock + "Z"
	}
	combined := time.Date(day.Year(), day.Month(), day.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
	return combined.Format("2006-01-02T15:04:05Z")
}

func humanClock(clock string) string {
	t, err := time.Parse("15:04:05", clock)
	if err != nil {
		return clock
	}
	return t.Format("3:04 PM")
}

