package grpc

import (
	"context"
	"fmt"

	pb "github.com/Araks1255/mangacage_protos"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s Server) NotifyAboutTitleOnModeration(ctx context.Context, title *pb.TitleOnModeration) (*pb.Empty, error) {
	var allowedUsersTgIds []int64
	s.DB.Raw(`SELECT users.tg_user_id FROM users
		INNER JOIN user_roles ON users.id = user_roles.user_id
		INNER JOIN roles ON user_roles.role_id = roles.id
		WHERE roles.name = 'moder' OR roles.name = 'admin'`,
	).Scan(&allowedUsersTgIds)

	var titleName string
	s.DB.Raw("SELECT name FROM titles WHERE id = ?", title.ID).Scan(&titleName)

	var msg tgbotapi.MessageConfig

	var response string
	if title.New {
		response = fmt.Sprintf("На модерацию пришёл новый тайтл\n\nНазвание: %s", titleName)
	} else {
		response = fmt.Sprintf("На модерацию пришли изменения для тайтла %s", titleName)
	}

	for i := 0; i < len(allowedUsersTgIds); i++ {
		msg = tgbotapi.NewMessage(allowedUsersTgIds[i], response)
		s.ServiceBot.Send(msg)
	}

	return &pb.Empty{}, nil
}

func (s Server) NotifyAboutVolumeOnModeration(ctx context.Context, volume *pb.VolumeOnModeration) (*pb.Empty, error) {
	var allowedUsersTgIds []int64
	s.DB.Raw(`SELECT users.tg_user_id FROM users
		INNER JOIN user_roles ON users.id = user_roles.user_id
		INNER JOIN roles ON user_roles.role_id = roles.id
		WHERE roles.name = 'moder' OR roles.name = 'admin'`,
	).Scan(&allowedUsersTgIds)

	var volumeName string
	s.DB.Raw("SELECT name FROM volumes WHERE id = ?", volume.ID).Scan(&volumeName)

	var msg tgbotapi.MessageConfig

	var response string
	if volume.New {
		response = fmt.Sprintf("На модерацию пришёл новый том\n\nНазвание: %s", volumeName)
	} else {
		response = fmt.Sprintf("На модерацию пришли изменения для тома %s", volumeName)
	}

	for i := 0; i < len(allowedUsersTgIds); i++ {
		msg = tgbotapi.NewMessage(allowedUsersTgIds[i], response)
		s.ServiceBot.Send(msg)
	}

	return &pb.Empty{}, nil
}

func (s Server) NotifyAboutChapterOnModeration(ctx context.Context, chapter *pb.ChapterOnModeration) (*pb.Empty, error) {
	var allowedUsersTgIds []int64
	s.DB.Raw(`SELECT users.tg_user_id FROM users
		INNER JOIN user_roles ON users.id = user_roles.user_id
		INNER JOIN roles ON user_roles.role_id = roles.id
		WHERE roles.name = 'moder' OR roles.name = 'admin'`,
	).Scan(&allowedUsersTgIds)

	var chapterName string
	s.DB.Raw("SELECT name FROM chapters WHERE id = ?", chapter.ID).Scan(&chapterName)

	var msg tgbotapi.MessageConfig

	var response string
	if chapter.New {
		response = fmt.Sprintf("На модерацию пришла новая глава\n\nНазвание: %s", chapterName)
	} else {
		response = fmt.Sprintf("На модерацию пришли изменения для главы %s", chapterName)
	}
	for i := 0; i < len(allowedUsersTgIds); i++ {
		msg = tgbotapi.NewMessage(allowedUsersTgIds[i], response)
		s.ServiceBot.Send(msg)
	}

	return &pb.Empty{}, nil
}

func (s Server) NotifyAboutUserOnModeration(ctx context.Context, user *pb.User) (*pb.Empty, error) {
	var allowedUsersTgIds []int64
	s.DB.Raw(`SELECT users.tg_user_id FROM users
		INNER JOIN user_roles ON users.id = user_roles.user_id
		INNER JOIN roles ON user_roles.role_id = roles.id
		WHERE roles.name = 'moder' OR roles.name = 'admin'`,
	).Scan(&allowedUsersTgIds)

	var userName string
	s.DB.Raw("SELECT user_name FROM users WHERE id = ?", user.ID).Scan(&userName)

	var msg tgbotapi.MessageConfig

	var response string
	if user.New {
		response = fmt.Sprintf("Верификации ожидает новый пользователь\n\nИмя: %s", userName)
	} else {
		response = fmt.Sprintf("Подтверждения изменений аккаунта ожидает пользователь %s", userName)
	}
	for i := 0; i < len(allowedUsersTgIds); i++ {
		msg = tgbotapi.NewMessage(allowedUsersTgIds[i], response)
		s.ServiceBot.Send(msg)
	}

	return &pb.Empty{}, nil
}
