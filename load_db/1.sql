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
