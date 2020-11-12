package main

import (
	"blacksmith-watcher/utils"
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/getsentry/sentry-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

func main() {
	var err error

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	switch os.Getenv("CWBW_LOGLEVEL") {
	case "DISABLED":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	case "PANIC":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "FATAL":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "TRACE":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	utils.Client, err = mongo.NewClient(options.Client().ApplyURI(os.Getenv("CWBW_MONGO_URI")))
	if err != nil {
		log.Panic().Err(err).Send()
	}

	err = utils.Client.Connect(context.Background())
	if err != nil {
		log.Panic().Err(err).Send()
	}

	err = utils.Client.Ping(context.Background(), nil)
	if err != nil {
		log.Panic().Err(err).Send()
	}

	utils.DB = utils.Client.Database(utils.GetEnv("CWBW_MONGO_DBNAME", "blacksmith-watcher"))

	SentryDSN := utils.GetEnv("CWBW_SENTRY_DSN", "")
	SentryEnvironment := utils.GetEnv("CWBW_ENVIRONMENT", "production")

	// Sentry init
	err = sentry.Init(sentry.ClientOptions{
		Dsn:         SentryDSN,
		Environment: SentryEnvironment,
	})
	if err != nil {
		log.Warn().Err(err).Send()
	}
	defer sentry.Flush(2 * time.Second)

	// Kafka consumer init
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": utils.GetEnv("CWBW_KAFKA_ADDRESS", "localhost"),
		"group.id":          "cw3",
		"auto.offset.reset": "latest",
	})

	if err != nil {
		log.Panic().Err(err).Send()
		sentry.CaptureException(err)
		return
	}

	err = InitBot(os.Getenv("CWBW_BOT_TOKEN"), consumer)
	if err != nil {
		sentry.CaptureException(err)
		log.Panic().Err(err).Send()
	}
}
