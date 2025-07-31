package site_notifications

import (
	"context"

	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s server) NotifyAboutNewUserOnVerification(ctx context.Context, user *pb.UserOnVerification) (*emptypb.Empty, error) {
	message := "Новый пользователь ожидает верификации"

	s.Mu.RLock()

	for i := 0; i < len(*s.ModersTgIDs); i++ {
		if _, err := s.Bot.Send(tgbotapi.NewMessage((*s.ModersTgIDs)[i], message)); err != nil {
			return nil, err
		}
	}

	s.Mu.RUnlock()

	return nil, nil
}
