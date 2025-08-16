package site_notifications

import (
	"sync"

	"github.com/Araks1255/mangacage_notifications/internal/sender"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedSiteNotificationsServer
	DB          *gorm.DB
	mu          *sync.RWMutex
	Sender      *sender.Sender
	ModersTgIDs *[]int64
}

func NewServer(db *gorm.DB, mu *sync.RWMutex, sender *sender.Sender, modersTgIDs *[]int64) server {
	return server{
		DB:          db,
		mu:          mu,
		Sender:      sender,
		ModersTgIDs: modersTgIDs,
	}
}
