package kafka

import (
	"encoding/json"
	"fmt"
	"mindgames/internal/domain"
)

const (
	UserCreated EventType = "user.created"
)

type UserEvent struct {
	User domain.User `json:"user"`
}

func DeserializeUserPayload(event *Event) (*UserEvent, error) {
	payloadJSON, err := json.Marshal(event.Payload)
	if err != nil {
		return nil, fmt.Errorf("error serializing payload: %w", err)
	}

	var userPayload UserEvent
	if err := json.Unmarshal(payloadJSON, &userPayload); err != nil {
		return nil, fmt.Errorf("error deserializing user payload: %w", err)
	}

	return &userPayload, nil
}
