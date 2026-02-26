package postgres

import (
	"basket-service/internal/adapters/out/postgres/basketrepo"
	"basket-service/internal/adapters/out/postgres/goodrepo"
	"basket-service/internal/core/domain/model/basket"
	"basket-service/internal/core/domain/model/good"
	"basket-service/internal/pkg/testcnts"
	"context"
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
	err = db.AutoMigrate(&goodrepo.GoodDTO{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&basketrepo.BasketDTO{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&basketrepo.ItemDTO{})
	assert.NoError(t, err)
	err = db.AutoMigrate(&basketrepo.DeliveryPeriodDTO{})
	assert.NoError(t, err)

	// Очистка выполняется после завершения теста
	t.Cleanup(func() {
		err := postgresContainer.Terminate(ctx)
		assert.NoError(t, err)
	})

	return ctx, db, nil
}

func Test_GoodRepositoryShouldCanAddGood(t *testing.T) {
	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	// Создаем UnitOfWork
	uow, err := NewUnitOfWork(db)
	assert.NoError(t, err)

	// Вызываем Add
	goodAggregate := good.Coffee()
	err = uow.GoodRepository().Add(ctx, goodAggregate)
	assert.NoError(t, err)

	// Считываем данные из БД
	var goodFromDb goodrepo.GoodDTO
	err = db.First(&goodFromDb, "id = ?", goodAggregate.ID()).Error
	assert.NoError(t, err)

	// Проверяем эквивалентность
	assert.Equal(t, goodAggregate.ID(), goodFromDb.ID)
	assert.Equal(t, goodAggregate.Title(), goodFromDb.Title)
	assert.Equal(t, goodAggregate.Description(), goodFromDb.Description)
	assert.Equal(t, goodAggregate.Price(), goodFromDb.Price)
	assert.Equal(t, goodAggregate.Quantity(), goodFromDb.Quantity)
	assert.Equal(t, goodAggregate.Weight().Value(), goodFromDb.Weight.Value)
}

func Test_BasketRepositoryShouldCanAddBasket(t *testing.T) {
	// Инициализируем окружение
	ctx, db, err := setupTest(t)
	assert.NoError(t, err)

	// Создаем UnitOfWork
	uow, err := NewUnitOfWork(db)
	assert.NoError(t, err)

	// Вызываем Add

	basketAggregate, err := basket.NewBasket(uuid.New())
	err = uow.BasketRepository().Add(ctx, basketAggregate)
	assert.NoError(t, err)

	// Считываем данные из БД
	var basketFromDb basketrepo.BasketDTO
	err = db.First(&basketFromDb, "id = ?", basketAggregate.ID()).Error
	assert.NoError(t, err)

	// Проверяем эквивалентность
	assert.Equal(t, basketAggregate.ID(), basketFromDb.ID)
	assert.Equal(t, basketAggregate.Status(), basketFromDb.Status)
}
