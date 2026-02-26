package cmd

type Config struct {
	HttpPort                  string
	DbHost                    string
	DbPort                    string
	DbUser                    string
	DbPassword                string
	DbName                    string
	DbSslMode                 string
	DiscountServiceGrpcHost   string
	KafkaHost                 string
	KafkaConsumerGroup        string
	KafkaStocksChangedTopic   string
	KafkaBasketConfirmedTopic string
}
