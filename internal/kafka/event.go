package kafka

import (
	"encoding/json"
	"fmt"
	"mindgames/internal/domain"
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	UserCreated EventType = "user.created"
)

type Event struct {
	ID        string      `json:"id"`
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

func NewEvent(eventType EventType, payload interface{}) Event {
	return Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Timestamp: time.Now(),
		Payload:   payload,
	}
}

type UserEvent struct {
	User domain.User `json:"user"`
}

type LoginEvent struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
}

func (e *Event) Serialize() ([]byte, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("\nerror serializing event: %w", err)
	}
	return data, nil
}
