package contracts

import "github.com/seanknox/myevent/lib/persistence"

// LocationCreatedEvent is emittd whenever a location is created
type LocationCreatedEvent struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Address string             `json:"address"`
	Country string             `json:"coutry"`
	Halls   []persistence.Hall `json:"halls"`
}

// EventName returns the events name
func (e *LocationCreatedEvent) EventName() string {
	return "locationCreated"
}
