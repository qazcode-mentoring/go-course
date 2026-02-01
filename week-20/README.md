# Неделя 20: Работа с PostgreSQL в Go

## Теория

В этом модуле мы научимся работать с реляционными базами данных из Go. Мы будем использовать **PostgreSQL** — одну из самых популярных и мощных СУБД, и библиотеку **pgx** — современный драйвер PostgreSQL для Go.

### Почему pgx, а не database/sql?

Go предоставляет стандартный пакет `database/sql` для работы с SQL базами данных. Он работает через драйверы и предоставляет универсальный интерфейс. Однако для PostgreSQL рекомендуется использовать **pgx/v5**:

| Критерий | database/sql | pgx |
|----------|-------------|-----|
| Производительность | Хорошая | Лучше (нативный протокол) |
| PostgreSQL-специфичные типы | Ограничено | Полная поддержка (uuid, json, arrays) |
| Подготовленные запросы | Да | Да + автоматическое кеширование |
| Пакетные запросы (batch) | Нет | Да |
| LISTEN/NOTIFY | Нет | Да |
| COPY | Нет | Да |
| Совместимость с database/sql | — | Полная (через pgx/stdlib) |

**Важно:** pgx можно использовать как напрямую, так и через интерфейс `database/sql`. В этом курсе мы используем прямой API pgx для лучшего понимания возможностей.

### Установка PostgreSQL через Docker

Самый простой способ запустить PostgreSQL локально — использовать Docker.

#### docker-compose.yml

```yaml
version: "3.9"

services:
  postgres:
    image: postgres:16-alpine
    container_name: go-course-postgres
    environment:
      POSTGRES_USER: gouser
      POSTGRES_PASSWORD: gopassword
      POSTGRES_DB: gocourse
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U gouser -d gocourse"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

Команды для работы:

```bash
# Запустить PostgreSQL
docker-compose up -d

# Проверить статус
docker-compose ps

# Подключиться к PostgreSQL через psql
docker exec -it go-course-postgres psql -U gouser -d gocourse

# Остановить
docker-compose down

# Остановить и удалить данные
docker-compose down -v
```

### Подключение к PostgreSQL

#### Строка подключения (DSN)

```
postgresql://user:password@host:port/database?sslmode=disable
```

Примеры:
```go
// Локальный Docker
dsn := "postgresql://gouser:gopassword@localhost:5432/gocourse?sslmode=disable"

// Со всеми параметрами
dsn := "postgresql://gouser:gopassword@localhost:5432/gocourse?sslmode=disable&connect_timeout=5"
```

#### Подключение через pgx

```go
import (
    "context"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

// Одиночное подключение (для простых скриптов)
func connectSingle(ctx context.Context, dsn string) (*pgx.Conn, error) {
    conn, err := pgx.Connect(ctx, dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to connect: %w", err)
    }
    return conn, nil
}

// Пул соединений (рекомендуется для серверов)
func connectPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }

    // Настройки пула
    config.MaxConns = 10              // максимум соединений
    config.MinConns = 2               // минимум поддерживаемых соединений
    config.MaxConnLifetime = time.Hour // время жизни соединения
    config.MaxConnIdleTime = 30 * time.Minute // таймаут простоя

    pool, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        return nil, fmt.Errorf("failed to create pool: %w", err)
    }

    // Проверяем подключение
    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping: %w", err)
    }

    return pool, nil
}
```

### CRUD операции

#### CREATE (INSERT)

```go
// Вставка одной записи
func createUser(ctx context.Context, pool *pgxpool.Pool, user User) (int, error) {
    var id int
    err := pool.QueryRow(ctx,
        `INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id`,
        user.Name, user.Email, user.Age,
    ).Scan(&id)

    if err != nil {
        return 0, fmt.Errorf("insert user: %w", err)
    }
    return id, nil
}

// Вставка нескольких записей
func createUsers(ctx context.Context, pool *pgxpool.Pool, users []User) error {
    batch := &pgx.Batch{}
    for _, u := range users {
        batch.Queue(
            `INSERT INTO users (name, email, age) VALUES ($1, $2, $3)`,
            u.Name, u.Email, u.Age,
        )
    }

    results := pool.SendBatch(ctx, batch)
    defer results.Close()

    for range users {
        _, err := results.Exec()
        if err != nil {
            return fmt.Errorf("batch insert: %w", err)
        }
    }
    return nil
}
```

#### READ (SELECT)

```go
// Получение одной записи — QueryRow
func getUserByID(ctx context.Context, pool *pgxpool.Pool, id int) (User, error) {
    var user User
    err := pool.QueryRow(ctx,
        `SELECT id, name, email, age, created_at FROM users WHERE id = $1`,
        id,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt)

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return User{}, ErrUserNotFound
        }
        return User{}, fmt.Errorf("query user: %w", err)
    }
    return user, nil
}

