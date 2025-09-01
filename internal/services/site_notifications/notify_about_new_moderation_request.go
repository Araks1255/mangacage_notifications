package site_notifications

import (
	"context"

	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s server) NotifyAboutNewModerationRequest(ctx context.Context, moderationRequest *pb.ModerationRequest) (*emptypb.Empty, error) {
	var message string

	switch moderationRequest.EntityOnModeration {
	case enums.EntityOnModeration_ENTITY_ON_MODERATION_AUTHOR:
		message = "На модерацию пришел автор"
	case enums.EntityOnModeration_ENTITY_ON_MODERATION_TITLE:
		message = "На модерацию пришел тайтл"
	case enums.EntityOnModeration_ENTITY_ON_MODERATION_CHAPTER:
		message = "На модерацию пришла глава"
	case enums.EntityOnModeration_ENTITY_ON_MODERATION_GENRE:
		message = "На модерацию пришел жанр"
	case enums.EntityOnModeration_ENTITY_ON_MODERATION_TAG:
		message = "На модерацию пришел тег"
	case enums.EntityOnModeration_ENTITY_ON_MODERATION_PROFILE_CHANGES:
		message = "На модерацию пришли изменения профиля"
	case enums.EntityOnModeration_ENTITY_ON_MODERATION_TEAM:
		message = "На модерацию пришла команда"
	}

	s.mu.RLock()

	msgs := make([]*tgbotapi.MessageConfig, len(*s.ModersTgIDs))

	for i := 0; i < len(msgs); i++ {
		msgs[i] = &tgbotapi.MessageConfig{
			BaseChat: tgbotapi.BaseChat{
				ChatID:           (*s.ModersTgIDs)[i],
				ReplyToMessageID: 0,
			},
			Text: message,
		}
	}

	s.mu.RUnlock()

	s.Sender.SendMassMessages(msgs)

	return nil, nil
}
