CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255),
    password VARCHAR(255)
);

INSERT INTO users (name, email, password) VALUES
    ('Mike', 'mike@example.com', 'secret'),
    ('John', 'john@example.com', 'secret'),
    ('Jane', 'jane@example.com', 'secret');