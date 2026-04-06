package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"os"

	"github.com/AnatolyKoltun/calculator-storage/config"
	"github.com/AnatolyKoltun/calculator-storage/database"
	"github.com/AnatolyKoltun/calculator-storage/models"
	pb "github.com/AnatolyKoltun/calculator-storage/proto"
	"github.com/AnatolyKoltun/calculator-storage/repositories"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var calcRepository = repositories.CalculationRepository{}

type server struct {
	pb.UnimplementedStorageServiceServer
}

// ListCalculations — реализация gRPC метода для получения данных
func (s *server) ListCalculations(ctx context.Context, req *pb.ListRequest) (*pb.CalculationListResponse, error) {
	log.Printf("gRPC запрос: ID=%s", req)

	filter := models.FilterRequest{}
	calculations := &pb.CalculationListResponse{}

	filter.DateTo = req.DateTo
	filter.DateFrom = req.DateFrom

	values, err := calcRepository.GetList(ctx, filter)

	if err == nil {
		calculations.Calculations = make([]*pb.Calculation, len(values))

		for ind, value := range values {
			calculations.Calculations[ind] = &pb.Calculation{
				Id:        int32(value.ID),
				Argument1: value.Argument1,
				Argument2: value.Argument2,
				Operator:  value.Operator,
				Result:    value.Result,
				CreatedAt: timestamppb.New(value.CreatedAt),
			}
		}
	}

	return calculations, err
}

func connectToDB() {

	dsn := new(config.DataSourceName)
	dsn.GetDatabaseURL()

	database.Connect(dsn.DatabaseURL)
}

func main() {
	defer database.Close()

	connectToDB()

	// 1. Подключение к NATS
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatal("Ошибка подключения к NATS:", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()

	if err != nil {
		log.Fatal("Ошибка JetStream:", err)
	}

	// 2. Создаем/проверяем Stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "CALCULATIONS",
		Subjects: []string{"calculations.*"},
		Storage:  nats.FileStorage,
	})

	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		log.Fatal("Ошибка создания Stream:", err)
	}

	// 3. Подписываемся на сообщения (durable consumer)
	_, err = js.Subscribe("calculations.*", func(msg *nats.Msg) {
		log.Printf("Получено сообщение: %s", string(msg.Data))

		// Парсим JSON
		var data *models.Calculation
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("Ошибка парсинга: %v", err)
			msg.Nak() // Не подтверждаем, NATS отправит повторно
			return
		}

		ctx := context.Background()
		err = calcRepository.Save(ctx, data)

		if err == nil {
			log.Printf("Успешно сохранено: %+v", data)
			msg.Ack()
		} else {
			log.Printf("Ошибка сохранения: %v", err)
			msg.Nak()
		}
	}, nats.Durable("storage-consumer"), nats.ManualAck())

	if err != nil {
		log.Fatal("Ошибка подписки:", err)
	}

	log.Println("NATS consumer запущен")

	// 4. Запуск gRPC сервера
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Ошибка создания listener:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStorageServiceServer(grpcServer, &server{})

	log.Println("gRPC сервер запущен на :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Ошибка gRPC сервера:", err)
	}
}
