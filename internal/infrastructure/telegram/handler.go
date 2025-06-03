package telegram

import (
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sunzhqr/golangulator/internal/domain"
)

type BotHandler struct {
	Bot        *tgbotapi.BotAPI
	Calculator domain.Calculator
}

func NewBotHandler(bot *tgbotapi.BotAPI, calculator domain.Calculator) *BotHandler {
	return &BotHandler{Bot: bot, Calculator: calculator}
}

func (h *BotHandler) HandleUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		text := update.Message.Text

		switch text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я голангулятор)\nПришли выражение:")
			h.Bot.Send(msg)
		case "/help":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Пример: 2 + 2 * (3 - 1)")
			h.Bot.Send(msg)
		default:
			result, err := h.Calculator.Eval(text)
			var reply string
			if err != nil {
				reply = "Ошибка в выражении: " + err.Error()
			} else {
				reply = "Результат: " + strconv.FormatFloat(result, 'f', -1, 64)
			}
			h.Bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, reply))
		}
	}
}
