package sender

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func (s *Sender) SendMassMessages(msgs []*tgbotapi.MessageConfig) {
	for i := 0; i < len(msgs); i++ {
		s.massMessages <- msgs[i]
	}
}
