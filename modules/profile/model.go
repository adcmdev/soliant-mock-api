package profile

// Profile is the singleton document returned by GET /profile. It bundles the
// authenticated user's personal info with aggregate stats.
type Profile struct {
	User            User `json:"user"`
	HoursWorked     int  `json:"hours_worked"`
	ShiftsCompleted int  `json:"shifts_completed"`
}

// User holds the personal details of the authenticated user.
type User struct {
	UUID        string `json:"uuid"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

