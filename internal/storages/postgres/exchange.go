package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gw-exchanger/internal/storages"
)

func (p *PostgresStorage) GetAllRates(ctx context.Context) (*storages.ExchangeRates, error) {
	const op = "postgres.GetAllRates"
	query := "SELECT from_currency, rate FROM exchange_rates WHERE to_currency = 'RUB'"

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		p.logger.Printf("%s: ошибка при выполнении запроса: %v", op, err)
		return nil, fmt.Errorf("ошибка при получении курсов: %w", err)
	}
	defer rows.Close()

	rates := make(map[string]float32)
	for rows.Next() {
		var currency string
		var rate float32
		if err := rows.Scan(&currency, &rate); err != nil {
			p.logger.Printf("%s: ошибка при сканировании строки: %v", op, err)
			return nil, fmt.Errorf("ошибка при сканировании строки: %w", err)
		}
		rates[currency] = rate
	}

	if len(rates) == 0 {
		p.logger.Printf("%s: курсы валют не найдены", op)
		return nil, errors.New("курсы валют не найдены")
	}

	return &storages.ExchangeRates{Rates: rates}, nil
}

func (p *PostgresStorage) GetRateForCurrency(ctx context.Context, fromCurrency, toCurrency string) (*storages.ExchangeRate, error) {
	const op = "postgres.GetRateForCurrency"
	query := "SELECT rate FROM exchange_rates WHERE from_currency = $1 AND to_currency = $2"

	var rate float32
	err := p.db.QueryRowContext(ctx, query, fromCurrency, toCurrency).Scan(&rate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p.logger.Printf("%s: курс валюты не найден для %s -> %s", op, fromCurrency, toCurrency)
			return nil, fmt.Errorf("пара валют %s -> %s не найдена", fromCurrency, toCurrency)
		}
		p.logger.Printf("%s: ошибка при выполнении запроса: %v", op, err)
		return nil, fmt.Errorf("ошибка при получении курса: %w", err)
	}

	return &storages.ExchangeRate{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         rate,
	}, nil
}
