package tests

import (
	"context"
	"net"
	"os"
	"sync"
	"testing"

	"github.com/Araks1255/mangacage_notifications/internal/sender"
	"github.com/Araks1255/mangacage_notifications/internal/services/moderation_notifications"
	"github.com/Araks1255/mangacage_notifications/internal/services/site_notifications"
	"github.com/Araks1255/mangacage_notifications/pkg/common/db"
	mn "github.com/Araks1255/mangacage_protos/gen/moderation_notifications"
	sn "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

var env struct {
	DB                            *gorm.DB
	Bot                           *tgbotapi.BotAPI
	SiteNotificationsClient       sn.SiteNotificationsClient
	ModerationNotificationsClient mn.ModerationNotificationsClient
	TgUserID                      int64
	Ctx                           context.Context
}

func TestMain(m *testing.M) {
	os.Chdir("./..")

	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	dbUrl := viper.Get("TEST_DB_URL").(string)
	token := viper.Get("TOKEN").(string)
	testTgUserID := viper.GetInt64("TEST_TG_USER_ID")

	db, err := db.Init(dbUrl)
	if err != nil {
		panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	sender := sender.NewSender(bot)

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	env.DB = db
	env.Bot = bot
	env.SiteNotificationsClient = sn.NewSiteNotificationsClient(conn)
	env.ModerationNotificationsClient = mn.NewModerationNotificationsClient(conn)
	env.TgUserID = testTgUserID
	env.Ctx = context.Background()

	modersTgIDs, err := getModersTgIDs(env.DB)
	if err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	siteNotificationsServer := site_notifications.NewServer(db, &sync.RWMutex{}, sender, &modersTgIDs)
	moderationNotificationsServer := moderation_notifications.NewServer(db, sender)

	sn.RegisterSiteNotificationsServer(s, siteNotificationsServer)
	mn.RegisterModerationNotificationsServer(s, moderationNotificationsServer)

	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()

	code := m.Run()

	cleanTestDB(env.DB)

	os.Exit(code)
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
			r.name = 'moderator' OR r.name = 'admin'`,
	).Scan(&res).Error

	if err != nil {
		return nil, err
	}

	return res, nil
}

func cleanTestDB(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE authors, chapters, titles, users, teams, genres, tags RESTART IDENTITY CASCADE")
}
