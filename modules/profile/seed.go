package profile

// Seed returns the default profile document used to populate an empty store.
func Seed() Profile {
	return Profile{
		User: User{
			UUID:        "mock-user-uuid",
			FirstName:   "Andres",
			LastName:    "Carrillo",
			Email:       "user@example.com",
			PhoneNumber: "(123)-456-7890",
			CreatedAt:   "2026-03-25T00:00:00Z",
			UpdatedAt:   "2026-04-24T00:00:00Z",
		},
		HoursWorked:     672,
		ShiftsCompleted: 84,
	}
}

