package kafka

import (
	"context"
	"fmt"
	"log"
)

type EventConsumer struct {
	kafkaClient KafkaClient
}

func NewEventConsumer(client KafkaClient) *EventConsumer {
	return &EventConsumer{
		kafkaClient: client,
	}
}

func (ec *EventConsumer) StartConsuming(ctx context.Context) {
	ec.kafkaClient.Subscribe(ctx, ec.processMessage)
	log.Println("Event consumer started")
}

func (ec *EventConsumer) processMessage(ctx context.Context, msgData []byte) error {
	event, err := DeserializeEvent(msgData)
	if err != nil {
		return fmt.Errorf("failed to deserialize event: %w", err)
	}

	log.Printf("Processing event: %s, id: %s", event.Type, event.ID)

	switch event.Type {
	case UserCreated:
		return ec.handleUserCreatedEvent(ctx, event)
	default:
		log.Printf("Unknown event type: %s", event.Type)
		return nil
	}
}

func (ec *EventConsumer) handleUserCreatedEvent(ctx context.Context, event *Event) error {
	userEvent, err := DeserializeUserPayload(event)
	if err != nil {
		return fmt.Errorf("failed to deserialize user created payload: %w", err)
	}

	log.Printf("User created with username: %s", userEvent.User.Username)
	return nil
}
