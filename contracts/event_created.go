package contracts

import "time"

// EventCreatedEvent is emitted whenever a new event is created
type EventCreatedEvent struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	LocationID string    `json:"location_id"`
	Start      time.Time `json:"start_time"`
	End        time.Time `json:"end_time"`
}

// EventName returns the event's name
func (e *EventCreatedEvent) EventName() string {
	return "eventCreated"
}
