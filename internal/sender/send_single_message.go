package sender

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (s *Sender) SendSingleMessage(msg *tgbotapi.MessageConfig) {
	s.singleMessages <- msg
}
