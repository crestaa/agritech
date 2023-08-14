GRANT ALL PRIVILEGES ON agritech_db.* TO 'user'@'%';
FLUSH PRIVILEGES;


CREATE DATABASE IF NOT EXISTS agritech_db;
USE agritech_db;

CREATE TABLE IF NOT EXISTS Campi (
    id_campo INT PRIMARY KEY AUTO_INCREMENT,
    nome_campo VARCHAR(255),
    latitudine DECIMAL(10, 8),
    longitudine DECIMAL(11, 8)
);

CREATE TABLE IF NOT EXISTS Sensori (
    id_sensore INT PRIMARY KEY AUTO_INCREMENT,
    mac VARCHAR(255),
    id_campo INT,
    latitudine DECIMAL(10, 8),
    longitudine DECIMAL(11, 8),
    FOREIGN KEY (id_campo) REFERENCES Campi(id_campo)
);

CREATE TABLE IF NOT EXISTS Tipi_Misurazione (
    id_tipo_misurazione INT PRIMARY KEY AUTO_INCREMENT,
    nome VARCHAR(32),
    unita_misura VARCHAR(8)
);

CREATE TABLE IF NOT EXISTS Misurazioni (
    id_misurazione INT PRIMARY KEY AUTO_INCREMENT,
    id_sensore INT,
    tipo_misurazione INT,
    valore DECIMAL(10, 2),
    data_ora DATETIME DEFAULT CURRENT_TIMESTAMP(),
    FOREIGN KEY (id_sensore) REFERENCES Sensori(id_sensore),
    FOREIGN KEY (tipo_misurazione) REFERENCES Tipi_Misurazione(id_tipo_misurazione)
);

