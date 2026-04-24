package auth

// Seed returns the default auth user used to respond to /auth/*/verify and
// /shifts/detail/verify. Timestamps are left nil so the JSON response emits
// `created_at: null` / `updated_at: null`, matching the spec for a user who
// has authenticated but not yet completed their profile.
func Seed() User {
	return User{
		UUID:        "mock-user-uuid",
		FirstName:   "",
		LastName:    "",
		Email:       "user@example.com",
		PhoneNumber: "",
		CreatedAt:   nil,
		UpdatedAt:   nil,
	}
}

