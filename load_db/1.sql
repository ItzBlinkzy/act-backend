CREATE DATABASE IF NOT EXISTS `c4u-go`;
USE `c4u-go`;

CREATE TABLE IF NOT EXISTS `contracts` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `nameOfContract` VARCHAR(255) NOT NULL,
    `gigaAllowed` INT NOT NULL,
    `numberUserAllowed` INT NOT NULL
);

CREATE TABLE IF NOT EXISTS `companies` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `name` VARCHAR(255) NOT NULL,
    `contractTypeId` INT UNSIGNED NOT NULL,
    `pivaCodiceFiscale` VARCHAR(255) NOT NULL,
    `spaceUsedMB` INT UNSIGNED NOT NULL,
    `intestazione` TEXT,
    `showIntestazione` BOOLEAN NOT NULL,
    `showLogo` BOOLEAN NOT NULL,
    `logoImg` VARCHAR(255),
    CONSTRAINT `fk_company_contract` FOREIGN KEY (`contractTypeId`) REFERENCES `contracts` (`id`)
);


CREATE TABLE IF NOT EXISTS `condomini` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `denominazione` VARCHAR(255) NOT NULL,
    `companyId` INT UNSIGNED NOT NULL,
    `idCodice` VARCHAR(255) NOT NULL,
    `codiceFiscale` VARCHAR(255) NOT NULL,
    `indirizzo` TEXT NOT NULL,
    `provincia` VARCHAR(100) NOT NULL,
    `comune` VARCHAR(100) NOT NULL,
    `email` VARCHAR(255),
    `codiceSDI` VARCHAR(100) NOT NULL,
    `contoCorrentePrincipaleIban` VARCHAR(255),
    `intestatarioContoPrincipale` VARCHAR(255),
    `contoCorrenteSecondarioIban` VARCHAR(255),
    `intestatarioContoSecondario` VARCHAR(255),
    `intestazione` TEXT,
    `note` TEXT,
    `created_at` DATETIME NOT NULL,
    `updated_at` DATETIME NOT NULL,
    `deleted_at` DATETIME,
    CONSTRAINT `fk_condominio_company` FOREIGN KEY (`companyId`) REFERENCES `companies` (`id`)
);

CREATE TABLE IF NOT EXISTS `users` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `firstName` VARCHAR(255) NOT NULL,
    `lastName` VARCHAR(255) NOT NULL,
    `email` VARCHAR(255) NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `companyId` INT UNSIGNED NOT NULL,
    `createdAt` DATETIME NOT NULL,
    `updatedAt` DATETIME NOT NULL,
    `deletedAt` DATETIME,
    CONSTRAINT `fk_user_company` FOREIGN KEY (`companyId`) REFERENCES `companies` (`id`)
);


CREATE TABLE IF NOT EXISTS `tabelle_millesimali` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `nomeTabella` VARCHAR(255) NOT NULL,
    `millesimiAssociati` DOUBLE PRECISION NOT NULL,
    `condominioId` INT UNSIGNED NOT NULL,
    `createdAt` DATETIME NOT NULL,
    `updatedAt` DATETIME NOT NULL,
    `deletedAt` DATETIME,
    CONSTRAINT `fk_tabella_millesimale_condominio`
        FOREIGN KEY (`condominioId`) 
        REFERENCES `condomini` (`id`)
);


CREATE TABLE IF NOT EXISTS `unita_immobiliari` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `condominioId` INT UNSIGNED NOT NULL,
    `tipoLocale` VARCHAR(255) NOT NULL,
    `piano` VARCHAR(255),
    `sezione` VARCHAR(255),
    `particella` VARCHAR(255),
    `scala` VARCHAR(255),
    `interno` VARCHAR(255),
    `rendita` DOUBLE PRECISION,
    `foglio` VARCHAR(255),
    `sub` VARCHAR(255),
    `note` TEXT,
    `createdAt` DATETIME NOT NULL,
    `updatedAt` DATETIME NOT NULL,
    `deletedAt` DATETIME,
    CONSTRAINT `fk_unita_immobiliari_condominio`
        FOREIGN KEY (`condominioId`) 
        REFERENCES `condomini` (`id`)
);


CREATE TABLE IF NOT EXISTS `unita_immobiliari_tab_mill` (
    `id` INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    `unitaImmobiliareId` INT UNSIGNED NOT NULL, 
    `tabellaMillesimaleId` INT UNSIGNED NOT NULL,
    `millesimiAssociati` DOUBLE PRECISION NOT NULL,
    `note` TEXT,
    `createdAt` DATETIME NOT NULL,
    `updatedAt` DATETIME NOT NULL,
    `deletedAt` DATETIME,
    CONSTRAINT `fk_unita_immobiliari_tab_mill_unita`
        FOREIGN KEY (`unitaImmobiliareId`) 
        REFERENCES `unita_immobiliari` (`id`),
    CONSTRAINT `fk_unita_immobiliari_tab_mill_tabella`
        FOREIGN KEY (`tabellaMillesimaleId`) 
        REFERENCES `tabelle_millesimali` (`id`)
);
