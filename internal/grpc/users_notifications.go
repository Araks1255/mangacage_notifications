package grpc

import (
	"context"
	//"errors"
	"fmt"
	//"log"

	//promoUtils "github.com/Araks1255/mangacage_promocodes/pkg/common/utils"
	pb "github.com/Araks1255/mangacage_protos"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//"github.com/spf13/viper"
)

func (s Server) NotifyAboutReleaseOfNewChapterInTitle(ctx context.Context, chapter *pb.ReleasedChapter) (*pb.Empty, error) {
	var inf struct {
		Chapter              string
		Title                string
		Volume               string
		SubscribedUsersTgIDs []int64
	}

	s.DB.Raw(
		`SELECT c.name AS chapter, volumes.name AS volume, titles.name AS title,
		(
			SELECT ARRAY (
				SELECT users.tg_user_id FROM users
				INNER JOIN user_titles_subscribed_to ON users.id = user_titles_subscribed_to.user_id
				INNER JOIN titles ON titles.id = user_titles_subscribed_to.title_id
				INNER JOIN volumes ON volumes.title_id = titles.id
				INNER JOIN chapters ON chapters.volume_id = volumes.id
				WHERE chapters.id = c.id
			) AS subscribed_users_tg_ids
		) FROM chapters AS c
		 INNER JOIN volumes ON c.volume_id = volumes.id
		 INNER JOIN titles ON volumes.title_id = titles.id
		 WHERE c.id = ?`, chapter.ID,
	).Scan(&inf)

	var msg tgbotapi.MessageConfig
	response := fmt.Sprintf("Вышла новая глава в тайтле %s\n\nНазвание - %s", inf.Title, inf.Chapter)

	for i := 0; i < len(inf.SubscribedUsersTgIDs); i++ {
		msg = tgbotapi.NewMessage(inf.SubscribedUsersTgIDs[i], response)
		s.UsersBot.Send(msg)
	}

	return &pb.Empty{}, nil
}

// func (s Server) SendPromocode(ctx context.Context, promocodeRequest *pb.PromocodeRequest) (*pb.Empty, error) {
// 	viper.SetConfigFile("./pkg/common/envs/.env")
// 	viper.ReadInConfig()

// 	key := []byte(viper.Get("PROMO_KEY").(string))

// 	var tgUserID int64
// 	s.DB.Raw("SELECT tg_user_id FROM users WHERE name = ?", promocodeRequest.User.Name).Scan(&tgUserID)
// 	if tgUserID == 0 {
// 		return &pb.Empty{}, errors.New("Пользователь не найден")
// 	}

// 	decryptedCode, err := promoUtils.DecryptPromocode(promocodeRequest.Promocode.Code, key)
// 	if err != nil {
// 		log.Println(err)
// 		return &pb.Empty{}, err
// 	}

// 	msg := tgbotapi.NewMessage(tgUserID,
// 		fmt.Sprintf("Здравствуйте %s!\n\nВот ваш промокод на %s алмазиков:\n%s",
// 			promocodeRequest.User.Name,
// 			promocodeRequest.Promocode.Amount,
// 			decryptedCode))

// 	if _, err = s.UsersBot.Send(msg); err != nil {
// 		log.Println(err)
// 		return &pb.Empty{}, err
// 	}

// 	return &pb.Empty{}, nil
// }
