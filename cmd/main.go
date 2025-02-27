package main

import (
	"net"

	mygrpc "github.com/Araks1255/mangacage_notifications/internal/grpc"
	"github.com/Araks1255/mangacage_notifications/pkg/common/db"
	pb "github.com/Araks1255/mangacage_protos"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	dbUrl := viper.Get("DB_URL").(string)
	serviceToken := viper.Get("SERVICE_TOKEN").(string)
	usersToken := viper.Get("USERS_TOKEN").(string)

	db, err := db.Init(dbUrl)
	if err != nil {
		panic(err)
	}

	serviceBot, err := tgbotapi.NewBotAPI(serviceToken)
	if err != nil {
		panic(err)
	}

	usersBot, err := tgbotapi.NewBotAPI(usersToken)
	if err != nil {
		panic(err)
	}

	mygrpc.InitBotsAndDB(serviceBot, usersBot, db)

	lis, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationsServer(grpcServer, mygrpc.Server{})
	grpcServer.Serve(lis)
}
