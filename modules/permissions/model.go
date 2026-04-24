package permissions

// Permissions is the singleton document that tracks which device-level
// permissions the user has granted to the app.
type Permissions struct {
	PushNotificationsEnabled bool `json:"push_notifications_enabled"`
	LocationAccessEnabled    bool `json:"location_access_enabled"`
}

