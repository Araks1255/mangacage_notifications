package moderation_notifications

import (
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedModerationNotificationsServer
	DB  *gorm.DB
	Bot *tgbotapi.BotAPI
}

func NewServer(db *gorm.DB, bot *tgbotapi.BotAPI) server {
	return server{
		DB:  db,
		Bot: bot,
	}
}
