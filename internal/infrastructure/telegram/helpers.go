package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

func defaultKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/history"),
			tgbotapi.NewKeyboardButton("/clear_history"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/help"),
		),
	)
}

func (h *BotHandler) InitBotCommands() {
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Начать работу"},
		{Command: "help", Description: "Справка"},
		{Command: "history", Description: "Показать историю"},
		{Command: "clear_history", Description: "Очистить историю"},
	}
	cfg := tgbotapi.NewSetMyCommands(commands...)
	if _, err := h.Bot.Request(cfg); err != nil {
		h.Logger.Error("Ошибка установки команд:", zap.Error(err))
	}
}
