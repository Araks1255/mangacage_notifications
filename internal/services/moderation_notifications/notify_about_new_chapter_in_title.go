package moderation_notifications

import (
	"context"
	"fmt"

	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lib/pq"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s server) NotifyAboutNewChapterInTitle(ctx context.Context, chapter *pb.Chapter) (*emptypb.Empty, error) {
	var data struct {
		SubscribedUsersTgIDs pq.Int64Array `gorm:"column:subscribed_users_tg_ids"`
		TitleName            *string
	}

	err := s.DB.Raw(
		`SELECT
			COALESCE(ARRAY_AGG(u.tg_user_id), '{}'::BIGINT[])::BIGINT[] AS subscribed_users_tg_ids,
			t.name AS title_name
		FROM
			chapters AS c
			INNER JOIN titles AS t ON t.id = c.title_id
			INNER JOIN user_titles_subscribed_to AS utst ON utst.title_id = t.id
			INNER JOIN users AS u ON u.id = utst.user_id
		WHERE
			c.id = ?
		GROUP BY
			t.id, c.id`,
		chapter.ID,
	).Scan(&data).Error

	if err != nil {
		return nil, err
	}

	if len(data.SubscribedUsersTgIDs) == 0 {
		return nil, nil
	}

	var dereferencedTitleName string
	if data.TitleName != nil {
		dereferencedTitleName = *data.TitleName
	}

	message := fmt.Sprintf("В тайтле \"%s\" вышла новая глава", dereferencedTitleName)

	for i := 0; i < len(data.SubscribedUsersTgIDs); i++ {
		if _, err := s.Bot.Send(tgbotapi.NewMessage(data.SubscribedUsersTgIDs[i], message)); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
