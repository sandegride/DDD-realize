package kafka

import (
	"basket-service/internal/core/application/usecases/commands"
	"basket-service/internal/generated/queues/stockschangedpb"
	"basket-service/internal/pkg/errs"
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"log"
)

type StocksChangedConsumer interface {
	Consume() error
	Close() error
}

var _ StocksChangedConsumer = &stocksChangedConsumer{}
var _ sarama.ConsumerGroupHandler = &stocksChangedConsumer{}

type stocksChangedConsumer struct {
	topic                      string
	consumerGroup              sarama.ConsumerGroup
	changeStocksCommandHandler commands.ChangeStocksCommandHandler
	ctx                        context.Context
	cancel                     context.CancelFunc
}

func NewStocksChangedConsumer(
	brokers []string,
	group string,
	topic string,
	changeStocksCommandHandler commands.ChangeStocksCommandHandler,
) (StocksChangedConsumer, error) {
	if len(brokers) == 0 {
		return nil, errs.NewValueIsRequiredError("brokers")
	}
	if group == "" {
		return nil, errs.NewValueIsRequiredError("group")
	}
	if topic == "" {
		return nil, errs.NewValueIsRequiredError("topic")
	}
	if changeStocksCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("changeStocksCommandHandler")
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.Version = sarama.V3_4_0_0
	saramaCfg.Consumer.Return.Errors = true
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &stocksChangedConsumer{
		topic:                      topic,
		consumerGroup:              consumerGroup,
		changeStocksCommandHandler: changeStocksCommandHandler,
		ctx:                        ctx,
		cancel:                     cancel,
	}, nil
}

func (c *stocksChangedConsumer) Close() error {
	c.cancel()
	return c.consumerGroup.Close()
}

func (c *stocksChangedConsumer) Consume() error {
	for {
		err := c.consumerGroup.Consume(c.ctx, []string{c.topic}, c)
		if err != nil {
			log.Printf("Error from consumer: %v", err)
			return err
		}
		if c.ctx.Err() != nil {
			return nil
		}
	}
}

// Реализация sarama.ConsumerGroupHandler:

func (c *stocksChangedConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *stocksChangedConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *stocksChangedConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		ctx := context.Background()
		fmt.Printf("Received: topic = %s, partition = %d, offset = %d, key = %s, value = %s\n",
			message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))

		var event stockschangedpb.StocksChangedIntegrationEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		cmd, err := commands.NewChangeStocksCommand(
			uuid.MustParse(event.GoodId), int(event.Quantity),
		)
		if err != nil {
			log.Printf("Failed to create changeStocks command: %v", err)
			session.MarkMessage(message, "")
			continue
		}

		if err := c.changeStocksCommandHandler.Handle(ctx, cmd); err != nil {
			log.Printf("Failed to handle changeStocks command: %v", err)
		}

		session.MarkMessage(message, "")
	}
	return nil
}
