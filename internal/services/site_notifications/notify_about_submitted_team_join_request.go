package site_notifications

import (
	"context"
	"fmt"

	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s server) NotifyAboutSubmittedTeamJoinRequest(ctx context.Context, joinRequest *pb.TeamJoinRequest) (*emptypb.Empty, error) {
	var data struct {
		TeamLeaderTgUserID *int64
		CandidateName      *string
	}

	err := s.DB.Raw(
		`SELECT
			(
				SELECT
					u.tg_user_id
				FROM
					users AS u
					INNER JOIN user_roles AS ur ON ur.user_id = u.id
					INNER JOIN roles AS r ON r.id = ur.role_id
				WHERE
					u.team_id = ? AND r.name = 'team_leader'
			) AS team_leader_tg_user_id,
			(SELECT user_name FROM users WHERE id = ?) AS candidate_name`,
		joinRequest.TeamID, joinRequest.CandidateID,
	).Scan(&data).Error

	if err != nil {
		return nil, err
	}

	if data.TeamLeaderTgUserID == nil {
		return nil, nil
	}

	var dereferencedCandidateName string
	if data.CandidateName != nil {
		dereferencedCandidateName = *data.CandidateName
	}

	message := fmt.Sprintf("Пользователь %s отправил заявку на вступление в вашу команду", dereferencedCandidateName)

	if joinRequest.Message != "" {
		message += fmt.Sprintf("\nВступительное сообщение: %s", joinRequest.Message)
	}

	if _, err := s.Bot.Send(tgbotapi.NewMessage(*data.TeamLeaderTgUserID, message)); err != nil {
		return nil, err
	}

	return nil, nil
}
