package postgres

import (
	"context"
	"delivery-service/internal/adapters/out/postgres/courierrepo"
	"delivery-service/internal/adapters/out/postgres/orderrepo"
	"delivery-service/internal/core/domain/model/courier"
	"delivery-service/internal/core/domain/model/kernel"
	"delivery-service/internal/core/domain/model/order"
	"delivery-service/internal/pkg/testcnts"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	postgresgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func setupTest(t *testing.T) (context.Context, *gorm.DB, error) {
	ctx := context.Background()
	postgresContainer, dsn, err := testcnts.StartPostgresContainer(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Подключаемся к БД через Gorm
	db, err := gorm.Open(postgresgorm.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// Авто миграция (создаём таблицу)
	err = db.AutoMigrate(&courierrepo.CourierDTO{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&courierrepo.StoragePlaceDTO{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&orderrepo.OrderDTO{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&courierrepo.StoragePlaceDTO{})
	assert.NoError(t, err)

	// Очистка выполняется после завершения теста
	t.Cleanup(func() {
		err := postgresContainer.Terminate(ctx)
		assert.NoError(t, err)
	})

	return ctx, db, nil
}

func Test_CourierRepositoryShouldCanAddCourier(t *testing.T) {
	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	// Создаем UnitOfWorkV2
	uow, err := NewUnitOfWork(db)
	assert.NoError(t, err)

	// Вызываем Add
	location := kernel.MaxLocation()
	courierAggregate, err := courier.NewCourier("Велосипедист", 2, location)
	err = uow.CourierRepository().Add(ctx, courierAggregate)
	assert.NoError(t, err)

	// Считываем данные из БД
	var courierFromDb courierrepo.CourierDTO
	err = db.First(&courierFromDb, "id = ?", courierAggregate.ID()).Error
	assert.NoError(t, err)

	// Проверяем эквивалентность
	assert.Equal(t, courierAggregate.ID(), courierFromDb.ID)
	assert.Equal(t, courierAggregate.Speed(), courierFromDb.Speed)
}

func Test_OrderRepositoryShouldCanAddOrder(t *testing.T) {
	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	// Создаем UnitOfWorkV2
	uow, err := NewUnitOfWork(db)
	assert.NoError(t, err)

	// Вызываем Add
	location := kernel.MinLocation()
	orderAggregate, err := order.NewOrder(uuid.New(), location, 10)
	err = uow.OrderRepository().Add(ctx, orderAggregate)
	assert.NoError(t, err)

	// Считываем данные из БД
	var orderFromDb orderrepo.OrderDTO
	err = db.First(&orderFromDb, "id = ?", orderAggregate.ID()).Error
	assert.NoError(t, err)

	// Проверяем эквивалентность
	assert.Equal(t, orderAggregate.ID(), orderFromDb.ID)
	assert.Equal(t, orderAggregate.Location().X(), orderFromDb.Location.X)
	assert.Equal(t, orderAggregate.Location().Y(), orderFromDb.Location.Y)
	assert.Equal(t, orderAggregate.Volume(), orderFromDb.Volume)
	assert.Equal(t, orderAggregate.Status(), orderFromDb.Status)
}
