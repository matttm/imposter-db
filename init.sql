CREATE DATABASE IF NOT EXISTS packagemain;

CREATE TABLE IF NOT EXISTS packagemain.orders_v2 (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
) ENGINE=InnoDB;

INSERT INTO packagemain.orders_v2 (name) VALUES ('order1');
