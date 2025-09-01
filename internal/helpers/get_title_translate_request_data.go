package helpers

import (
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	"gorm.io/gorm"
)

func GetTitleTranslateRequestData(db *gorm.DB, translateRequest *pb.TitleTranslateRequest) (teamLeaderTgID *int64, titleName string, err error) {
	var data struct {
		TeamLeaderTgID *int64
		TitleName      *string
	}

	query :=
		`SELECT 
			(
				SELECT
					u.id
				FROM
					users AS u
					INNER JOIN teams AS t ON t.id = u.team_id
					INNER JOIN user_roles AS ur ON ur.user_id = u.id
					INNER JOIN roles AS r ON r.id = ur.role_id
				WHERE
					t.id = ? AND r.name = 'team_leader'
			) AS team_leader_tg_id,
			name AS title_name
		FROM
			titles
		WHERE
			id = ?`

	if err := db.Raw(query, translateRequest.TeamID, translateRequest.TitleID).Scan(&data).Error; err != nil {
		return nil, "", err
	}

	if data.TeamLeaderTgID == nil {
		return nil, "", nil
	}

	var dereferencedTitleName string
	if data.TitleName != nil {
		dereferencedTitleName = *data.TitleName
	}

	return data.TeamLeaderTgID, dereferencedTitleName, nil
}
