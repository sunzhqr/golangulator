package telegram

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sunzhqr/golangulator/internal/domain"
)

type BotHandler struct {
	Bot        *tgbotapi.BotAPI
	Calculator domain.Calculator
	History    domain.HistoryUseCase
}

func NewBotHandler(bot *tgbotapi.BotAPI, calculator domain.Calculator, history domain.HistoryUseCase) *BotHandler {
	return &BotHandler{
		Bot:        bot,
		Calculator: calculator,
		History:    history,
	}
}

func (h *BotHandler) HandleUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.Bot.GetUpdatesChan(u)

	for update := range updates {
		// Inline mode
		if update.InlineQuery != nil {
			query := update.InlineQuery.Query
			if query == "" {
				continue
			}

			result, err := h.Calculator.Eval(query)
			var contentText string
			if err != nil {
				contentText = "Ошибка: " + err.Error()
			} else {
				contentText = fmt.Sprintf("%s = %.10g", query, result)
			}

			inlineResult := tgbotapi.NewInlineQueryResultArticle(
				"calc_"+strconv.Itoa(int(update.InlineQuery.From.ID)), // уникальный ID
				"Вычислить: "+query,
				contentText,
			)
			inlineResult.Description = contentText
			inlineResult.InputMessageContent = tgbotapi.InputTextMessageContent{
				Text: contentText,
			}

			inlineConfig := tgbotapi.InlineConfig{
				InlineQueryID: update.InlineQuery.ID,
				IsPersonal:    true,
				CacheTime:     0,
				Results:       []interface{}{inlineResult},
			}

			if _, err := h.Bot.Request(inlineConfig); err != nil {
				log.Println("Ошибка при отправке inline-ответа:", err)
			}
			continue
		}

		if update.Message == nil || update.Message.Text == "" {
			continue
		}

		text := update.Message.Text
		chatID := update.Message.Chat.ID

		switch text {
		case "/start":
			msg := tgbotapi.NewMessage(chatID, "Привет! Я Голангулятор)\nПришли арифметическое выражение, и я его вычислю.")
			msg.ReplyMarkup = defaultKeyboard()
			h.Bot.Send(msg)

		case "/help":
			msg := tgbotapi.NewMessage(chatID, "Пример: 2 + 2 * (3 - 1)\nКоманды:\n/history - показать последние вычисления\n/clear_history - очистить историю")
			msg.ReplyMarkup = defaultKeyboard()
			h.Bot.Send(msg)

		case "/history":
			history, err := h.History.GetUserHistory(chatID)
			if err != nil || len(history) == 0 {
				h.Bot.Send(tgbotapi.NewMessage(chatID, "История пуста"))
				continue
			}

			var sb strings.Builder
			sb.WriteString("Последние вычисления:\n")
			for index, item := range history {
				sb.WriteString(fmt.Sprintf("|%d| %s = %v\n", index+1, item.Expression, item.Result))
			}
			h.Bot.Send(tgbotapi.NewMessage(chatID, sb.String()))

		case "/clear_history":
			err := h.History.ClearUserHistory(chatID)
			if err != nil {
				h.Bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при очистке истории"))
			} else {
				h.Bot.Send(tgbotapi.NewMessage(chatID, "История очищена"))
			}

		default:
			result, err := h.Calculator.Eval(text)
			var reply string
			if err != nil {
				reply = "Ошибка в выражении: " + err.Error()
			} else {
				reply = "Результат: " + strconv.FormatFloat(result, 'f', -1, 64)
				saveErr := h.History.SaveEntry(chatID, text, result)
				if saveErr != nil {
					log.Println(saveErr)
				}
			}
			h.Bot.Send(tgbotapi.NewMessage(chatID, reply))
		}
	}
}
