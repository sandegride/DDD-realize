package main

import (
	"basket-service/cmd"
	httpin "basket-service/internal/adapters/in/http"
	"basket-service/internal/adapters/out/postgres/basketrepo"
	"basket-service/internal/adapters/out/postgres/goodrepo"
	"basket-service/internal/core/domain/model/good"
	"basket-service/internal/generated/servers"
	"basket-service/internal/pkg/errs"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	oam "github.com/oapi-codegen/echo-middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"os"
)

func main() {
	configs := getConfigs()

	connectionString, err := makeConnectionString(
		configs.DbHost,
		configs.DbPort,
		configs.DbUser,
		configs.DbPassword,
		configs.DbName,
		configs.DbSslMode)
	if err != nil {
		log.Fatal(err.Error())
	}

	crateDbIfNotExists(configs.DbHost,
		configs.DbPort,
		configs.DbUser,
		configs.DbPassword,
		configs.DbName,
		configs.DbSslMode)
	gormDb := mustGormOpen(connectionString)
	mustAutoMigrate(gormDb)

	compositionRoot := cmd.NewCompositionRoot(
		configs,
		gormDb,
	)
	defer compositionRoot.CloseAll()

	startKafkaConsumer(compositionRoot)
	startWebServer(compositionRoot, configs.HttpPort)
}

func getConfigs() cmd.Config {
	config := cmd.Config{
		HttpPort:                  goDotEnvVariable("HTTP_PORT"),
		DbHost:                    goDotEnvVariable("DB_HOST"),
		DbPort:                    goDotEnvVariable("DB_PORT"),
		DbUser:                    goDotEnvVariable("DB_USER"),
		DbPassword:                goDotEnvVariable("DB_PASSWORD"),
		DbName:                    goDotEnvVariable("DB_NAME"),
		DbSslMode:                 goDotEnvVariable("DB_SSLMODE"),
		DiscountServiceGrpcHost:   goDotEnvVariable("DISCOUNT_SERVICE_GRPC_HOST"),
		KafkaHost:                 goDotEnvVariable("KAFKA_HOST"),
		KafkaConsumerGroup:        goDotEnvVariable("KAFKA_CONSUMER_GROUP"),
		KafkaStocksChangedTopic:   goDotEnvVariable("KAFKA_STOCKS_CHANGED_TOPIC"),
		KafkaBasketConfirmedTopic: goDotEnvVariable("KAFKA_BASKET_CONFIRMED_TOPIC"),
	}
	return config
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func crateDbIfNotExists(host string, port string, user string,
	password string, dbName string, sslMode string) {
	dsn, err := makeConnectionString(host, port, user, password, "postgres", sslMode)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к PostgreSQL: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("ошибка при закрытии db: %v", err)
		}
	}()

	// Создаём базу данных, если её нет
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		log.Printf("Ошибка создания БД (возможно, уже существует): %v", err)
	}
}

func makeConnectionString(host string, port string, user string,
	password string, dbName string, sslMode string) (string, error) {
	if host == "" {
		return "", errs.NewValueIsRequiredError(host)
	}
	if port == "" {
		return "", errs.NewValueIsRequiredError(port)
	}
	if user == "" {
		return "", errs.NewValueIsRequiredError(user)
	}
	if password == "" {
		return "", errs.NewValueIsRequiredError(password)
	}
	if dbName == "" {
		return "", errs.NewValueIsRequiredError(dbName)
	}
	if sslMode == "" {
		return "", errs.NewValueIsRequiredError(sslMode)
	}
	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		host,
		port,
		user,
		password,
		dbName,
		sslMode), nil
}

func mustGormOpen(connectionString string) *gorm.DB {
	pgGorm, err := gorm.Open(postgres.New(
		postgres.Config{
			DSN:                  connectionString,
			PreferSimpleProtocol: true,
		},
	), &gorm.Config{})
	if err != nil {
		log.Fatalf("connection to postgres through gorm\n: %s", err)
	}
	return pgGorm
}

func mustAutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&basketrepo.BasketDTO{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	err = db.AutoMigrate(&basketrepo.ItemDTO{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	err = db.AutoMigrate(&basketrepo.DeliveryPeriodDTO{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	err = db.AutoMigrate(&goodrepo.GoodDTO{})
	if err != nil {
		log.Fatalf("Ошибка миграции: %v", err)
	}

	// Seed
	goodDTOs := make([]goodrepo.GoodDTO, len(good.Goods))
	for i, goodAggregate := range good.Goods {
		goodDTOs[i] = goodrepo.DomainToDTO(goodAggregate)
	}

	err = db.Clauses(clause.OnConflict{
		UpdateAll: true, // обновит все поля, кроме PK
	}).Create(&goodDTOs).Error
	if err != nil {
		log.Fatalf("Ошибка при вставке товаров: %v", err)
	}
}

func startWebServer(compositionRoot *cmd.CompositionRoot, port string) {
	handlers, err := httpin.NewServer(
		compositionRoot.NewAddAddressCommandHandler(),
		compositionRoot.NewAddDeliveryPeriodCommandHandler(),
		compositionRoot.NewChangeItemsCommandHandler(),
		compositionRoot.NewCheckoutCommandHandler(),
		compositionRoot.NewGetBasketQueryHandler(),
	)
	if err != nil {
		log.Fatalf("Ошибка инициализации HTTP Server: %v", err)
	}

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
	}))

	spec, err := servers.GetSwagger()
	if err != nil {
		log.Fatalf("Error reading OpenAPI spec: %v", err)
	}
	e.Use(oam.OapiRequestValidator(spec))
	e.Pre(middleware.RemoveTrailingSlash())
	registerSwaggerOpenApi(e)
	registerSwaggerUi(e)
	servers.RegisterHandlers(e, handlers)
	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%s", port)))
}

func registerSwaggerOpenApi(e *echo.Echo) {
	e.GET("/openapi.json", func(c echo.Context) error {
		swagger, err := servers.GetSwagger()
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to load swagger: "+err.Error())
		}

		data, err := swagger.MarshalJSON()
		if err != nil {
			return c.String(http.StatusInternalServerError, "failed to marshal swagger: "+err.Error())
		}

		return c.Blob(http.StatusOK, "application/json", data)
	})
}

func registerSwaggerUi(e *echo.Echo) {
	e.GET("/docs", func(c echo.Context) error {
		html := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		  <meta charset="UTF-8">
		  <title>Swagger UI</title>
		  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css">
		</head>
		<body>
		  <div id="swagger-ui"></div>
		  <script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
		  <script>
			window.onload = () => {
			  SwaggerUIBundle({
				url: "/openapi.json",
				dom_id: "#swagger-ui",
			  });
			};
		  </script>
		</body>
		</html>`
		return c.HTML(http.StatusOK, html)
	})
}

func startKafkaConsumer(compositionRoot *cmd.CompositionRoot) {
	go func() {
		if err := compositionRoot.NewStocksChangedConsumer().Consume(); err != nil {
			log.Fatalf("Kafka consumer error: %v", err)
		}
	}()
}
