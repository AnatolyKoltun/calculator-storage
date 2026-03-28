package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/AnatolyKoltun/calculator-storage/config"
	"github.com/AnatolyKoltun/calculator-storage/database"
	"github.com/AnatolyKoltun/calculator-storage/proto"
	"github.com/AnatolyKoltun/calculator-storage/repositories"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

func connectToDB() {

	dsn := new(config.DataSourceName)
	dsn.GetDatabaseURL()

	database.Connect(dsn.DatabaseURL)
}

type server struct {
	proto.UnimplementedStorageServiceServer
	// Добавьте поля для БД и NATS
}

func (s *server) GetCalculation(ctx context.Context, req *proto.ListRequest) (*proto.CalculationListResponse, error) {
	calcRepository := repositories.CalculationRepository{}
	//repositories.GetList(ctx context.Context, req *proto.CalculationListResponse{})
	return calcRepository.GetList(ctx context.Context, req *proto.ListRequest)
}

func main() {

	defer database.Close()
	connectToDB()

	// 1. Подключение к NATS (consumer)
	nc, _ := nats.Connect(os.Getenv("NATS_URL"))
	js, _ := nc.JetStream()

	// Подписка на сообщения
	js.Subscribe("calculations.*", func(msg *nats.Msg) {
		// Обработка входящих данных
		msg.Ack()
	}, nats.Durable("storage-consumer"), nats.ManualAck())

	// 2. Запуск gRPC сервера
	lis, _ := net.Listen("tcp", ":50051")
	grpcServer := grpc.NewServer()
	proto.RegisterStorageServiceServer(grpcServer, &server{})

	log.Println("Storage service started on :50051")
	grpcServer.Serve(lis)
}
