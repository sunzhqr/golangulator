package telegram

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sunzhqr/golangulator/internal/domain"
	"go.uber.org/zap"
)

type BotHandler struct {
	Bot        *tgbotapi.BotAPI
	Calculator domain.Calculator
	History    domain.HistoryUseCase
	Logger     *zap.Logger
}

func NewBotHandler(bot *tgbotapi.BotAPI, calculator domain.Calculator, history domain.HistoryUseCase, logger *zap.Logger) *BotHandler {
	return &BotHandler{
		Bot:        bot,
		Calculator: calculator,
		History:    history,
		Logger:     logger,
	}
}

func (h *BotHandler) HandleUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := h.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.InlineQuery != nil {
			h.handleInline(update.InlineQuery)
			continue
		}

		if update.Message == nil || update.Message.Text == "" {
			continue
		}

		h.handleCommand(update.Message.Chat.ID, update.Message.Text)
	}
}

func (h *BotHandler) handleCommand(chatID int64, text string) {
	switch text {
	case "/start":
		h.sendMessage(chatID, "Привет! Я Голангулятор)\nПришли арифметическое выражение, и я его вычислю.", defaultKeyboard())

	case "/help":
		h.sendMessage(chatID, "Пример: 2 + 2 * (3 - 1)\nКоманды:\n/history - показать последние вычисления\n/clear_history - очистить историю", defaultKeyboard())

	case "/history":
		h.handleHistory(chatID)

	case "/clear_history":
		if err := h.History.ClearUserHistory(chatID); err != nil {
			h.sendMessage(chatID, "Ошибка при очистке истории")
		} else {
			h.sendMessage(chatID, "История очищена")
		}

	default:
		h.handleExpression(chatID, text)
	}
}

func (h *BotHandler) handleExpression(chatID int64, expression string) {
	result, err := h.Calculator.Eval(expression)
	if err != nil {
		h.sendMessage(chatID, err.Error())
		return
	}

	reply := "Результат: " + strconv.FormatFloat(result, 'f', -1, 64)
	h.sendMessage(chatID, reply)

	if err := h.History.SaveEntry(chatID, expression, result); err != nil {
		h.Logger.Error("Ошибка при сохранении истории", zap.Error(err))
	}
}

func (h *BotHandler) handleHistory(chatID int64) {
	history, err := h.History.GetUserHistory(chatID)
	if err != nil || len(history) == 0 {
		h.sendMessage(chatID, "История пуста")
		return
	}

	var sb strings.Builder
	sb.WriteString("Последние вычисления:\n")
	for i, entry := range history {
		sb.WriteString(fmt.Sprintf("|%d| %s = %v\n", i+1, entry.Expression, entry.Result))
	}

	h.sendMessage(chatID, sb.String())
}

func (h *BotHandler) handleInline(query *tgbotapi.InlineQuery) {
	if query.Query == "" {
		return
	}

	result, err := h.Calculator.Eval(query.Query)
	var contentText string
	if err != nil {
		contentText = "Ошибка: " + err.Error()
	} else {
		contentText = fmt.Sprintf("%s = %.10g", query.Query, result)
	}

	hash := fmt.Sprintf("%x", md5.Sum([]byte(query.Query)))
	inlineResult := tgbotapi.NewInlineQueryResultArticle(
		"calc_"+strconv.Itoa(int(query.From.ID))+"_"+hash,
		"Вычислить: "+query.Query,
		contentText,
	)
	inlineResult.Description = contentText
	inlineResult.InputMessageContent = tgbotapi.InputTextMessageContent{Text: contentText}

	inlineConfig := tgbotapi.InlineConfig{
		InlineQueryID: query.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       []interface{}{inlineResult},
	}

	if _, err := h.Bot.Request(inlineConfig); err != nil {
		h.Logger.Error("Ошибка при отправке inline-ответа", zap.Error(err))
	}
}

func (h *BotHandler) sendMessage(chatID int64, text string, markup ...tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	if len(markup) > 0 {
		msg.ReplyMarkup = markup[0]
	}
	if _, err := h.Bot.Send(msg); err != nil {
		h.Logger.Error("Ошибка при отправке сообщения", zap.Error(err))
	}
}
