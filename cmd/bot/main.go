package main

import (
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/sunzhqr/golangulator/internal/infrastructure/storage"
	"github.com/sunzhqr/golangulator/internal/infrastructure/telegram"
	"github.com/sunzhqr/golangulator/internal/usecase"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Ошибка загрузки .env файла", zap.Error(err))
	}

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		logger.Fatal("BOT_TOKEN не найден в переменных окружения")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logger.Fatal("Ошибка создания Telegram-бота", zap.Error(err))
	}
	logger.Info("Бот авторизован", zap.String("bot_username", bot.Self.UserName))

	calculator := usecase.NewCalculatorUseCase()

	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("Ошибка подключения к базе данных", zap.Error(err))
	}
	if err := db.AutoMigrate(&storage.HistoryEntryModel{}); err != nil {
		logger.Fatal("Ошибка миграции таблицы", zap.Error(err))
	}
	logger.Info("Подключение к базе данных успешно")

	historyRepo := storage.NewPostgresRepo(db)
	historyUseCase := usecase.NewHistoryUseCase(historyRepo)

	handler := telegram.NewBotHandler(bot, calculator, historyUseCase, logger)
	handler.InitBotCommands()
	logger.Info("Бот запущен и готов к работе")
	handler.HandleUpdates()
}