// Получение нескольких записей — Query
func getAllUsers(ctx context.Context, pool *pgxpool.Pool) ([]User, error) {
    rows, err := pool.Query(ctx,
        `SELECT id, name, email, age, created_at FROM users ORDER BY id`,
    )
    if err != nil {
        return nil, fmt.Errorf("query users: %w", err)
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt)
        if err != nil {
            return nil, fmt.Errorf("scan user: %w", err)
        }
        users = append(users, user)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }

    return users, nil
}

// Использование pgx.CollectRows для упрощения
func getAllUsersSimple(ctx context.Context, pool *pgxpool.Pool) ([]User, error) {
    rows, err := pool.Query(ctx,
        `SELECT id, name, email, age, created_at FROM users ORDER BY id`,
    )
    if err != nil {
        return nil, err
    }

    users, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
    if err != nil {
        return nil, err
    }
    return users, nil
}
```

#### UPDATE

```go
func updateUser(ctx context.Context, pool *pgxpool.Pool, id int, user User) error {
    result, err := pool.Exec(ctx,
        `UPDATE users SET name = $1, email = $2, age = $3 WHERE id = $4`,
        user.Name, user.Email, user.Age, id,
    )
    if err != nil {
        return fmt.Errorf("update user: %w", err)
    }

    if result.RowsAffected() == 0 {
        return ErrUserNotFound
    }
    return nil
}
```

#### DELETE

```go
func deleteUser(ctx context.Context, pool *pgxpool.Pool, id int) error {
    result, err := pool.Exec(ctx,
        `DELETE FROM users WHERE id = $1`,
        id,
    )
    if err != nil {
        return fmt.Errorf("delete user: %w", err)
    }

    if result.RowsAffected() == 0 {
        return ErrUserNotFound
    }
    return nil
}
```

### Prepared Statements

pgx автоматически кеширует prepared statements, но можно управлять этим явно:

```go
// Явная подготовка запроса
func prepareStatements(ctx context.Context, conn *pgx.Conn) error {
    _, err := conn.Prepare(ctx, "get_user",
        `SELECT id, name, email, age FROM users WHERE id = $1`)
    if err != nil {
        return err
    }

    _, err = conn.Prepare(ctx, "list_users",
        `SELECT id, name, email, age FROM users ORDER BY id LIMIT $1 OFFSET $2`)
    return err
}

// Использование подготовленного запроса
func getUserPrepared(ctx context.Context, conn *pgx.Conn, id int) (User, error) {
    var user User
    err := conn.QueryRow(ctx, "get_user", id).
        Scan(&user.ID, &user.Name, &user.Email, &user.Age)
    return user, err
}
```

### Транзакции

Транзакции обеспечивают ACID-свойства (Atomicity, Consistency, Isolation, Durability).

#### Базовое использование

```go
func transferMoney(ctx context.Context, pool *pgxpool.Pool, fromID, toID int, amount float64) error {
    // Начинаем транзакцию
    tx, err := pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    // defer Rollback безопасен — если уже был Commit, ничего не произойдёт
    defer tx.Rollback(ctx)

    // Списываем с первого счёта
    result, err := tx.Exec(ctx,
        `UPDATE accounts SET balance = balance - $1 WHERE id = $2 AND balance >= $1`,
        amount, fromID,
    )
    if err != nil {
        return fmt.Errorf("debit: %w", err)
    }
    if result.RowsAffected() == 0 {
        return errors.New("insufficient funds or account not found")
    }

    // Зачисляем на второй счёт
    result, err = tx.Exec(ctx,
        `UPDATE accounts SET balance = balance + $1 WHERE id = $2`,
        amount, toID,
    )
    if err != nil {
        return fmt.Errorf("credit: %w", err)
    }
    if result.RowsAffected() == 0 {
        return errors.New("destination account not found")
    }

    // Фиксируем транзакцию
    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("commit: %w", err)
    }

    return nil
}
```

#### Уровни изоляции

```go
// Транзакция с определённым уровнем изоляции
tx, err := pool.BeginTx(ctx, pgx.TxOptions{
    IsoLevel:   pgx.Serializable,  // ReadCommitted, RepeatableRead, Serializable
    AccessMode: pgx.ReadWrite,      // ReadWrite, ReadOnly
})
```

#### Функция-обёртка для транзакций

```go
// WithTx выполняет функцию в транзакции
func WithTx(ctx context.Context, pool *pgxpool.Pool, fn func(tx pgx.Tx) error) error {
    tx, err := pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    defer tx.Rollback(ctx)

    if err := fn(tx); err != nil {
        return err
    }

    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("commit tx: %w", err)
    }
    return nil
}

