package rpc

import (
	"context"
	"log"
	"net"

	"github.com/AnatolyKoltun/calculator-storage/models"
	pb "github.com/AnatolyKoltun/calculator-storage/proto"
	"github.com/AnatolyKoltun/calculator-storage/repositories"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	pb.UnimplementedStorageServiceServer
	repositories.CalculationRepository
}

// ListCalculations — реализация gRPC метода для получения данных
func (s *server) ListCalculations(ctx context.Context, req *pb.ListRequest) (*pb.CalculationListResponse, error) {
	log.Printf("gRPC запрос: ID=%s", req)

	filter := models.FilterRequest{}
	calculations := &pb.CalculationListResponse{}

	filter.DateTo = req.DateTo
	filter.DateFrom = req.DateFrom

	values, err := s.GetList(ctx, filter)

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

func RunningGrpc() {
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
