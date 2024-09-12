CREATE DATABASE IF NOT EXISTS `act-db`;
USE `act-db`;

CREATE TABLE IF NOT EXISTS `type_user` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `typeOfUser` VARCHAR(255) NOT NULL,
    `companyAllowed` INT NOT NULL
);

CREATE TABLE IF NOT EXISTS `companies` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `companyName` VARCHAR(255) NOT NULL,
    `userId` INT UNSIGNED NOT NULL,
    CONSTRAINT `fk_company_user` FOREIGN KEY (`userId`) REFERENCES `user` (`id`)
);

CREATE TABLE IF NOT EXISTS `users` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `firstName` VARCHAR(255) NOT NULL,
    `lastName` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `typeUserId` INT UNSIGNED NOT NULL,
    `createdAt` DATETIME NOT NULL,
    `updatedAt` DATETIME NOT NULL,
    `deletedAt` DATETIME,
    CONSTRAINT `fk_user_type_user` FOREIGN KEY (`typeUserId`) REFERENCES `type_user` (`id`)
);
