CREATE TABLE IF NOT EXISTS clients (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

INSERT INTO clients (name) VALUES
    ('Monsters Inc.'),
    ('A1A Car Wash'),
    ('Atlantis Casions');

CREATE TABLE IF NOT EXISTS currencies (
    id INT AUTO_INCREMENT PRIMARY KEY,
    iso VARCHAR(255) NOT NULL UNIQUE
);

INSERT INTO currencies (iso) VALUES
    ('BTC'),
    ('ETH'),
    ('USDT');

CREATE TABLE IF NOT EXISTS addresses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    hash VARCHAR(255) NOT NULL,
    client_id INT NOT NULL,
    currency_id INT NOT NULL,
    CONSTRAINT fk_client FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE,
    CONSTRAINT fk_currency FOREIGN KEY (currency_id) REFERENCES currencies(id) ON DELETE CASCADE,
    UNIQUE KEY unique_currency_id_hash (currency_id, hash)
);

INSERT INTO addresses (hash, client_id, currency_id) VALUES
    ('2N3fEj5ZunRpuu4rFAydd27WEPopxCFKY5x',1,1),
    ('2MxFcfjTCf76QuC7wpVvdt8DzMggtTmrnYa',1,1),
    ('2MtzBo4LPoQEKooC81Pmnyu8eoQ7YMRRFPj',1,2);

CREATE TABLE IF NOT EXISTS transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    address_id INT,
    txid VARCHAR(255),
    amount INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_address FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE CASCADE,
    UNIQUE KEY unique_address_id_txid (address_id, txid)
);

CREATE TABLE IF NOT EXISTS messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    type VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    currency VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    txid VARCHAR(255) NOT NULL,
    amount INT NOT NULL,
    timestamp DATETIME NOT NULL,
    cretaed_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
