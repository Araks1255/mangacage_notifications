package site_notifications

import (
	"context"

	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s server) NotifyAboutTitleTranslateRequest(ctx context.Context, translateRequest *pb.TitleTranslateRequest) (*emptypb.Empty, error) {
	message := "Пришел новый запрос на перевод тайтла"

	s.mu.RLock()

	msgs := make([]*tgbotapi.MessageConfig, len(*s.ModersTgIDs))

	for i := 0; i < len(msgs); i++ {
		msgs[i] = &tgbotapi.MessageConfig{
			Text:     message,
			BaseChat: tgbotapi.BaseChat{ChatID: (*s.ModersTgIDs)[i]},
		}
	}

	s.mu.RUnlock()

	s.Sender.SendMassMessages(msgs)

	return nil, nil
}
