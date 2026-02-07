-- Создание таблиц
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    amount DECIMAL(15,2) NOT NULL,
    user_id UUID NOT NULL,
    category_id UUID NOT NULL REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT NOW(),

    UNIQUE(user_id, name)
);