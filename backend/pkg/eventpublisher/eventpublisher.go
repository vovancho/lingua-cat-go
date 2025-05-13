package eventpublisher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-sql/v3/pkg/sql"
	"github.com/ThreeDotsLabs/watermill/message"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type PublisherInterface interface {
	Publish(ctx context.Context, tx sql.ContextExecutor, event any) error
}

type EventTopicName string

type EventPublisher struct {
	topic EventTopicName
}

func NewEventPublisher(topic EventTopicName) *EventPublisher {
	return &EventPublisher{
		topic: topic,
	}
}

func (p *EventPublisher) Publish(ctx context.Context, tx sql.ContextExecutor, event any) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal Event: %w", err)
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)

	// Трейсинг
	propagator := otel.GetTextMapPropagator()
	carrier := propagation.MapCarrier{}
	propagator.Inject(ctx, carrier)
	for k, v := range carrier {
		msg.Metadata[k] = v
	}

	// Publisher на транзакции
	publisherConfig := sql.PublisherConfig{
		SchemaAdapter: sql.DefaultPostgreSQLSchema{},
	}

	txPublisher, err := sql.NewPublisher(tx, publisherConfig, nil)
	if err != nil {
		return fmt.Errorf("create txPublisher: %w", err)
	}

	if err := txPublisher.Publish(string(p.topic), msg); err != nil {
		return fmt.Errorf("publish event: %w", err)
	}

	return nil
}
