package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sunzhqr/golangulator/internal/infrastructure/telegram"
	"github.com/sunzhqr/golangulator/internal/usecase"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	botToken := os.Getenv("BOT_TOKEN")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	calculator := usecase.NewCalculatorUseCase()
	handler := telegram.NewBotHandler(bot, calculator)

	log.Println("Бот запущен")
	handler.HandleUpdates()
}
