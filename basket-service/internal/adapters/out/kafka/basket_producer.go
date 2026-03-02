package kafka

import (
	"basket-service/internal/core/domain/model/basket"
	"basket-service/internal/core/ports"
	"basket-service/internal/generated/queues/basketconfirmedpb"
	"basket-service/internal/pkg/ddd"
	"basket-service/internal/pkg/errs"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

type basketProducer struct {
	topic    string
	producer sarama.SyncProducer
}

func NewBasketProducer(brokers []string, topic string) (ports.BasketProducer, error) {
	if len(brokers) == 0 {
		return nil, errs.NewValueIsRequiredError("brokers")
	}
	if topic == "" {
		return nil, errs.NewValueIsRequiredError("topic")
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V3_4_0_0
	saramaCfg.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("create sync producer: %w", err)
	}

	return &basketProducer{
		topic:    topic,
		producer: producer,
	}, nil
}

func (p *basketProducer) Close() error {
	return p.producer.Close()
}

func (p *basketProducer) Publish(ctx context.Context, domainEvent ddd.DomainEvent) error {
	confirmedDomainEvent, ok := domainEvent.(*basket.ConfirmedDomainEvent)
	if !ok {
		return fmt.Errorf("unexpected domain event type: %T", domainEvent)
	}

	integrationEvent := p.mapDomainEventToIntegrationEvent(confirmedDomainEvent)

	bytes, err := json.Marshal(integrationEvent)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(confirmedDomainEvent.Payload.ID.String()),
		Value: sarama.ByteEncoder(bytes),
	}

	resultCh := make(chan error, 1)

	go func() {
		_, _, err := p.producer.SendMessage(msg)
		resultCh <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-resultCh:
		return err
	}
}

func (p *basketProducer) mapDomainEventToIntegrationEvent(domainEvent *basket.ConfirmedDomainEvent) *basketconfirmedpb.BasketConfirmedIntegrationEvent {
	items := make([]*basketconfirmedpb.Item, 0)
	for _, item := range domainEvent.Payload.Items {
		items = append(items, &basketconfirmedpb.Item{
			Id:       item.ID.String(),
			GoodId:   item.GoodID.String(),
			Title:    item.Title,
			Price:    item.Price,
			Quantity: int32(item.Quantity),
		})
	}
	integrationEvent := basketconfirmedpb.BasketConfirmedIntegrationEvent{
		BasketId: domainEvent.Payload.ID.String(),
		Address: &basketconfirmedpb.Address{
			Country:   domainEvent.Payload.Address.Country,
			City:      domainEvent.Payload.Address.City,
			Street:    domainEvent.Payload.Address.Street,
			House:     domainEvent.Payload.Address.House,
			Apartment: domainEvent.Payload.Address.Apartment,
		},
		Items: items,
		DeliveryPeriod: &basketconfirmedpb.DeliveryPeriod{
			From: int32(domainEvent.Payload.DeliveryPeriod.From),
			To:   int32(domainEvent.Payload.DeliveryPeriod.To),
		},
	}
	return &integrationEvent
}
