package moderation_notifications

import (
	"context"
	"errors"
	"fmt"

	"github.com/Araks1255/mangacage_notifications/internal/helpers"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s server) SendModerationRequestDeclineReason(ctx context.Context, reason *pb.ModerationRequestDeclineReason) (*emptypb.Empty, error) {
	tgReceiverID, err := helpers.GetTgUserID(s.DB, uint(reason.CreatorID))
	if err != nil {
		return nil, err
	}

	if tgReceiverID == 0 {
		return nil, errors.New("Пользователь не может получать сообщения")
	}

	var entityForMessage string

	switch reason.EntityOnModeration {
	case enums.EntityOnModeration_ENTITY_ON_MODERATION_AUTHOR:
		entityForMessage = "автора"
	case enums.EntityOnModeration_ENTITY_ON_MODERATION_TITLE:
		entityForMessage = "тайтла"
	case enums.EntityOnModeration_ENTITY_ON_MODERATION_CHAPTER:
		entityForMessage = "главы"
	case enums.EntityOnModeration_ENTITY_ON_MODERATION_GENRE:
		entityForMessage = "жанра"
	case enums.EntityOnModeration_ENTITY_ON_MODERATION_TAG:
		entityForMessage = "тега"
	case enums.EntityOnModeration_ENTITY_ON_MODERATION_PROFILE_CHANGES:
		entityForMessage = "изменений профиля"
	case enums.EntityOnModeration_ENTITY_ON_MODERATION_TEAM:
		entityForMessage = "команды"
	}

	message := fmt.Sprintf("Заявка на модерацию %s \"%s\" была отклонена по причине: %s", entityForMessage, reason.EntityName, reason.Reason)

	msg := tgbotapi.NewMessage(tgReceiverID, message)

	s.Sender.SendSingleMessage(&msg)

	return nil, nil
}
