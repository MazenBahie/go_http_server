CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(10) NOT NULL DEFAULT 'user'
);


CREATE TABLE credit_cards (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    cardNumber bigint NOT NULL,
    cvv VARCHAR(3) NOT NULL,
    expirationDate DATE NOT NULL,
    user_id INTEGER REFERENCES users(id)
);


CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    description TEXT
);
