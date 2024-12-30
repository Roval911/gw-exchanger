package handlers

import (
	pb "github.com/Roval911/proto-exchange/exchange"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gw-exchanger/internal/storages"
)

type Server struct {
	pb.UnimplementedExchangeServiceServer
	storage storages.Storages
	logger  *logrus.Logger
}

func Register(gRPC *grpc.Server, logger *logrus.Logger, storage storages.Storages) {
	pb.RegisterExchangeServiceServer(gRPC, &Server{
		logger:  logger,
		storage: storage,
	})
}
