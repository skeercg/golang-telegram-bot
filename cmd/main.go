package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"golang-telegram-bot/pkg"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	unsplashKey := viper.GetString("unsplashkey")
	imageController := pkg.ImageController{
		UnsplashKey: unsplashKey,
	}
	bot, err := tgbotapi.NewBotAPI(viper.GetString("tgtoken"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		go func() {
			if update.Message != nil {
				chatID := update.Message.Chat.ID

				switch update.Message.Text {
				case "/start":
					msg := tgbotapi.NewMessage(chatID, "Hello, I'm Imager bot!")

					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
				case "/image":
					photoConfig := imageController.GetPhoto(chatID)

					if _, err := bot.Send(photoConfig); err != nil {
						log.Panic(err)
					}
				default:
					msg := tgbotapi.NewMessage(chatID, "Unknown command")

					if _, err := bot.Send(msg); err != nil {
						log.Panic(err)
					}
				}

			}
		}()
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
