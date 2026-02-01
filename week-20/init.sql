-- Инициализация базы данных для курса Go
-- Этот файл выполняется автоматически при первом запуске контейнера

-- Таблица пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    age INT CHECK (age >= 0 AND age <= 150),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Таблица счетов (для примера с транзакциями)
CREATE TABLE IF NOT EXISTS accounts (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0 CHECK (balance >= 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Таблица транзакций (для истории переводов)
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    from_account_id INT REFERENCES accounts(id),
    to_account_id INT REFERENCES accounts(id),
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для оптимизации
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_from_account ON transactions(from_account_id);
CREATE INDEX IF NOT EXISTS idx_transactions_to_account ON transactions(to_account_id);

-- Функция для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Триггеры для автоматического обновления updated_at
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_accounts_updated_at ON accounts;
CREATE TRIGGER update_accounts_updated_at
    BEFORE UPDATE ON accounts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Тестовые данные
INSERT INTO users (name, email, age) VALUES
    ('Alice Johnson', 'alice@example.com', 25),
    ('Bob Smith', 'bob@example.com', 30),
    ('Charlie Brown', 'charlie@example.com', 35)
ON CONFLICT (email) DO NOTHING;

-- Создаём счета для тестовых пользователей
INSERT INTO accounts (user_id, balance, currency)
SELECT id, 1000.00, 'USD' FROM users WHERE email = 'alice@example.com'
ON CONFLICT DO NOTHING;

INSERT INTO accounts (user_id, balance, currency)
SELECT id, 500.00, 'USD' FROM users WHERE email = 'bob@example.com'
ON CONFLICT DO NOTHING;

INSERT INTO accounts (user_id, balance, currency)
SELECT id, 2000.00, 'USD' FROM users WHERE email = 'charlie@example.com'
ON CONFLICT DO NOTHING;
