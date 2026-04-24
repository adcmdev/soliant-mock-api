package permissions

// Seed returns the default permissions document used to populate an empty
// store.
func Seed() Permissions {
	return Permissions{
		PushNotificationsEnabled: true,
		LocationAccessEnabled:    true,
	}
}

