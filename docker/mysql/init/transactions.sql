CREATE TABLE IF NOT EXISTS clients (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS currencies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    iso VARCHAR(255) NOT NULL
);

INSERT INTO currencies (iso) VALUES
    ('BTC'),
    ('ETH'),
    ('USDT');

INSERT INTO clients (name) VALUES
    ('Monsters Inc.'),
    ('A1A Car Wash'),
    ('Atlantis Casions');

CREATE TABLE IF NOT EXISTS transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    txid VARCHAR(255),
    amount INT NOT NULL,
    currency_id INT,
    client_id INT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_currency FOREIGN KEY (currency_id) REFERENCES currencies(id) ON DELETE CASCADE,
    CONSTRAINT fk_client FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS events (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    txid VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);