package main

import (
	"cw-broker/messages"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
	"time"
)

func InitBot(
	telegramToken string,
	logger *logrus.Logger,
	consumer *kafka.Consumer,
) error {
	bot, err := telebot.NewBot(
		telebot.Settings{
			Token:  telegramToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})
	if err != nil {
		return err
	}

	bot.Handle("/start", func(message *telebot.Message) {

	})

	consumer.SubscribeTopics([]string{"cw3-yellow_pages"}, nil)

	defer bot.Start()

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			var message messages.BlacksmithShop
			err = json.Unmarshal([]byte(msg.Value), &message)
			if err != nil {
				sentry.CaptureException(err)
				logger.Error(fmt.Sprintf("Decoder error: %v (%v)\n", err, msg))
			}

			// TODO: Shops parsing

			if err != nil {
				sentry.CaptureException(err)
				logger.Error(err)
			}
			logger.Trace(fmt.Sprintf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value)))
		} else {
			sentry.CaptureException(err)
			logger.Error(fmt.Sprintf("Consumer error: %v (%v)\n", err, msg))
		}
	}
}
