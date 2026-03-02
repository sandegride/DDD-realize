package cmd

import (
	kafkain "basket-service/internal/adapters/in/kafka"
	grpcout "basket-service/internal/adapters/out/grpc/discount"
	kafkaout "basket-service/internal/adapters/out/kafka"
	"basket-service/internal/adapters/out/postgres"
	"basket-service/internal/adapters/out/postgres/outboxrepo"
	"basket-service/internal/core/application/eventhandlers"
	"basket-service/internal/core/application/usecases/commands"
	"basket-service/internal/core/application/usecases/queries"
	"basket-service/internal/core/domain/model/basket"
	"basket-service/internal/core/domain/services"
	"basket-service/internal/core/ports"
	"basket-service/internal/jobs"
	"basket-service/internal/pkg/ddd"
	"basket-service/internal/pkg/outbox"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
	"reflect"
	"sync"
)

type CompositionRoot struct {
	configs        Config
	gormDb         *gorm.DB
	discountClient ports.DiscountClient

	closers      []Closer
	onceDiscount sync.Once
}

func NewCompositionRoot(configs Config, gormDb *gorm.DB) *CompositionRoot {
	return &CompositionRoot{
		configs: configs,
		gormDb:  gormDb,
	}
}

func (cr *CompositionRoot) NewPromoGoodService() services.PromoGoodService {
	promoGoodService := services.NewPromoGoodService()
	return promoGoodService
}

func (cr *CompositionRoot) NewUnitOfWork() ports.UnitOfWork {
	unitOfWork, err := postgres.NewUnitOfWork(cr.gormDb)
	if err != nil {
		log.Fatalf("cannot create UnitOfWork: %v", err)
	}
	return unitOfWork
}

func (cr *CompositionRoot) NewAddAddressCommandHandler() commands.AddAddressCommandHandler {
	commandHandler, err := commands.NewAddAddressCommandHandler(cr.NewUnitOfWork())
	if err != nil {
		log.Fatalf("cannot create AddAddressCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewAddDeliveryPeriodCommandHandler() commands.AddDeliveryPeriodCommandHandler {
	commandHandler, err := commands.NewAddDeliveryPeriodCommandHandler(cr.NewUnitOfWork())
	if err != nil {
		log.Fatalf("cannot create AddDeliveryPeriodCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewChangeItemsCommandHandler() commands.ChangeItemsCommandHandler {
	commandHandler, err := commands.NewChangeItemsCommandHandler(
		cr.NewUnitOfWork())
	if err != nil {
		log.Fatalf("cannot create ChangeItemsCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewChangeStocksCommandHandler() commands.ChangeStocksCommandHandler {
	commandHandler, err := commands.NewChangeStocksCommandHandler(
		cr.NewUnitOfWork())
	if err != nil {
		log.Fatalf("cannot create ChangeStocksCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewCheckoutCommandHandler() commands.CheckoutCommandHandler {
	commandHandler, err := commands.NewCheckoutCommandHandler(
		cr.NewUnitOfWork(), cr.NewPromoGoodService(), cr.NewDiscountClient())
	if err != nil {
		log.Fatalf("cannot create CheckoutCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewGetBasketQueryHandler() queries.GetBasketQueryHandler {
	queryHandler, err := queries.NewGetBasketQueryHandler(cr.gormDb)
	if err != nil {
		log.Fatalf("cannot create GetBasketQueryHandler: %v", err)
	}
	return queryHandler
}

func (cr *CompositionRoot) NewDiscountClient() ports.DiscountClient {
	cr.onceDiscount.Do(func() {
		client, err := grpcout.NewClient(cr.configs.DiscountServiceGrpcHost)
		if err != nil {
			log.Fatalf("cannot create DiscountClient: %v", err)
		}
		cr.RegisterCloser(client)
		cr.discountClient = client
	})
	return cr.discountClient
}

func (cr *CompositionRoot) NewStocksChangedConsumer() kafkain.StocksChangedConsumer {
	consumer, err := kafkain.NewStocksChangedConsumer(
		[]string{cr.configs.KafkaHost},
		cr.configs.KafkaConsumerGroup,
		cr.configs.KafkaStocksChangedTopic,
		cr.NewChangeStocksCommandHandler(),
	)
	if err != nil {
		log.Fatalf("cannot create StocksChangedConsumer: %v", err)
	}
	cr.RegisterCloser(consumer)
	return consumer
}

func (cr *CompositionRoot) NewBasketCompletedDomainEventHandler() ddd.EventHandler {
	producer := cr.NewBasketProducer()
	handler, err := eventhandlers.NewBasketCompletedDomainEventHandler(producer)
	if err != nil {
		log.Fatalf("cannot create OrderCompletedDomainEventHandler: %v", err)
	}
	return handler
}

func (cr *CompositionRoot) NewMediatrWithSubscriptions() ddd.Mediatr {
	mediatr := ddd.NewMediatr()
	mediatr.Subscribe(cr.NewBasketCompletedDomainEventHandler(), basket.NewEmptyConfirmedDomainEvent())
	return mediatr
}

func (cr *CompositionRoot) NewBasketProducer() ports.BasketProducer {
	producer, err := kafkaout.NewBasketProducer([]string{cr.configs.KafkaHost}, cr.configs.KafkaBasketConfirmedTopic)
	if err != nil {
		log.Fatalf("cannot create OrderProducer: %v", err)
	}
	cr.RegisterCloser(producer)
	return producer
}

func (cr *CompositionRoot) NewEventRegistry() outbox.EventRegistry {
	registry, err := outbox.NewEventRegistry()
	if err != nil {
		log.Fatalf("cannot create EventRegistry: %v", err)
	}
	err = registry.RegisterDomainEvent(reflect.TypeOf(basket.ConfirmedDomainEvent{}))
	if err != nil {
		log.Fatalf("cannot register domain event: %v", err)
	}
	return registry
}

func (cr *CompositionRoot) NewOutboxRepository() outboxrepo.OutboxRepository {
	repository, err := outboxrepo.NewRepository(cr.gormDb)
	if err != nil {
		log.Fatalf("cannot create OutboxRepository: %v", err)
	}
	return repository
}

func (cr *CompositionRoot) NewOutboxJob() cron.Job {
	job, err := jobs.NewOutboxJob(cr.NewOutboxRepository(), cr.NewEventRegistry(), cr.NewMediatrWithSubscriptions())
	if err != nil {
		log.Fatalf("cannot create OutboxJob: %v", err)
	}
	return job
}
