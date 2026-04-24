package location

// Location is the singleton document that holds the user's saved address and
// preferred travel radius.
type Location struct {
	AddressLine1 string `json:"address_line_1"`
	AddressLine2 string `json:"address_line_2"`
	TravelRadius int    `json:"travel_radius"`
}

