package sender

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Sender struct {
	bot            *tgbotapi.BotAPI
	singleMessages chan *tgbotapi.MessageConfig // Сообщения для одного пользователя. Одобрение заявки на модерацию, приём в команду и т.п.
	massMessages   chan *tgbotapi.MessageConfig // Сообщения для множества пользователей. Новая заявка на модерацию, новая глава в тайтле и т.п.
}

func NewSender(bot *tgbotapi.BotAPI) *Sender {
	return &Sender{
		bot:            bot,
		singleMessages: make(chan *tgbotapi.MessageConfig, 512),
		massMessages:   make(chan *tgbotapi.MessageConfig, 8192),
	}
}
