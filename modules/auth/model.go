package auth

// User is the payload returned by the auth/verify endpoints. Timestamps are
// pointers so they can be emitted as JSON `null` when absent, matching the
// spec for newly authenticated users who haven't completed signup yet.
type User struct {
	UUID        string  `json:"uuid"`
	FirstName   string  `json:"first_name"`
	LastName    string  `json:"last_name"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phone_number"`
	CreatedAt   *string `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
}

