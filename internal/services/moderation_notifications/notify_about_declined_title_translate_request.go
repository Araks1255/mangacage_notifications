package moderation_notifications

import (
	"context"
	"fmt"

	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s server) NotifyAboutDeclinedTitleTranslateRequest(ctx context.Context, translateRequest *pb.TitleTranslateRequest) (*emptypb.Empty, error) {
	var data struct {
		SenderTgID *int64
		TitleName  *string
	}

	err := s.DB.Raw(
		`SELECT
			(SELECT tg_user_id FROM users WHERE id = ?) AS sender_tg_id,
			(SELECT name FROM titles WHERE id = ?) AS title_name`,
		translateRequest.SenderID, translateRequest.TitleID,
	).Scan(&data).Error

	if err != nil {
		return nil, err
	}

	if data.SenderTgID == nil {
		return nil, nil
	}

	var dereferencedTitleName string
	if data.TitleName != nil {
		dereferencedTitleName = *data.TitleName
	}

	message := fmt.Sprintf("Ваша заявка на перевод тайтла \"%s\" отклонена", dereferencedTitleName)

	msg := tgbotapi.NewMessage(*data.SenderTgID, message)

	s.Sender.SendSingleMessage(&msg)

	return nil, nil
}
