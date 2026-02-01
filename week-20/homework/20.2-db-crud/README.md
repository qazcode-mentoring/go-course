# ДЗ 20.2: CRUD операции с PostgreSQL

## Цель

Научиться выполнять CRUD операции (Create, Read, Update, Delete) с PostgreSQL используя pgx. Освоить работу с QueryRow, Query и Exec.

## Подготовка

1. Убедись, что PostgreSQL запущен:

```bash
cd week-20
docker-compose up -d
```

2. Проверь, что таблица users существует (она создаётся автоматически из init.sql):

```bash
docker exec -it go-course-postgres psql -U gouser -d gocourse -c "SELECT * FROM users;"
```

## SQL-схема

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    age INT CHECK (age >= 0 AND age <= 150),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

## Что нужно сделать

Реализовать CRUD операции для работы с таблицей users:

### Структура данных

```go
type User struct {
    ID        int
    Name      string
    Email     string
    Age       int
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 1. CreateUser — создание пользователя

```go
func CreateUser(ctx context.Context, pool *pgxpool.Pool, name, email string, age int) (User, error)
```

- Вставляет нового пользователя в таблицу
- Использует `INSERT ... RETURNING id, created_at, updated_at`
- Возвращает созданного пользователя со всеми полями

### 2. GetUserByID — получение по ID

```go
func GetUserByID(ctx context.Context, pool *pgxpool.Pool, id int) (User, error)
```

- Ищет пользователя по ID
- Использует `QueryRow` для получения одной записи
- Возвращает `ErrUserNotFound`, если пользователь не найден

### 3. GetUserByEmail — получение по email

```go
func GetUserByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (User, error)
```

- Ищет пользователя по email
- Возвращает `ErrUserNotFound`, если пользователь не найден

### 4. GetAllUsers — получение всех пользователей

```go
func GetAllUsers(ctx context.Context, pool *pgxpool.Pool) ([]User, error)
```

- Возвращает список всех пользователей
- Использует `Query` для получения нескольких записей
- Сортирует по ID

### 5. UpdateUser — обновление пользователя

```go
func UpdateUser(ctx context.Context, pool *pgxpool.Pool, id int, name, email string, age int) (User, error)
```

- Обновляет данные пользователя
- Использует `RETURNING` для получения обновлённых данных
- Возвращает `ErrUserNotFound`, если пользователь не найден

### 6. DeleteUser — удаление пользователя

```go
func DeleteUser(ctx context.Context, pool *pgxpool.Pool, id int) error
```

- Удаляет пользователя по ID
- Возвращает `ErrUserNotFound`, если пользователь не существовал

## SQL запросы

```sql
-- CREATE
INSERT INTO users (name, email, age) VALUES ($1, $2, $3)
RETURNING id, name, email, age, created_at, updated_at;

-- READ (один)
SELECT id, name, email, age, created_at, updated_at
FROM users WHERE id = $1;

-- READ (все)
SELECT id, name, email, age, created_at, updated_at
FROM users ORDER BY id;

-- UPDATE
UPDATE users SET name = $1, email = $2, age = $3, updated_at = NOW()
WHERE id = $4
RETURNING id, name, email, age, created_at, updated_at;

-- DELETE
DELETE FROM users WHERE id = $1;
```

## Пример использования

```go
func main() {
    ctx := context.Background()
    pool, _ := pgxpool.New(ctx, dsn)
    defer pool.Close()

    // Создание пользователя
    user, err := CreateUser(ctx, pool, "John Doe", "john@example.com", 28)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created: %+v\n", user)

    // Получение по ID
    user, err = GetUserByID(ctx, pool, user.ID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found: %+v\n", user)

    // Обновление
    user, err = UpdateUser(ctx, pool, user.ID, "John Updated", "john.updated@example.com", 29)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Updated: %+v\n", user)

    // Получение всех
    users, err := GetAllUsers(ctx, pool)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("All users: %d\n", len(users))

    // Удаление
    err = DeleteUser(ctx, pool, user.ID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Deleted successfully")
}
```

Ожидаемый вывод:
```
=== PostgreSQL CRUD Demo ===

--- Create User ---
Created user: ID=4, Name=John Doe, Email=john@example.com, Age=28

--- Get User by ID ---
Found user: ID=4, Name=John Doe, Email=john@example.com

--- Get User by Email ---
Found user: ID=4, Name=John Doe

--- Update User ---
Updated user: ID=4, Name=John Updated, Email=john.updated@example.com, Age=29

--- Get All Users ---
Total users: 4
  - Alice Johnson (alice@example.com)
  - Bob Smith (bob@example.com)
  - Charlie Brown (charlie@example.com)
  - John Updated (john.updated@example.com)

--- Delete User ---
User deleted successfully

--- Verify Deletion ---
GetUserByID: user not found (expected)
```

## Подсказки

### QueryRow для одной записи

```go
var user User
err := pool.QueryRow(ctx, query, args...).Scan(
    &user.ID,
    &user.Name,
    &user.Email,
    &user.Age,
    &user.CreatedAt,
    &user.UpdatedAt,
)
if errors.Is(err, pgx.ErrNoRows) {
    return User{}, ErrUserNotFound
}
```

### Query для нескольких записей

```go
rows, err := pool.Query(ctx, query, args...)
if err != nil {
    return nil, err
}
defer rows.Close()

var users []User
for rows.Next() {
    var user User
    err := rows.Scan(&user.ID, &user.Name, ...)
    if err != nil {
        return nil, err
    }
    users = append(users, user)
}

// Обязательно проверяем ошибки после цикла
if err := rows.Err(); err != nil {
    return nil, err
}
```

### Exec для изменения данных

```go
result, err := pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
if err != nil {
    return err
}
if result.RowsAffected() == 0 {
    return ErrUserNotFound
}
```

## Обработка ошибок PostgreSQL

```go
import "github.com/jackc/pgx/v5/pgconn"

// Проверка уникальности email
var pgErr *pgconn.PgError
if errors.As(err, &pgErr) && pgErr.Code == "23505" {
    return User{}, ErrEmailAlreadyExists
}
```

## Критерии выполнения

- [ ] CreateUser создаёт пользователя и возвращает его со всеми полями
- [ ] CreateUser возвращает ошибку при дублировании email
- [ ] GetUserByID возвращает пользователя или ErrUserNotFound
- [ ] GetUserByEmail возвращает пользователя или ErrUserNotFound
- [ ] GetAllUsers возвращает список всех пользователей
- [ ] UpdateUser обновляет пользователя или возвращает ErrUserNotFound
- [ ] DeleteUser удаляет пользователя или возвращает ErrUserNotFound
- [ ] Все функции корректно обрабатывают ошибки
- [ ] Код компилируется без ошибок
