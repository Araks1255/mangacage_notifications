package moderation_notifications

import (
	"github.com/Araks1255/mangacage_notifications/internal/sender"
	pb "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedModerationNotificationsServer
	DB     *gorm.DB
	Sender *sender.Sender
}

func NewServer(db *gorm.DB, sender *sender.Sender) server {
	return server{
		DB:     db,
		Sender: sender,
	}
}
