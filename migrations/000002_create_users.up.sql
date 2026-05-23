CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,        -- Автоинкрементный ID
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'customer',
    created_at TIMESTAMP DEFAULT NOW() -- Дата создания по умолчанию
);