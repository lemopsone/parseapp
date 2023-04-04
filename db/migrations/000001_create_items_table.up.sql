CREATE TABLE IF NOT EXISTS items(
    href VARCHAR(255) NOT NULL,
    title VARCHAR(255),
    current_price_yen VARCHAR(50),
    current_price_rur VARCHAR(50),
    blitz_price_yen VARCHAR(50),
    blitz_price_rur VARCHAR(50),
    time_left VARCHAR(50),
    telegram_id VARCHAR(255)
);