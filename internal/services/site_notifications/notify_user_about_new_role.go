package site_notifications

import (
	"context"
	"fmt"

	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s server) NotifyUserAboutNewRole(ctx context.Context, role *pb.NewRole) (*emptypb.Empty, error) {
	var data struct {
		TgUserID *int64
		RoleName *string
	}

	err := s.DB.Raw(
		`SELECT
			(SELECT tg_user_id FROM users WHERE id = ?),
			(SELECT name FROM roles WHERE id = ?) AS role_name`,
		role.UserID, role.RoleID,
	).Scan(&data).Error

	if err != nil {
		return nil, err
	}

	if data.TgUserID == nil {
		return nil, nil
	}

	var dereferencedRoleName string
	if data.RoleName != nil {
		dereferencedRoleName = *data.RoleName
	}

	message := fmt.Sprintf("Вам назначена новая роль - %s", dereferencedRoleName)

	if _, err := s.Bot.Send(tgbotapi.NewMessage(*data.TgUserID, message)); err != nil {
		return nil, err
	}

	return nil, nil
}
