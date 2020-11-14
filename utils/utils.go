package utils

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/tucnak/telebot.v2"
	"os"
)

var Client *mongo.Client
var DB *mongo.Database
var Bot *telebot.Bot
var (
	Unsubscribe       = &telebot.ReplyMarkup{}
	UnsubscribeButton = Unsubscribe.Data("Отписаться", "unsub")

	Subscribe       = &telebot.ReplyMarkup{}
	SubscribeButton = Subscribe.Data("Отписаться", "unsub")
)

func GetEnv(key string, def string) string {
	env := os.Getenv(key)
	if env == "" {
		return def
	}
	return env
}

func Contains(array []string, find string) bool {
	for _, a := range array {
		if a == find {
			return true
		}
	}
	return false
}

func Remove(s []string, item string) []string {
	index := 0
	for _, i := range s {
		if i != item {
			s[index] = i
			index++
		}
	}
	return s[:index]
}
