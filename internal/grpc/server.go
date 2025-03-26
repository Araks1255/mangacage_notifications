package grpc

import (
	pb "github.com/Araks1255/mangacage_protos"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type Server struct {
	pb.UnimplementedNotificationsServer
	DB         *gorm.DB
	UsersBot   *tgbotapi.BotAPI
	ServiceBot *tgbotapi.BotAPI
}

func InitServer(existingServiceBot *tgbotapi.BotAPI, existingUsersBot *tgbotapi.BotAPI, existingDB *gorm.DB) Server {
	return Server{
		DB:         existingDB,
		UsersBot:   existingUsersBot,
		ServiceBot: existingServiceBot,
	}
}
