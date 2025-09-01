package moderation_notifications

import (
	"context"
	"fmt"
	"log"

	"github.com/Araks1255/mangacage_notifications/internal/helpers"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s server) NotifyAboutDeclinedTitleTranslateRequest(ctx context.Context, translateRequest *pb.TitleTranslateRequest) (*emptypb.Empty, error) {
	teamLeaderTgID, titleName, err := helpers.GetTitleTranslateRequestData(s.DB, translateRequest)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if teamLeaderTgID == nil {
		return nil, nil
	}

	message := fmt.Sprintf("Ваша заявка на перевод тайтла \"%s\" отклонена", titleName)

	msg := tgbotapi.NewMessage(*teamLeaderTgID, message)

	s.Sender.SendSingleMessage(&msg)

	return nil, nil
}
