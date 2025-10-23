package kafka

import "mindgames/internal/domain"

const (
	UserCreated EventType = "user.created"
)

type UserEvent struct {
	User domain.User `json:"user"`
}
