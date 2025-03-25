package grpc

import (
	"context"
	"fmt"

	pb "github.com/Araks1255/mangacage_protos"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (s Server) NotifyAboutTitleOnModeration(ctx context.Context, title *pb.TitleOnModeration) (*pb.Empty, error) {
	var allowedUsersTgIds []int64
	db.Raw(`SELECT users.tg_user_id FROM users
		INNER JOIN user_roles ON users.id = user_roles.user_id
		INNER JOIN roles ON user_roles.role_id = roles.id
		WHERE roles.name = 'moder' OR roles.name = 'admin'`,
	).Scan(&allowedUsersTgIds)

	var msg tgbotapi.MessageConfig

	var response string
	if title.New {
		response = fmt.Sprintf("На модерацию пришёл новый тайтл\n\nНазвание: %s", title.Name)
	} else {
		response = fmt.Sprintf("На модерацию пришли изменения для тайтла %s", title.Name)
	}

	for i := 0; i < len(allowedUsersTgIds); i++ {
		msg = tgbotapi.NewMessage(allowedUsersTgIds[i], response)
		serviceBot.Send(msg)
	}

	return &pb.Empty{}, nil
}

func (s Server) NotifyAboutVolumeOnModeration(ctx context.Context, volume *pb.VolumeOnModeration) (*pb.Empty, error) {
	var allowedUsersTgIds []int64
	db.Raw(`SELECT users.tg_user_id FROM users
		INNER JOIN user_roles ON users.id = user_roles.user_id
		INNER JOIN roles ON user_roles.role_id = roles.id
		WHERE roles.name = 'moder' OR roles.name = 'admin'`,
	).Scan(&allowedUsersTgIds)

	var msg tgbotapi.MessageConfig

	var response string
	if volume.New {
		response = fmt.Sprintf("На модерацию пришёл новый том\n\nНазвание: %s", volume.Name)
	} else {
		response = fmt.Sprintf("На модерацию пришли изменения для тома %s", volume.Name)
	}

	for i := 0; i < len(allowedUsersTgIds); i++ {
		msg = tgbotapi.NewMessage(allowedUsersTgIds[i], response)
		serviceBot.Send(msg)
	}

	return &pb.Empty{}, nil
}

func (s Server) NotifyAboutChapterOnModeration(ctx context.Context, chapter *pb.ChapterOnModeration) (*pb.Empty, error) {
	var allowedUsersTgIds []int64
	db.Raw(`SELECT users.tg_user_id FROM users
		INNER JOIN user_roles ON users.id = user_roles.user_id
		INNER JOIN roles ON user_roles.role_id = roles.id
		WHERE roles.name = 'moder' OR roles.name = 'admin'`,
	).Scan(&allowedUsersTgIds)

	var msg tgbotapi.MessageConfig

	var response string
	if chapter.New {
		response = fmt.Sprintf("На модерацию пришла новая глава\n\nНазвание: %s", chapter.Name)
	} else {
		response = fmt.Sprintf("На модерацию пришли изменения для главы %s", chapter.Name)
	}
	for i := 0; i < len(allowedUsersTgIds); i++ {
		msg = tgbotapi.NewMessage(allowedUsersTgIds[i], response)
		serviceBot.Send(msg)
	}

	return &pb.Empty{}, nil
}

func (s Server) NotifyAboutUserOnModeration(ctx context.Context, user *pb.User) (*pb.Empty, error) {
	var allowedUsersTgIds []int64
	db.Raw(`SELECT users.tg_user_id FROM users
		INNER JOIN user_roles ON users.id = user_roles.user_id
		INNER JOIN roles ON user_roles.role_id = roles.id
		WHERE roles.name = 'moder' OR roles.name = 'admin'`,
	).Scan(&allowedUsersTgIds)

	var msg tgbotapi.MessageConfig

	var response string
	if user.New {
		response = fmt.Sprintf("Верификации ожидает новый пользователь\n\nИмя: %s", user.Name)
	} else {
		response = fmt.Sprintf("Подтверждения изменений аккаунта ожидает пользователь %s", user.Name)
	}
	for i := 0; i < len(allowedUsersTgIds); i++ {
		msg = tgbotapi.NewMessage(allowedUsersTgIds[i], response)
		serviceBot.Send(msg)
	}

	return &pb.Empty{}, nil
}
