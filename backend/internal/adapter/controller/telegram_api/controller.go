package telegram_api

import (
	"fmt"
	"strings"
	"sync"

	"github.com/andreychh/coopera-backend/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramController struct {
	bot    *tgbotapi.BotAPI
	logger *logger.Logger
	_      sync.Mutex
}

func NewTelegramController(logger *logger.Logger, botToken string) (*TelegramController, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Telegram bot: %w", err)
	}

	return &TelegramController{
		bot:    bot,
		logger: logger,
	}, nil
}

func (c *TelegramController) Start() error {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error("Recovered from panic: %v", r)
		}
	}()

	c.logger.Info("Starting Telegram bot...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := c.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || update.Message.From == nil {
			continue
		}

		if strings.EqualFold(update.Message.Text, "привет") {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, сервис работает")
			_, err := c.bot.Send(msg)
			if err != nil {
				c.logger.Error("Failed to send message: %v", err)
			}
		}
	}

	return nil
}
