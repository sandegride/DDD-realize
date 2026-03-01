package cmd

import (
	kafkain "basket-service/internal/adapters/in/kafka"
	grpcout "basket-service/internal/adapters/out/grpc/discount"
	"basket-service/internal/adapters/out/postgres"
	"basket-service/internal/core/application/usecases/commands"
	"basket-service/internal/core/application/usecases/queries"
	"basket-service/internal/core/domain/services"
	"basket-service/internal/core/ports"
	"gorm.io/gorm"
	"log"
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

func (cr *CompositionRoot) NewUnitOfWork() ports.UnitOfWork {
	unitOfWork, err := postgres.NewUnitOfWork(cr.gormDb)
	if err != nil {
		log.Fatalf("cannot create UnitOfWork: %v", err)
	}
	return unitOfWork
}

func (cr *CompositionRoot) NewUnitOfWorkFactory() ports.UnitOfWorkFactory {
	unitOfWorkFactory, err := postgres.NewUnitOfWorkFactory(cr.gormDb)
	if err != nil {
		log.Fatalf("cannot create UnitOfWorkFactory: %v", err)
	}
	return unitOfWorkFactory
}

func (cr *CompositionRoot) NewAddAddressCommandHandler() commands.AddAddressCommandHandler {
	commandHandler, err := commands.NewAddAddressCommandHandler(cr.NewUnitOfWorkFactory())
	if err != nil {
		log.Fatalf("cannot create AddAddressCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewAddDeliveryPeriodCommandHandler() commands.AddDeliveryPeriodCommandHandler {
	commandHandler, err := commands.NewAddDeliveryPeriodCommandHandler(cr.NewUnitOfWorkFactory())
	if err != nil {
		log.Fatalf("cannot create AddDeliveryPeriodCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewChangeItemsCommandHandler() commands.ChangeItemsCommandHandler {
	commandHandler, err := commands.NewChangeItemsCommandHandler(
		cr.NewUnitOfWorkFactory())
	if err != nil {
		log.Fatalf("cannot create ChangeItemsCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewChangeStocksCommandHandler() commands.ChangeStocksCommandHandler {
	commandHandler, err := commands.NewChangeStocksCommandHandler(
		cr.NewUnitOfWorkFactory())
	if err != nil {
		log.Fatalf("cannot create ChangeStocksCommandHandler: %v", err)
	}
	return commandHandler
}

func (cr *CompositionRoot) NewCheckoutCommandHandler() commands.CheckoutCommandHandler {
	commandHandler, err := commands.NewCheckoutCommandHandler(
		cr.NewUnitOfWorkFactory(), cr.NewPromoGoodService(), cr.NewDiscountClient())
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
