package handlers

import (
	"context"
	"fmt"
	pb "github.com/Roval911/proto-exchange/exchange"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetExchangeRates - обработчик для получения всех курсов валют.
// Этот метод обрабатывает запросы на получение всех доступных курсов валют и возвращает их в ответе.
func (s *Server) GetExchangeRates(ctx context.Context, req *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	// Получаем все курсы валют из хранилища (например, из базы данных или кэшированного источника).
	exchange, err := s.storage.GetAllRates(ctx)
	if err != nil {
		// Если при получении курсов возникла ошибка, логируем её и возвращаем ошибку через gRPC.
		// Используем код ошибки Internal, чтобы указать на внутреннюю ошибку на сервере.
		s.logger.Printf("Ошибка при получении курсов валют: %v", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get exchange rates: %v", err))
	}

	// Если все прошло успешно, формируем ответ с курсами валют и возвращаем его клиенту.
	return &pb.ExchangeRatesResponse{Rates: exchange.Rates}, nil
}

// GetExchangeRateForCurrency - обработчик для получения курса валюты по заданной паре.
// Этот метод обрабатывает запросы для получения курса обмена между двумя валютами (например, USD -> EUR).
func (s *Server) GetExchangeRateForCurrency(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	// Проверяем, что поля валют в запросе не пустые.
	// Если одно из полей пустое, логируем ошибку и возвращаем статус ошибки InvalidArgument.
	if req.FromCurrency == "" || req.ToCurrency == "" {
		s.logger.Printf("Ошибка: пустое значение валют в запросе: %s -> %s", req.FromCurrency, req.ToCurrency)
		return nil, status.Error(codes.InvalidArgument, "currency fields must not be empty")
	}

	// Получаем курс обмена для указанной валютной пары из хранилища.
	exchangeRate, err := s.storage.GetRateForCurrency(ctx, req.FromCurrency, req.ToCurrency)
	if err != nil {
		// Если при получении курса валют произошла ошибка, логируем её и возвращаем ошибку NotFound,
		// указывая, что такая валютная пара не поддерживается.
		s.logger.Printf("Ошибка при получении курса валют %s -> %s: %v", req.FromCurrency, req.ToCurrency, err)
		return nil, status.Error(codes.NotFound, fmt.Sprintf("currency pair %s -> %s not supported", req.FromCurrency, req.ToCurrency))
	}

	// Если курс найден, формируем ответ с деталями валютной пары и возвращаем его.
	return &pb.ExchangeRateResponse{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         exchangeRate.Rate, // Указываем курс обмена.
	}, nil
}
