package moderation_notifications

import (
	"context"

	"github.com/Araks1255/mangacage_notifications/internal/helpers"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s server) NotifyUserAboutVerificatedAccount(ctx context.Context, user *pb.VerificatedUser) (*emptypb.Empty, error) {
	tgUserID, err := helpers.GetTgUserID(s.DB, uint(user.ID))

	if err != nil {
		return nil, err
	}

	if tgUserID == 0 {
		return nil, nil
	}

	msg := tgbotapi.NewMessage(tgUserID, "Ваш аккаунт прошел верификацию")

	s.Sender.SendSingleMessage(&msg)

	return nil, nil
}
