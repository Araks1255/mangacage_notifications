package main

import (
	"net"

	mygrpc "github.com/Araks1255/mangacage_notifications/internal/grpc"
	"github.com/Araks1255/mangacage_notifications/pkg/common/db"
	pb "github.com/Araks1255/mangacage_service_protos"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
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

	mygrpc.InitBotAndDB(bot, db)

	lis, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterServiceNotificationsServer(grpcServer, mygrpc.Server{})
	grpcServer.Serve(lis)
}
