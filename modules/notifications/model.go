package notifications

// Notification represents an in-app notification for a user.
type Notification struct {
	UUID       string         `json:"uuid"`
	UserUUID   string         `json:"user_uuid"`
	Type       string         `json:"type"`
	Title      string         `json:"title"`
	Body       string         `json:"body"`
	CreatedAt  string         `json:"created_at"`
	UpdatedAt  string         `json:"updated_at"`
	IsRead     bool           `json:"is_read"`
	ActionText string         `json:"action_text"`
	Data       map[string]any `json:"data"`
}

// GetID implements data.Entity.
func (n *Notification) GetID() string { return n.UUID }

// SetID implements data.Entity.
func (n *Notification) SetID(id string) { n.UUID = id }

