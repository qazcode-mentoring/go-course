# ДЗ 20.1: Подключение к PostgreSQL

## Цель

Научиться подключаться к PostgreSQL из Go с использованием библиотеки pgx. Освоить создание пула соединений и настройку его параметров.

## Подготовка

1. Убедись, что Docker установлен и запущен
2. Перейди в директорию `week-20/` и запусти PostgreSQL:

```bash
cd week-20
docker-compose up -d
```

3. Проверь, что контейнер работает:

```bash
docker-compose ps
```

4. Установи pgx:

```bash
go get github.com/jackc/pgx/v5
```

## Что нужно сделать

Реализовать функции для подключения к PostgreSQL:

### 1. Connect — простое подключение

```go
func Connect(ctx context.Context, dsn string) (*pgx.Conn, error)
```

- Подключается к базе данных
- Проверяет соединение через Ping
- Возвращает соединение или ошибку

### 2. ConnectPool — создание пула соединений

```go
func ConnectPool(ctx context.Context, dsn string) (*pgxpool.Pool, error)
```

- Создаёт пул соединений
- Проверяет подключение через Ping
- Возвращает пул или ошибку

### 3. ConnectPoolWithConfig — пул с настройками

```go
type PoolConfig struct {
    DSN             string
    MaxConns        int32
    MinConns        int32
    MaxConnLifetime time.Duration
    MaxConnIdleTime time.Duration
}

func ConnectPoolWithConfig(ctx context.Context, cfg PoolConfig) (*pgxpool.Pool, error)
```

- Парсит DSN в конфигурацию
- Применяет настройки из PoolConfig
- Создаёт и проверяет пул

### 4. GetDatabaseInfo — информация о подключении

```go
type DatabaseInfo struct {
    Version     string // версия PostgreSQL
    Database    string // имя базы данных
    User        string // текущий пользователь
    MaxConns    int32  // максимум соединений в пуле
    CurrentConns int32 // текущее количество соединений
}

func GetDatabaseInfo(ctx context.Context, pool *pgxpool.Pool) (DatabaseInfo, error)
```

- Выполняет запросы для получения информации
- `SELECT version()` — версия PostgreSQL
- `SELECT current_database()` — имя базы данных
- `SELECT current_user` — текущий пользователь
- Получает статистику пула через методы pgxpool

## Строка подключения (DSN)

Формат:
```
postgresql://user:password@host:port/database?sslmode=disable
```

Для локального Docker:
```
postgresql://gouser:gopassword@localhost:5432/gocourse?sslmode=disable
```

## Пример использования

```go
func main() {
    ctx := context.Background()
    dsn := "postgresql://gouser:gopassword@localhost:5432/gocourse?sslmode=disable"

    // Простое подключение
    conn, err := Connect(ctx, dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close(ctx)
    fmt.Println("Connected successfully!")

    // Пул соединений
    pool, err := ConnectPool(ctx, dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()
    fmt.Println("Pool created successfully!")

    // Пул с настройками
    poolCfg := PoolConfig{
        DSN:             dsn,
        MaxConns:        10,
        MinConns:        2,
        MaxConnLifetime: time.Hour,
        MaxConnIdleTime: 30 * time.Minute,
    }
    configuredPool, err := ConnectPoolWithConfig(ctx, poolCfg)
    if err != nil {
        log.Fatal(err)
    }
    defer configuredPool.Close()

    // Информация о подключении
    info, err := GetDatabaseInfo(ctx, configuredPool)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("PostgreSQL Version: %s\n", info.Version)
    fmt.Printf("Database: %s\n", info.Database)
    fmt.Printf("User: %s\n", info.User)
    fmt.Printf("Max Connections: %d\n", info.MaxConns)
}
```

Ожидаемый вывод:
```
=== PostgreSQL Connection Demo ===

--- Simple Connection ---
Connected successfully!

--- Connection Pool ---
Pool created successfully!

--- Pool with Config ---
Configured pool created successfully!

--- Database Info ---
PostgreSQL Version: PostgreSQL 16.x ...
Database: gocourse
User: gouser
Max Connections: 10
Current Connections: 1
```

## Подсказки

- Для подключения: `pgx.Connect(ctx, dsn)`
- Для создания пула: `pgxpool.New(ctx, dsn)`
- Для парсинга конфигурации: `pgxpool.ParseConfig(dsn)`
- Для проверки соединения: `conn.Ping(ctx)` или `pool.Ping(ctx)`
- Для настройки пула используй поля структуры `*pgxpool.Config`
- Статистика пула: `pool.Stat().MaxConns()`, `pool.Stat().TotalConns()`

## SQL запросы для получения информации

```sql
-- Версия PostgreSQL
SELECT version();

-- Текущая база данных
SELECT current_database();

-- Текущий пользователь
SELECT current_user;
```

## Критерии выполнения

- [ ] Docker с PostgreSQL запущен и работает
- [ ] Connect успешно подключается и возвращает соединение
- [ ] ConnectPool создаёт работающий пул соединений
- [ ] ConnectPoolWithConfig применяет все настройки из PoolConfig
- [ ] GetDatabaseInfo возвращает корректную информацию
- [ ] При ошибках подключения возвращаются понятные ошибки
- [ ] Код компилируется без ошибок
