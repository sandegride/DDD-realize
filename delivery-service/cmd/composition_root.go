package cmd

import (
	kafkain "delivery-service/internal/adapters/in/kafka"
	grpcout "delivery-service/internal/adapters/out/grpc/geo"
	kafkaout "delivery-service/internal/adapters/out/kafka"
	"delivery-service/internal/adapters/out/postgres"
	"delivery-service/internal/adapters/out/postgres/outboxrepo"
	"delivery-service/internal/core/application/eventhandlers"
	"delivery-service/internal/core/application/usecases/commands"
	"delivery-service/internal/core/application/usecases/queries"
	"delivery-service/internal/core/domain/model/order"
	"delivery-service/internal/core/domain/services"
	"delivery-service/internal/core/ports"
	"delivery-service/internal/jobs"
	"delivery-service/internal/pkg/ddd"
	"delivery-service/internal/pkg/outbox"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"reflect"
	"sync"

	_ "gorm.io/gorm"
	"log"
)

type CompositionRoot struct {
	configs   Config
	gormDb    *gorm.DB
	geoClient ports.GeoClient

	closers []Closer
	onceGeo sync.Once
}

func NewCompositionRoot(configs Config, gormDb *gorm.DB) *CompositionRoot {
	return &CompositionRoot{
		configs: configs,
		gormDb:  gormDb,
	}
}

func (cr *CompositionRoot) NewOrderDispatcher() services.OrderDispatcher {
	orderDispatcher := services.NewOrderDispatcher()
	return orderDispatcher
}

func (cr *CompositionRoot) NewUnitOfWork() ports.UnitOfWork {
	unitOfWork, err := postgres.NewUnitOfWork(cr.gormDb)
	if err != nil {
		log.Fatalf("cannot create UnitOfWork: %v", err)
	}
	return unitOfWork
}

func (cr *CompositionRoot) NewCreateOrderCommandHandler() commands.CreateOrderCommandHandler {
	commandHandler, err := commands.NewCreateOrderCommandHandler(cr.NewUnitOfWork(), cr.NewGeoClient())
	if err != nil {
		log.Fatalf("cannot create CreateOrderCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewCreateCourierCommandHandler() commands.CreateCourierCommandHandler {
	commandHandler, err := commands.NewCreateCourierCommandHandler(cr.NewUnitOfWork())
	if err != nil {
		log.Fatalf("cannot create CreateCourierCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewAssignOrdersCommandHandler() commands.AssignOrdersCommandHandler {
	commandHandler, err := commands.NewAssignOrdersCommandHandler(
		cr.NewUnitOfWork(), cr.NewOrderDispatcher())
	if err != nil {
		log.Fatalf("cannot create AssignOrdersCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewMoveCouriersCommandHandler() commands.MoveCouriersCommandHandler {
	commandHandler, err := commands.NewMoveCouriersCommandHandler(
		cr.NewUnitOfWork())
	if err != nil {
		log.Fatalf("cannot create MoveCouriersCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewGetAllCouriersQueryHandler() queries.GetAllCouriersQueryHandler {
	queryHandler, err := queries.NewGetAllCouriersQueryHandler(cr.gormDb)
	if err != nil {
		log.Fatalf("cannot create GetAllCouriersQueryHandler: %v", err)
	}
	return queryHandler
}

func (cr *CompositionRoot) NewGetNotCompletedOrdersQueryHandler() queries.GetNotCompletedOrdersQueryHandler {
	queryHandler, err := queries.NewGetNotCompletedOrdersQueryHandler(cr.gormDb)
	if err != nil {
		log.Fatalf("cannot create GetNotCompletedOrdersQueryHandler: %v", err)
	}
	return queryHandler
}

func (cr *CompositionRoot) NewAssignOrdersJob() cron.Job {
	job, err := jobs.NewAssignOrdersJob(cr.NewAssignOrdersCommandHandler())
	if err != nil {
		log.Fatalf("cannot create AssignOrdersJob: %v", err)
	}
	return job
}

func (cr *CompositionRoot) NewMoveCouriersJob() cron.Job {
	job, err := jobs.NewMoveCouriersJob(cr.NewMoveCouriersCommandHandler())
	if err != nil {
		log.Fatalf("cannot create MoveCouriersJob: %v", err)
	}
	return job
}

func (cr *CompositionRoot) NewGeoClient() ports.GeoClient {
	cr.onceGeo.Do(func() {
		client, err := grpcout.NewClient(cr.configs.GeoServiceGrpcHost)
		if err != nil {
			log.Fatalf("cannot create GeoClient: %v", err)
		}
		cr.RegisterCloser(client)
		cr.geoClient = client
	})
	return cr.geoClient
}

func (cr *CompositionRoot) NewBasketConfirmedConsumer() kafkain.BasketConfirmedConsumer {
	consumer, err := kafkain.NewBasketConfirmedConsumer(
		[]string{cr.configs.KafkaHost},
		cr.configs.KafkaConsumerGroup,
		cr.configs.KafkaBasketConfirmedTopic,
		cr.NewCreateOrderCommandHandler(),
	)
	if err != nil {
		log.Fatalf("cannot create BasketConfirmedConsumer: %v", err)
	}
	cr.RegisterCloser(consumer)
	return consumer
}

func (cr *CompositionRoot) NewOrderCompletedDomainEventHandler() ddd.EventHandler {
	producer := cr.NewOrderProducer()
	handler, err := eventhandlers.NewOrderCompletedDomainEventHandler(producer)
	if err != nil {
		log.Fatalf("cannot create OrderCompletedDomainEventHandler: %v", err)
	}
	return handler
}

func (cr *CompositionRoot) NewMediatrWithSubscriptions() ddd.Mediatr {
	mediatr := ddd.NewMediatr()
	mediatr.Subscribe(cr.NewOrderCompletedDomainEventHandler(), order.NewEmptyCompletedDomainEvent())
	return mediatr
}

func (cr *CompositionRoot) NewOrderProducer() ports.OrderProducer {
	producer, err := kafkaout.NewOrderProducer([]string{cr.configs.KafkaHost}, cr.configs.KafkaOrderChangedTopic)
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
	err = registry.RegisterDomainEvent(reflect.TypeOf(order.CompletedDomainEvent{}))
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
