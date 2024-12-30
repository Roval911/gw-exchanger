package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gw-exchanger/internal/config"
	"gw-exchanger/internal/handlers"
	"gw-exchanger/internal/storages/postgres"
	"gw-exchanger/pkg/logger"
	"log"
	"net"
)

func main() {
	// Инициализация конфигурации
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Ошибка при загрузке конфигурации: %v", err)
	}

	// Инициализация логгера
	log := logger.InitLogger()

	// Подключение к базе данных PostgreSQL
	connInfo := postgres.ConnectionInfo{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		DBName:   cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.DB.Password,
	}

	db, err := postgres.NewPostgresConnection(connInfo)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	defer db.Close()

	// Инициализация хранилища PostgreSQL
	storage := postgres.NewPostgresStorage(db)

	// Запуск миграций
	//postgres.RunMigrations()

	// Создание gRPC сервера
	grpcServer := grpc.NewServer()

	// Регистрируем сервисы
	handlers.Register(grpcServer, log, storage)

	reflection.Register(grpcServer)

	// Открываем порт для сервера
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

	// Запуск gRPC сервера
	log.Infof("Сервер запущен на порту %d", cfg.Server.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка при запуске gRPC сервера: %v", err)
	}
}
