package site_notifications

import (
	"context"
	"fmt"

	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s server) NotifyAboutTeamJoinRequestResponse(ctx context.Context, response *pb.TeamJoinRequestResponse) (*emptypb.Empty, error) {
	var data struct {
		TeamName *string
		TgUserID *int64
	}

	err := s.DB.Raw(
		`SELECT
			(SELECT name FROM teams WHERE id = ?) AS team_name,
			(SELECT tg_user_id FROM users WHERE id = ?)`,
		response.TeamID, response.UserID,
	).Scan(&data).Error

	if err != nil {
		return nil, err
	}

	if data.TgUserID == nil {
		return nil, nil
	}

	var dereferenedTeamName string
	if data.TeamName != nil {
		dereferenedTeamName = *data.TeamName
	}

	var message string

	switch response.Result {
	case enums.ResultOfTeamJoinRequest_RESULT_OF_TEAM_JOIN_REQUEST_DECLINED:
		message = fmt.Sprintf("Вас не приняли в команду перевода %s по причине: \"%s\"", dereferenedTeamName, response.Reason)

	case enums.ResultOfTeamJoinRequest_RESULT_OF_TEAM_JOIN_REQUEST_APPROVED:
		message = fmt.Sprintf("Вас приняли в команду перевода %s", dereferenedTeamName)
	}

	if _, err := s.Bot.Send(tgbotapi.NewMessage(*data.TgUserID, message)); err != nil {
		return nil, err
	}

	return nil, nil
}
