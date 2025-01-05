package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gw-exchanger/internal/storages"
)

// GetAllRates - метод для получения всех курсов валют для обмена в рубли (RUB).
// Запрашивает данные о курсах валют из базы данных, где все курсы конвертируются в рубли.
func (p *PostgresStorage) GetAllRates(ctx context.Context) (*storages.ExchangeRates, error) {
	const op = "postgres.GetAllRates" // Константа для указания имени операции (используется для логирования).

	// SQL-запрос для получения всех валют, которые обмениваются на рубли.
	query := "SELECT from_currency, rate FROM exchange_rates WHERE to_currency = 'RUB'"

	// Выполнение запроса к базе данных.
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		// Логирование ошибки при выполнении запроса.
		p.logger.Printf("%s: ошибка при выполнении запроса: %v", op, err)
		// Возвращаем ошибку с описанием проблемы.
		return nil, fmt.Errorf("ошибка при получении курсов: %w", err)
	}
	defer rows.Close() // Обеспечиваем закрытие rows после выполнения функции.

	// Мапа для хранения курсов валют.
	rates := make(map[string]float32)

	// Итерируем по строкам результата запроса.
	for rows.Next() {
		var currency string
		var rate float32
		// Сканируем текущую строку в переменные.
		if err := rows.Scan(&currency, &rate); err != nil {
			// Логирование ошибки при сканировании строки.
			p.logger.Printf("%s: ошибка при сканировании строки: %v", op, err)
			// Возвращаем ошибку с описанием проблемы.
			return nil, fmt.Errorf("ошибка при сканировании строки: %w", err)
		}
		// Добавляем валюту и её курс в мапу.
		rates[currency] = rate
	}

	// Если курсы валют не найдены, логируем ошибку.
	if len(rates) == 0 {
		p.logger.Printf("%s: курсы валют не найдены", op)
		// Возвращаем ошибку, что курсы валют не найдены.
		return nil, errors.New("курсы валют не найдены")
	}

	// Возвращаем найденные курсы валют в структуре ExchangeRates.
	return &storages.ExchangeRates{Rates: rates}, nil
}

// GetRateForCurrency - метод для получения курса обмена для конкретной валютной пары.
// Принимает валюту-источник (fromCurrency) и валюту-цель (toCurrency) и возвращает курс обмена.
func (p *PostgresStorage) GetRateForCurrency(ctx context.Context, fromCurrency, toCurrency string) (*storages.ExchangeRate, error) {
	const op = "postgres.GetRateForCurrency" // Константа для указания имени операции (для логирования).

	// SQL-запрос для получения курса обмена между двумя валютами.
	query := "SELECT rate FROM exchange_rates WHERE from_currency = $1 AND to_currency = $2"

	var rate float32
	// Выполнение запроса и сканирование результата в переменную rate.
	err := p.db.QueryRowContext(ctx, query, fromCurrency, toCurrency).Scan(&rate)
	if err != nil {
		// Если ошибка связана с отсутствием данных (Нет строк в результате).
		if errors.Is(err, sql.ErrNoRows) {
			// Логирование ошибки, что валюта не найдена для данной валютной пары.
			p.logger.Printf("%s: курс валюты не найден для %s -> %s", op, fromCurrency, toCurrency)
			// Возвращаем ошибку, что пара валют не найдена.
			return nil, fmt.Errorf("пара валют %s -> %s не найдена", fromCurrency, toCurrency)
		}
		// Логирование других типов ошибок при выполнении запроса.
		p.logger.Printf("%s: ошибка при выполнении запроса: %v", op, err)
		// Возвращаем ошибку с описанием проблемы.
		return nil, fmt.Errorf("ошибка при получении курса: %w", err)
	}

	// Возвращаем данные о валютной паре и её курсе в нужном формате.
	return &storages.ExchangeRate{
		FromCurrency: fromCurrency,
		ToCurrency:   toCurrency,
		Rate:         rate,
	}, nil
}
