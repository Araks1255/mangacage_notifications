package moderation_notifications

import (
	"context"
	"fmt"

	"github.com/Araks1255/mangacage_notifications/internal/sender"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s server) NotifyAboutApprovedModerationRequest(ctx context.Context, approvedRequest *pb.ApprovedEntity) (*emptypb.Empty, error) {
	if approvedRequest.Entity == enums.Entity_ENTITY_PROFILE {
		err := notifyAboutApprovedProfileChanges(s.DB, s.Sender, uint(approvedRequest.CreatorID))
		return nil, err
	}

	var data struct {
		TgUserID   *int64
		EntityName *string
	}

	query := `SELECT
				(SELECT tg_user_id FROM users WHERE id = ?),
				(SELECT name FROM %s WHERE id = ?) AS entity_name`

	message := "Ваша заявка на модерацию %s \"%s\" одобрена"

	var (
		entityForQuery   string
		entityForMessage string
	)

	switch approvedRequest.Entity {
	case enums.Entity_ENTITY_TITLE:
		entityForQuery, entityForMessage = "titles", "тайтла"

	case enums.Entity_ENTITY_CHAPTER:
		entityForQuery, entityForMessage = "chapters", "главы"

	case enums.Entity_ENTITY_AUTHOR:
		entityForQuery, entityForMessage = "authors", "автора"

	case enums.Entity_ENTITY_TEAM:
		entityForQuery, entityForMessage = "teams", "команды"

	case enums.Entity_ENTITY_GENRE:
		entityForQuery, entityForMessage = "genres", "жанра"

	case enums.Entity_ENTITY_TAG:
		entityForQuery, entityForMessage = "tags", "тега"

	default:
		return nil, fmt.Errorf("неподдерживаемый тип сущности: %s", approvedRequest.Entity.String())
	}

	query = fmt.Sprintf(query, entityForQuery)

	if err := s.DB.Raw(query, approvedRequest.CreatorID, approvedRequest.ID).Scan(&data).Error; err != nil {
		return nil, err
	}

	if data.TgUserID == nil {
		return nil, nil
	}

	var dereferencedEntityName string
	if data.EntityName != nil {
		dereferencedEntityName = *data.EntityName
	}

	message = fmt.Sprintf(message, entityForMessage, dereferencedEntityName)

	msg := tgbotapi.NewMessage(*data.TgUserID, message)

	s.Sender.SendSingleMessage(&msg)

	return nil, nil
}

func notifyAboutApprovedProfileChanges(db *gorm.DB, sender *sender.Sender, creatorID uint) error {
	var tgUserID *int64

	if err := db.Raw("SELECT tg_user_id FROM users WHERE id = ?", creatorID).Scan(&tgUserID).Error; err != nil {
		return err
	}

	if tgUserID == nil {
		return nil
	}

	msg := tgbotapi.NewMessage(*tgUserID, "Ваша заявка на изменение профиля одобрена")

	sender.SendSingleMessage(&msg)

	return nil
}
