CREATE DATABASE IF NOT EXISTS `act-db`;
USE `act-db`;

CREATE TABLE IF NOT EXISTS type_user (
    id SERIAL PRIMARY KEY,
    type_of_user VARCHAR(255) NOT NULL 
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    credit DECIMAL(10, 2) DEFAULT 0.00;
    type_user_id INTEGER NOT NULL,  
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,        
    login_method VARCHAR(255) NOT NUL
    CONSTRAINT fk_user_type_user FOREIGN KEY (type_user_id) REFERENCES type_user (id)
);

CREATE TABLE IF NOT EXISTS bought_stocks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    ticker VARCHAR(255) NOT NULL,
    quantity_owned INTEGER,
    quantity_sold INTEGER,
    client_id INTEGER,
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,        
    CONSTRAINT fk_user_bought_stocks FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_client_bought_stocks FOREIGN KEY (client_id) REFERENCES clients (id)
);

CREATE TABLE IF NOT EXISTS logs_bought_stocks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    bought_stock_id INTEGER NOT NULL,
    quantity_bought INTEGER DEFAULT NULL,
    quantity_sold INTEGER DEFAULT NULL,
    client_id INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_bought_stocks FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_logs_bought_stocks FOREIGN KEY (bought_stock_id) REFERENCES bought_stocks (id),
    CONSTRAINT fk_logs_clients_stocks FOREIGN KEY (client_id) REFERENCES clients (id)
);

CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    description TEXT NOT NULL,
    stars INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,   
    CONSTRAINT fk_review_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS client_manager_association (
    id SERIAL PRIMARY KEY,
    manager_id INTEGER NOT NULL,
    client_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,   
    CONSTRAINT fk_client_user FOREIGN KEY (client_id) REFERENCES clients (id),
    CONSTRAINT fk_manager_user FOREIGN KEY (manager_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS clients (
    id SERIAL PRIMARY KEY,
    company_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);
