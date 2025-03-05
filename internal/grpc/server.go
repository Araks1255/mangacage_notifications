package grpc

import (
	pb "github.com/Araks1255/mangacage_protos"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type Server struct {
	pb.UnimplementedNotificationsServer
}

var serviceBot *tgbotapi.BotAPI
var usersBot *tgbotapi.BotAPI
var db *gorm.DB

func InitBotsAndDB(existingServiceBot *tgbotapi.BotAPI, existingUsersBot *tgbotapi.BotAPI, existingDB *gorm.DB) {
	serviceBot = existingServiceBot
	usersBot = existingUsersBot
	db = existingDB
}
