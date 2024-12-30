package handlers

import (
	"context"
	"fmt"
	pb "github.com/Roval911/proto-exchange/exchange"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetExchangeRates - обработчик для получения всех курсов валют
func (s *Server) GetExchangeRates(ctx context.Context, req *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	exchange, err := s.storage.GetAllRates(ctx)
	if err != nil {
		s.logger.Printf("Ошибка при получении курсов валют: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get exchange rates: %v", err))
	}

	return &pb.ExchangeRatesResponse{Rates: exchange.Rates}, nil
}

// GetExchangeRateForCurrency - обработчик для получения курса валюты по заданной паре
func (s *Server) GetExchangeRateForCurrency(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	if req.FromCurrency == "" || req.ToCurrency == "" {
		s.logger.Printf("Ошибка: пустое значение валют в запросе: %s -> %s", req.FromCurrency, req.ToCurrency)
		return nil, status.Error(codes.InvalidArgument, "currency fields must not be empty")
	}

	exchangeRate, err := s.storage.GetRateForCurrency(ctx, req.FromCurrency, req.ToCurrency)
	if err != nil {
		s.logger.Printf("Ошибка при получении курса валют %s -> %s: %v", req.FromCurrency, req.ToCurrency, err)
		return nil, status.Error(codes.NotFound, fmt.Sprintf("currency pair %s -> %s not supported", req.FromCurrency, req.ToCurrency))
	}

	return &pb.ExchangeRateResponse{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         exchangeRate.Rate,
	}, nil
}
