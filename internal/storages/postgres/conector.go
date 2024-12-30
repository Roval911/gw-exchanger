package postgres

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"log"
)

var db *sql.DB

type ConnectionInfo struct {
	Host     string
	Port     int
	Username string
	DBName   string
	SSLMode  string
	Password string
}

type PostgresStorage struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewPostgresStorage(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{db: db}
}

func NewPostgresConnection(info ConnectionInfo) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
		info.Host, info.Port, info.Username, info.DBName, info.SSLMode, info.Password))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Закрытие соединения с базой данных
func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func RunMigrations() {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Error creating migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migration", // Путь к папке с миграциями
		"postgres",         // Имя базы данных
		driver,
	)
	if err != nil {
		log.Fatalf("Error initializing migration: %v", err)
	}

	// Применяем все миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Error applying migration: %v", err)
	}

	log.Println("Migrations applied successfully")
}

// Откат последней миграции
func RollbackLastMigration() {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Error creating migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migration",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Error initializing migration: %v", err)
	}

	// Откат последней миграции
	if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Error rolling back migration: %v", err)
	}

	log.Println("Last migration rolled back successfully")
}

func SetDB(mockDB *sql.DB) {
	db = mockDB
}
