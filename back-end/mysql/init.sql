CREATE DATABASE IF NOT EXISTS `Rates`;
Go
use `Rates`;
Go
CREATE TABLE  Rates(
    id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    date VARCHAR(255) NOT NULL,
    currency VARCHAR(255) NOT NULL, 
    rate float NOT NULL
)