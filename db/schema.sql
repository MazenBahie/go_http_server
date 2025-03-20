CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    stripe_customer_id VARCHAR(255) UNIQUE,
    role VARCHAR(10) NOT NULL DEFAULT 'user'
);


CREATE TABLE IF NOT EXISTS credit_cards (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    card_number bigint NOT NULL,
    cvv VARCHAR(3) NOT NULL,
    expiration_date DATE NOT NULL,
    user_id INTEGER REFERENCES users(id)
);


CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    description TEXT
);

CREATE TABLE IF NOT EXISTS user_orders (
    id SERIAL PRIMARY KEY,
    order_date DATE NOT NULL,
    user_id INTEGER REFERENCES users(id),
	product_id INTEGER REFERENCES products(id)
);
