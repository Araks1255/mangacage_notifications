package grpc

import (
	"context"
	"fmt"

	pb "github.com/Araks1255/mangacage_protos"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func (s Server) PasswordRecovery(ctx context.Context, user *pb.User) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
