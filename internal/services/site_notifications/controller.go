package site_notifications

import (
	"sync"

	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedSiteNotificationsServer
	DB          *gorm.DB
	Mu          *sync.RWMutex
	Bot         *tgbotapi.BotAPI
	ModersTgIDs *[]int64
}

func NewServer(db *gorm.DB, mu *sync.RWMutex, bot *tgbotapi.BotAPI, modersTgIDs *[]int64) server {
	return server{
		DB:          db,
		Mu:          mu,
		Bot:         bot,
		ModersTgIDs: modersTgIDs,
	}
}
