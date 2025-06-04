package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sunzhqr/golangulator/internal/infrastructure/storage"
	"github.com/sunzhqr/golangulator/internal/infrastructure/telegram"
	"github.com/sunzhqr/golangulator/internal/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	db.AutoMigrate(&storage.HistoryEntryModel{})

	historyRepo := storage.NewPostgresRepo(db)
	historyUseCase := usecase.NewHistoryUseCase(historyRepo)

	handler := telegram.NewBotHandler(bot, calculator, historyUseCase)
	handler.InitBotCommands()
	log.Println("Бот запущен")
	handler.HandleUpdates()
}
