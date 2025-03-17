package models

// AlarmState represents the possible states of an alarm.
type AlarmState string

const (
	Triggered AlarmState = "Triggered"
	Active    AlarmState = "Active"
	ACKed     AlarmState = "ACKed"
	Cleared   AlarmState = "Cleared"
)

// Alarm represents the structure for an alarm with essential details.
type Alarm struct {
	ID        string     `json:"id"`        // Unique identifier for the alarm
	Name      string     `json:"name"`      // Descriptive name of the alarm
	State     AlarmState `json:"state"`     // Current state of the alarm
	CreatedAt string     `json:"created_at"` // Creation timestamp of the alarm
	UpdatedAt string     `json:"updated_at"` // Last updated timestamp of the alarm
	ACKedAt   string     `json:"acked_at"`   // Timestamp for when the alarm was acknowledged
}

// IsValid checks if the provided alarm state is valid.
func (a AlarmState) IsValid() bool {
	switch a {
	case Triggered, Active, ACKed, Cleared:
		return true
	default:
		return false
	}
}
