package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DSN для подключения к локальному PostgreSQL в Docker
const defaultDSN = "postgresql://gouser:gopassword@localhost:5432/gocourse?sslmode=disable"

// PoolConfig содержит настройки пула соединений
type PoolConfig struct {
	DSN             string        // строка подключения
	MaxConns        int32         // максимальное количество соединений
	MinConns        int32         // минимальное количество соединений
	MaxConnLifetime time.Duration // максимальное время жизни соединения
	MaxConnIdleTime time.Duration // максимальное время простоя соединения
}

// DatabaseInfo содержит информацию о подключении к базе данных
type DatabaseInfo struct {
	Version      string // версия PostgreSQL
	Database     string // имя базы данных
	User         string // текущий пользователь
	MaxConns     int32  // максимум соединений в пуле
	CurrentConns int32  // текущее количество соединений
}

// Connect создаёт простое подключение к PostgreSQL
func Connect(ctx context.Context, dsn string) (*pgx.Conn, error) {
	// TODO: реализуй функцию
	// 1. Подключись к базе данных: conn, err := pgx.Connect(ctx, dsn)
	// 2. Проверь ошибку и верни её, если есть
	// 3. Проверь соединение: err = conn.Ping(ctx)
	// 4. Если Ping не удался, закрой соединение и верни ошибку
	// 5. Верни соединение
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	if err = conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("failed to ping: %w", err)
	}

	return conn, nil
}

// ConnectPool создаёт пул соединений к PostgreSQL
func ConnectPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	// TODO: реализуй функцию
	// 1. Создай пул: pool, err := pgxpool.New(ctx, dsn)
	// 2. Проверь ошибку
	// 3. Проверь соединение: err = pool.Ping(ctx)
	// 4. Если Ping не удался, закрой пул и верни ошибку
	// 5. Верни пул
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping: %w", err)
	}
	return pool, nil
}

// ConnectPoolWithConfig создаёт пул соединений с пользовательскими настройками
func ConnectPoolWithConfig(ctx context.Context, cfg PoolConfig) (*pgxpool.Pool, error) {
	// TODO: реализуй функцию
	// 1. Распарси DSN в конфигурацию: config, err := pgxpool.ParseConfig(cfg.DSN)
	// 2. Проверь ошибку
	// 3. Применить настройки из cfg:
	//    config.MaxConns = cfg.MaxConns
	//    config.MinConns = cfg.MinConns
	//    config.MaxConnLifetime = cfg.MaxConnLifetime
	//    config.MaxConnIdleTime = cfg.MaxConnIdleTime
	// 4. Создай пул с конфигурацией: pool, err := pgxpool.NewWithConfig(ctx, config)
	// 5. Проверь соединение через Ping
	// 6. Верни пул или ошибку
	config, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	config.MaxConns = cfg.MaxConns
	config.MinConns = cfg.MinConns
	config.MaxConnLifetime = cfg.MaxConnLifetime
	config.MaxConnIdleTime = cfg.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping: %w", err)
	}

	return pool, nil
}

// GetDatabaseInfo возвращает информацию о подключении к базе данных
func GetDatabaseInfo(ctx context.Context, pool *pgxpool.Pool) (DatabaseInfo, error) {
	// TODO: реализуй функцию
	// 1. Создай переменную info DatabaseInfo
	// 2. Выполни запрос для получения версии:
	//    err := pool.QueryRow(ctx, "SELECT version()").Scan(&info.Version)
	// 3. Выполни запрос для получения имени базы данных:
	//    err = pool.QueryRow(ctx, "SELECT current_database()").Scan(&info.Database)
	// 4. Выполни запрос для получения пользователя:
	//    err = pool.QueryRow(ctx, "SELECT current_user").Scan(&info.User)
	// 5. Получи статистику пула:
	//    stat := pool.Stat()
	//    info.MaxConns = stat.MaxConns()
	//    info.CurrentConns = stat.TotalConns()
	// 6. Верни info
	var info DatabaseInfo

	err := pool.QueryRow(ctx, "SELECT version()").Scan(&info.Version)
	if err != nil {
		return DatabaseInfo{}, fmt.Errorf("failed to get version: %w", err)
	}

	err = pool.QueryRow(ctx, "SELECT current_database()").Scan(&info.Database)
	if err != nil {
		return DatabaseInfo{}, fmt.Errorf("failed to get database name: %w", err)
	}

	err = pool.QueryRow(ctx, "SELECT current_user").Scan(&info.User)
	if err != nil {
		return DatabaseInfo{}, fmt.Errorf("failed to get user: %w", err)
	}

	stat := pool.Stat()
	info.MaxConns = stat.MaxConns()
	info.CurrentConns = stat.TotalConns()

	return info, nil
}

func main() {
	fmt.Println("=== PostgreSQL Connection Demo ===")

	ctx := context.Background()

	// --- Простое подключение ---
	fmt.Println("\n--- Simple Connection ---")
	conn, err := Connect(ctx, defaultDSN)
	if err != nil {
		log.Printf("Failed to connect: %v", err)
	} else {
		fmt.Println("Connected successfully!")
		conn.Close(ctx)
	}

	// --- Пул соединений ---
	fmt.Println("\n--- Connection Pool ---")
	pool, err := ConnectPool(ctx, defaultDSN)
	if err != nil {
		log.Printf("Failed to create pool: %v", err)
	} else {
		fmt.Println("Pool created successfully!")
		pool.Close()
	}

	// --- Пул с настройками ---
	fmt.Println("\n--- Pool with Config ---")
	poolCfg := PoolConfig{
		DSN:             defaultDSN,
		MaxConns:        10,
		MinConns:        2,
		MaxConnLifetime: time.Hour,
		MaxConnIdleTime: 30 * time.Minute,
	}
	configuredPool, err := ConnectPoolWithConfig(ctx, poolCfg)
	if err != nil {
		log.Printf("Failed to create configured pool: %v", err)
	} else {
		fmt.Println("Configured pool created successfully!")

		// --- Информация о базе данных ---
		fmt.Println("\n--- Database Info ---")
		info, err := GetDatabaseInfo(ctx, configuredPool)
		if err != nil {
			log.Printf("Failed to get database info: %v", err)
		} else {
			fmt.Printf("PostgreSQL Version: %s\n", info.Version)
			fmt.Printf("Database: %s\n", info.Database)
			fmt.Printf("User: %s\n", info.User)
			fmt.Printf("Max Connections: %d\n", info.MaxConns)
			fmt.Printf("Current Connections: %d\n", info.CurrentConns)
		}

		configuredPool.Close()
	}
}
