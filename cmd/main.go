package main

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/Araks1255/mangacage_notifications/internal/sender"
	"github.com/Araks1255/mangacage_notifications/internal/services/moderation_notifications"
	"github.com/Araks1255/mangacage_notifications/internal/services/site_notifications"
	"github.com/Araks1255/mangacage_notifications/pkg/common/db"
	mn "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	sn "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	dbUrl := viper.Get("DB_URL").(string)
	token := viper.Get("TOKEN").(string)

	db, err := db.Init(dbUrl)
	if err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	sender := sender.NewSender(bot)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	modersTgIDs, err := getModersTgIDs(db)
	if err != nil {
		panic(err)
	}

	var mu sync.RWMutex

	go startUpdatingModersTgIDs(db, &mu, &modersTgIDs)
	go sender.Start()

	s := grpc.NewServer()

	siteNotificationsServer := site_notifications.NewServer(db, &mu, sender, &modersTgIDs)
	moderationNotificationsServer := moderation_notifications.NewServer(db, sender)

	sn.RegisterSiteNotificationsServer(s, siteNotificationsServer)
	mn.RegisterModerationNotificationsServer(s, moderationNotificationsServer)

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

func getModersTgIDs(db *gorm.DB) ([]int64, error) {
	var res []int64

	err := db.Raw(
		`SELECT
			u.tg_user_id
		FROM
			users AS u
			INNER JOIN user_roles AS ur ON ur.user_id = u.id
			INNER JOIN roles AS r ON r.id = ur.role_id
		WHERE
			r.name = 'moderator' OR r.name = 'admin' AND u.tg_user_id IS NOT NULL`,
	).Scan(&res).Error

	if err != nil {
		return nil, err
	}

	return res, nil
}

func startUpdatingModersTgIDs(db *gorm.DB, mu *sync.RWMutex, modersTgIDs *[]int64) {
	for {
		time.Sleep(24 * time.Hour)
		mu.Lock()
		if err := db.Raw(
			`SELECT
				u.tg_user_id
			FROM
				users AS u
				INNER JOIN user_roles AS ur ON ur.user_id = u.id
				INNER JOIN roles AS r ON r.id = ur.role_id
			WHERE
				r.name = 'moderator' OR r.name = 'admin' AND u.tg_user_id IS NOT NULL`,
		).Scan(modersTgIDs).Error; err != nil {
			log.Println(err)
		}
		mu.Unlock()
	}
}
