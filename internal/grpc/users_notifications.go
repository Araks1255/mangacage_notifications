package grpc

import (
	"context"
	"errors"
	"fmt"
	"log"

	promoUtils "github.com/Araks1255/mangacage_promocodes/pkg/common/utils"
	pb "github.com/Araks1255/mangacage_protos"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

func (s Server) NotifyAboutReleaseOfNewChapterInTitle(ctx context.Context, chapter *pb.ReleasedChapter) (*pb.Empty, error) {
	var chapterTitleID uint
	db.Raw("SELECT title_id FROM chapters WHERE name = ?", chapter.Name).Scan(&chapterTitleID)

	var chapterTitleName string
	db.Raw("SELECT name FROM titles WHERE id = ?", chapterTitleID).Scan(&chapterTitleName)

	var subscribedUsersIDs []int64
	db.Raw("SELECT users.tg_user_id FROM users INNER JOIN user_titles_subscribed_to ON users.id = user_titles_subscribed_to.user_id WHERE user_titles_subscribed_to.title_id = ?", chapterTitleID).Scan(&subscribedUsersIDs)

	var msg tgbotapi.MessageConfig
	response := fmt.Sprintf("Вышла новая глава в тайтле %s\n\nНазвание - %s", chapterTitleName, chapter.Name)
	for i := 0; i < len(subscribedUsersIDs); i++ {
		msg = tgbotapi.NewMessage(subscribedUsersIDs[i], response)
		if _, err := usersBot.Send(msg); err != nil {
			return &pb.Empty{}, err
		}
	}

	return &pb.Empty{}, nil
}

func (s Server) SendPromocode(ctx context.Context, promocodeRequest *pb.PromocodeRequest) (*pb.Empty, error) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	key := []byte(viper.Get("PROMO_KEY").(string))

	var tgUserID int64
	db.Raw("SELECT tg_user_id FROM users WHERE name = ?", promocodeRequest.User.Name).Scan(&tgUserID)
	if tgUserID == 0 {
		return &pb.Empty{}, errors.New("Пользователь не найден")
	}

	decryptedCode, err := promoUtils.DecryptPromocode(promocodeRequest.Promocode.Code, key)
	if err != nil {
		log.Println(err)
		return &pb.Empty{}, err
	}

	msg := tgbotapi.NewMessage(tgUserID,
		fmt.Sprintf("Здравствуйте %s!\n\nВот ваш промокод на %s алмазиков:\n%s",
			promocodeRequest.User.Name,
			promocodeRequest.Promocode.Amount,
			decryptedCode))

	if _, err = usersBot.Send(msg); err != nil {
		log.Println(err)
		return &pb.Empty{}, err
	}

	return &pb.Empty{}, nil
}
