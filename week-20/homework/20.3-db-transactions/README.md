# ДЗ 20.3: Транзакции в PostgreSQL

## Цель

Научиться работать с транзакциями в PostgreSQL. Освоить Begin, Commit, Rollback и понять, когда и зачем нужны транзакции.

## Подготовка

1. Убедись, что PostgreSQL запущен:

```bash
cd week-20
docker-compose up -d
```

2. Проверь, что таблицы accounts и transactions существуют:

```bash
docker exec -it go-course-postgres psql -U gouser -d gocourse -c "SELECT * FROM accounts;"
```

## SQL-схема

```sql
-- Счета пользователей
CREATE TABLE accounts (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0 CHECK (balance >= 0),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- История транзакций
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    from_account_id INT REFERENCES accounts(id),
    to_account_id INT REFERENCES accounts(id),
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## Зачем нужны транзакции?

Транзакции обеспечивают ACID-свойства:

- **Atomicity (Атомарность)** — все операции выполняются как единое целое
- **Consistency (Согласованность)** — база переходит из одного корректного состояния в другое
- **Isolation (Изолированность)** — параллельные транзакции не влияют друг на друга
- **Durability (Долговечность)** — после Commit данные сохранены надёжно

**Пример:** При переводе денег нужно списать с одного счёта и зачислить на другой. Если что-то пойдёт не так между этими операциями, деньги не должны потеряться.

## Что нужно сделать

### Структуры данных

```go
type Account struct {
    ID        int
    UserID    int
    Balance   float64
    Currency  string
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Transaction struct {
    ID            int
    FromAccountID int
    ToAccountID   int
    Amount        float64
    Description   string
    CreatedAt     time.Time
}
```

### 1. GetAccountByID — получение счёта

```go
func GetAccountByID(ctx context.Context, pool *pgxpool.Pool, id int) (Account, error)
```

- Возвращает информацию о счёте по ID

### 2. TransferMoney — перевод денег (основное задание)

```go
func TransferMoney(ctx context.Context, pool *pgxpool.Pool, fromID, toID int, amount float64, description string) error
```

**Алгоритм:**
1. Начать транзакцию: `tx, err := pool.Begin(ctx)`
2. Добавить `defer tx.Rollback(ctx)` — безопасный откат при любой ошибке
3. Списать деньги с исходного счёта (проверить, что баланс достаточен)
4. Зачислить деньги на целевой счёт
5. Записать транзакцию в таблицу transactions
6. Зафиксировать: `tx.Commit(ctx)`

**Ошибки:**
- `ErrAccountNotFound` — счёт не найден
- `ErrInsufficientFunds` — недостаточно средств
- `ErrSameAccount` — перевод на тот же счёт

### 3. GetAccountTransactions — история транзакций

```go
func GetAccountTransactions(ctx context.Context, pool *pgxpool.Pool, accountID int) ([]Transaction, error)
```

- Возвращает все транзакции, где счёт является отправителем или получателем
- Сортировка по дате (новые первые)

### 4. WithTransaction — обёртка для транзакций (дополнительно)

```go
func WithTransaction(ctx context.Context, pool *pgxpool.Pool, fn func(tx pgx.Tx) error) error
```

- Универсальная функция для выполнения кода в транзакции
- Автоматический Rollback при ошибке
- Автоматический Commit при успехе

## SQL запросы

```sql
-- Получить счёт
SELECT id, user_id, balance, currency, created_at, updated_at
FROM accounts WHERE id = $1;

-- Списать деньги (с проверкой баланса)
UPDATE accounts SET balance = balance - $1, updated_at = NOW()
WHERE id = $2 AND balance >= $1;

-- Зачислить деньги
UPDATE accounts SET balance = balance + $1, updated_at = NOW()
WHERE id = $2;

-- Записать транзакцию
INSERT INTO transactions (from_account_id, to_account_id, amount, description)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at;

-- История транзакций счёта
SELECT id, from_account_id, to_account_id, amount, description, created_at
FROM transactions
WHERE from_account_id = $1 OR to_account_id = $1
ORDER BY created_at DESC;
```

## Пример использования

```go
func main() {
    ctx := context.Background()
    pool, _ := pgxpool.New(ctx, dsn)
    defer pool.Close()

    // Проверяем балансы до перевода
    acc1, _ := GetAccountByID(ctx, pool, 1)
    acc2, _ := GetAccountByID(ctx, pool, 2)
    fmt.Printf("Before: Account 1: %.2f, Account 2: %.2f\n", acc1.Balance, acc2.Balance)

    // Переводим деньги
    err := TransferMoney(ctx, pool, 1, 2, 100.00, "Payment for services")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Transfer successful!")

    // Проверяем балансы после перевода
    acc1, _ = GetAccountByID(ctx, pool, 1)
    acc2, _ = GetAccountByID(ctx, pool, 2)
    fmt.Printf("After: Account 1: %.2f, Account 2: %.2f\n", acc1.Balance, acc2.Balance)

    // История транзакций
    transactions, _ := GetAccountTransactions(ctx, pool, 1)
    fmt.Printf("Transactions for account 1: %d\n", len(transactions))
}
```

Ожидаемый вывод:
```
=== PostgreSQL Transactions Demo ===

--- Initial Balances ---
Account 1: 1000.00 USD
Account 2: 500.00 USD

--- Transfer Money ---
Transferring 100.00 from account 1 to account 2...
Transfer successful!

--- Balances After Transfer ---
Account 1: 900.00 USD
Account 2: 600.00 USD

--- Transaction History ---
Account 1 transactions:
  - Sent 100.00 to account 2: Payment for services

--- Test Insufficient Funds ---
Trying to transfer 10000.00 (more than balance)...
Error: insufficient funds (expected)

--- Test Same Account ---
Trying to transfer to same account...
Error: cannot transfer to same account (expected)
```

## Подсказки

### Базовый паттерн транзакции

```go
func TransferMoney(ctx context.Context, pool *pgxpool.Pool, fromID, toID int, amount float64, description string) error {
    // Проверка: нельзя переводить на тот же счёт
    if fromID == toID {
        return ErrSameAccount
    }

    // Начинаем транзакцию
    tx, err := pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    // defer Rollback безопасен — если был Commit, Rollback ничего не сделает
    defer tx.Rollback(ctx)

    // Операции внутри транзакции используют tx, а не pool
    result, err := tx.Exec(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1", amount, fromID)
    if err != nil {
        return fmt.Errorf("debit: %w", err)
    }
    if result.RowsAffected() == 0 {
        return ErrInsufficientFunds // или ErrAccountNotFound
    }

    // ... остальные операции

    // Фиксируем транзакцию
    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("commit: %w", err)
    }

    return nil
}
```

### Использование обёртки

```go
func WithTransaction(ctx context.Context, pool *pgxpool.Pool, fn func(tx pgx.Tx) error) error {
    tx, err := pool.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    if err := fn(tx); err != nil {
        return err
    }

    return tx.Commit(ctx)
}

// Использование
err := WithTransaction(ctx, pool, func(tx pgx.Tx) error {
    // все операции здесь
    return nil
})
```

## Критерии выполнения

- [ ] GetAccountByID возвращает информацию о счёте
- [ ] TransferMoney выполняет перевод атомарно (либо всё, либо ничего)
- [ ] TransferMoney проверяет достаточность средств
- [ ] TransferMoney записывает транзакцию в таблицу transactions
- [ ] TransferMoney возвращает ErrInsufficientFunds при нехватке средств
- [ ] TransferMoney возвращает ErrSameAccount при переводе на тот же счёт
- [ ] GetAccountTransactions возвращает историю транзакций
- [ ] При ошибке между операциями изменения откатываются (Rollback)
- [ ] Код компилируется без ошибок