// Использование
err := WithTx(ctx, pool, func(tx pgx.Tx) error {
    if _, err := tx.Exec(ctx, "UPDATE ..."); err != nil {
        return err
    }
    if _, err := tx.Exec(ctx, "INSERT ..."); err != nil {
        return err
    }
    return nil
})
```

### Repository Pattern

Паттерн Repository абстрагирует логику доступа к данным от бизнес-логики:

```go
// UserRepository определяет контракт для работы с пользователями
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id int) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    GetAll(ctx context.Context) ([]User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int) error
}

// PostgresUserRepository реализует UserRepository для PostgreSQL
type PostgresUserRepository struct {
    pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
    return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *User) error {
    err := r.pool.QueryRow(ctx,
        `INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id, created_at`,
        user.Name, user.Email, user.Age,
    ).Scan(&user.ID, &user.CreatedAt)
    return err
}

// ... остальные методы
```

### Обработка ошибок PostgreSQL

```go
import "github.com/jackc/pgx/v5/pgconn"

func handlePgError(err error) error {
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
        switch pgErr.Code {
        case "23505": // unique_violation
            return fmt.Errorf("duplicate key: %s", pgErr.ConstraintName)
        case "23503": // foreign_key_violation
            return fmt.Errorf("foreign key violation: %s", pgErr.ConstraintName)
        case "23502": // not_null_violation
            return fmt.Errorf("null value in column: %s", pgErr.ColumnName)
        case "22P02": // invalid_text_representation
            return fmt.Errorf("invalid input syntax")
        default:
            return fmt.Errorf("database error: %s (code: %s)", pgErr.Message, pgErr.Code)
        }
    }
    return err
}
```

### NULL значения

```go
import "github.com/jackc/pgx/v5/pgtype"

type User struct {
    ID        int
    Name      string
    Email     pgtype.Text  // может быть NULL
    Age       pgtype.Int4  // может быть NULL
    DeletedAt pgtype.Timestamp
}

// Проверка на NULL
if user.Email.Valid {
    fmt.Println("Email:", user.Email.String)
} else {
    fmt.Println("Email is NULL")
}

// Установка значения
user.Email = pgtype.Text{String: "test@example.com", Valid: true}

// Установка NULL
user.Email = pgtype.Text{Valid: false}
```

## Установка зависимостей

```bash
# Установка pgx
go get github.com/jackc/pgx/v5
```

## Домашние задания

| # | Задание | Описание |
|---|---------|----------|
| 20.1 | [db-connect](./homework/20.1-db-connect/) | Подключение к PostgreSQL, проверка соединения, настройка пула |
| 20.2 | [db-crud](./homework/20.2-db-crud/) | CRUD операции для User (QueryRow, Query, Exec) |
| 20.3 | [db-transactions](./homework/20.3-db-transactions/) | Транзакции: Tx, Commit, Rollback, перевод денег |
| 20.4 | [db-repository](./homework/20.4-db-repository/) | Repository pattern для абстракции работы с БД |

## Вопросы для самопроверки

Создай файл `answers-20.txt` и напиши ответы на вопросы:

1. В чём разница между `pgx.Connect` и `pgxpool.New`? Когда какой способ подключения использовать?

2. Какой метод использовать для SELECT одной записи, а какой для нескольких? Что вернёт QueryRow, если запись не найдена?

3. Зачем нужны транзакции? Что произойдёт с данными, если между BEGIN и COMMIT произойдёт ошибка?

4. Почему после `pool.Begin(ctx)` следует делать `defer tx.Rollback(ctx)`? Что произойдёт, если Rollback вызвать после Commit?

5. Какие преимущества даёт Repository pattern? Как он помогает в тестировании?

## Дополнительные материалы

- [pgx GitHub](https://github.com/jackc/pgx)
- [pgx Documentation](https://pkg.go.dev/github.com/jackc/pgx/v5)
- [PostgreSQL Tutorial](https://www.postgresqltutorial.com/)
- [Go Database/SQL Tutorial](http://go-database-sql.org/)
- [PostgreSQL Error Codes](https://www.postgresql.org/docs/current/errcodes-appendix.html)
- [Docker Hub: PostgreSQL](https://hub.docker.com/_/postgres)
