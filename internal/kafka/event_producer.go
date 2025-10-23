package kafka

import (
	"context"
	"fmt"
	"log"
	"mindgames/internal/domain"
)

type EventProducer interface {
	PublishCreateUser(ctx context.Context, user *domain.User) error
}

type EventProducerImpl struct {
	client IKafkaClient
}

func NewEventProducer(client IKafkaClient) *EventProducerImpl {
	return &EventProducerImpl{
		client: client,
	}
}

func (p *EventProducerImpl) PublishCreateUser(ctx context.Context, user *domain.User) error {
	payload := UserEvent{
		User: *user,
	}

	event := NewEvent(UserCreated, payload)

	eventData, err := event.Serialize()
	if err != nil {
		return fmt.Errorf("event serialization error: %w", err)
	}

	err = p.client.Publish(ctx, string(event.Type), eventData)
	if err != nil {
		return fmt.Errorf("event publishing error: %w", err)
	}

	log.Printf("Event published %s ID %s", event.Type, event.ID)
	return nil
}
