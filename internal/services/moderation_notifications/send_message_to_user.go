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

func (s server) SendMessageToUser(ctx context.Context, messageFromModerator *pb.MessageFromModerator) (*emptypb.Empty, error) {
	tgReceiverID, err := helpers.GetTgUserID(s.DB, uint(messageFromModerator.ReceiverID))
	if err != nil {
		return nil, err
	}

	if tgReceiverID == 0 {
		return nil, errors.New("Пользователь не может принимать сообщения")
	}

	var entityForLink string

	switch messageFromModerator.EntityOnModeration {
	case enums.EntityOnModeration_ENTITY_ON_MODERATION_TITLE:
		entityForLink = "titles"
	case enums.EntityOnModeration_ENTITY_ON_MODERATION_CHAPTER:
		entityForLink = "chapters"
	}

	message := fmt.Sprintf("У вас новое сообщение от модератора: %s", messageFromModerator.Text)

	// if entityForLink != "" {
	// 	linkButton := tgbotapi.NewInlineKeyboardButtonURL("Ссылка на заявку", attachedLink)
	// 	row := tgbotapi.NewInlineKeyboardRow(linkButton)
	// 	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(row)
	// 	msg.ReplyMarkup = inlineKeyboard
	// } // В кнопки нельзя localhost вставлять

	if entityForLink != "" {
		attachedLink := fmt.Sprintf("http://localhost:8080/users/me/moderation/%s/%d", entityForLink, messageFromModerator.EntityOnModerationID)
		message += fmt.Sprintf("\nСсылка на заявку - %s", attachedLink)
	}

	msg := tgbotapi.NewMessage(tgReceiverID, message)

	s.Sender.SendSingleMessage(&msg)

	return nil, nil
}
