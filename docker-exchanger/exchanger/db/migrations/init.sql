CREATE TABLE exchange_rates (
    id SERIAL PRIMARY KEY,
    from_currency VARCHAR(3) NOT NULL,
    to_currency VARCHAR(3) NOT NULL,
    rate NUMERIC(10, 5) NOT NULL
);

INSERT INTO exchange_rates (from_currency, to_currency, rate)
VALUES
('USD', 'EUR', 0.85),
('EUR', 'USD', 1.18),
('USD', 'RUB', 70.00),
('RUB', 'USD', 0.014),
('EUR', 'RUB', 80.00),
('RUB', 'EUR', 0.0125);