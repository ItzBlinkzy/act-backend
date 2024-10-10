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
    type_user_id INTEGER NOT NULL,  
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,        
    CONSTRAINT fk_user_type_user FOREIGN KEY (type_user_id) REFERENCES type_user (id)
);

CREATE TABLE IF NOT EXISTS bought_stocks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    ticker VARCHAR(255) NOT NULL,
    quantity_owned INTEGER,
    quantity_sold INTEGER,
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,        
    CONSTRAINT fk_user_bought_stocks FOREIGN KEY (user_id) REFERENCES users (id)
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