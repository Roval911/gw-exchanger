CREATE TABLE exchange_rates (
                                from_currency VARCHAR(3) NOT NULL,
                                to_currency VARCHAR(3) NOT NULL,
                                rate FLOAT NOT NULL,
                                PRIMARY KEY (from_currency, to_currency)
);
INSERT INTO exchange_rates (from_currency, to_currency, rate) VALUES
                                                                  ('USD', 'EUR', 0.94),
                                                                  ('USD', 'RUB', 85.5),
                                                                  ('EUR', 'USD', 1.06),
                                                                  ('EUR', 'RUB', 90.8),
                                                                  ('RUB', 'USD', 0.0117),
                                                                  ('RUB', 'EUR', 0.011),
                                                                  ('RUB', 'GBP', 0.0096)