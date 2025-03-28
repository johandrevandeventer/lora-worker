package engine

import (
	"log"

	"github.com/johandrevandeventer/kafkaclient/config"
	"github.com/johandrevandeventer/kafkaclient/consumer"
	"github.com/johandrevandeventer/kafkaclient/producer"
	"github.com/johandrevandeventer/logging"
	"github.com/johandrevandeventer/lora-worker/internal/flags"
	"go.uber.org/zap"
)

func (e *Engine) startKafkaProducer() {
	e.logger.Info("Starting Kafka producer")

	var kafkaProducerLogger *zap.Logger
	if flags.FlagKafkaLogging {
		kafkaProducerLogger = logging.GetLogger("kafka.producer")
	} else {
		kafkaProducerLogger = zap.NewNop()
	}

	// Define Kafka producer config
	producerConfig := config.NewKafkaProducerConfig("localhost:9092", 5, 5)

	// Initialize Kafka Producer Pool
	kafkaProducerPool, err := producer.NewKafkaProducerPool(e.ctx, producerConfig, kafkaProducerLogger)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer pool: %v", err)
	}

	e.kafkaProducerPool = kafkaProducerPool
}

func (e *Engine) startKafkaConsumer() {
	e.logger.Info("Starting Kafka consumer")

	var kafkaConsumerLogger *zap.Logger
	if flags.FlagKafkaLogging {
		kafkaConsumerLogger = logging.GetLogger("kafka.consumer")
	} else {
		kafkaConsumerLogger = zap.NewNop()
	}

	// Define Kafka consumer config
	var consumerConfig *config.KafkaConsumerConfig
	if flags.FlagEnvironment == "development" {
		consumerConfig = config.NewKafkaConsumerConfig("localhost:9092", "rubicon_kafka_lora_development", "lora-development-consumer-group")
	} else {
		consumerConfig = config.NewKafkaConsumerConfig("localhost:9092", "rubicon_kafka_lora", "lora-consumer-group")
	}

	// Initialize Kafka Consumer Pool
	kafkaConsumer, err := consumer.NewKafkaConsumer(e.ctx, consumerConfig, kafkaConsumerLogger)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer pool: %v", err)
	}

	e.kafkaConsumer = kafkaConsumer

	// // Start Prometheus Metrics Server
	// e.wg.Add(1)
	// go func() {
	// 	defer e.wg.Done()
	// 	prometheusserver.StartPrometheusServer(":2113", e.ctx)
	// }()

	// Start Kafka consumer
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		e.kafkaConsumer.Start()
	}()
}
