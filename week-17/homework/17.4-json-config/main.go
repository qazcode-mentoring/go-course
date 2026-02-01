package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ServerConfig содержит настройки HTTP-сервера
type ServerConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// DatabaseConfig содержит настройки подключения к БД
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	SSLMode  string `json:"ssl_mode"`
}

// LoggingConfig содержит настройки логирования
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// Config — главная структура конфигурации приложения
type Config struct {
	AppName  string         `json:"app_name"`
	Version  string         `json:"version"`
	Debug    bool           `json:"debug"`
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Logging  LoggingConfig  `json:"logging"`
	Features []string       `json:"features,omitempty"`
}

// LoadConfig загружает конфигурацию из JSON-файла
func LoadConfig(filename string) (*Config, error) {
	// TODO:
	// 1. Прочитай файл с помощью os.ReadFile
	// 2. Создай переменную Config
	// 3. Распарси JSON в структуру
	// 4. Верни указатель на конфигурацию
	// При ошибке используй fmt.Errorf с %w для wrapping

	// Подавляем предупреждение о неиспользуемых импортах
	_ = os.ReadFile
	_ = json.Unmarshal

	return nil, nil
}

// SaveConfig сохраняет конфигурацию в JSON-файл с форматированием
func SaveConfig(filename string, cfg *Config) error {
	// TODO:
	// 1. Сериализуй конфигурацию с помощью json.MarshalIndent
	// 2. Запиши в файл с помощью os.WriteFile
	// Права доступа: 0644
	return nil
}

// ValidateConfig проверяет обязательные поля и корректность значений
func ValidateConfig(cfg *Config) error {
	// TODO: проверь обязательные поля:
	// - AppName не пустой
	// - Server.Host не пустой
	// - Server.Port в диапазоне 1-65535
	// - Database.Host не пустой
	// - Database.Port в диапазоне 1-65535
	// - Database.User не пустой
	// - Database.DBName не пустой
	//
	// Собери все ошибки в слайс и верни их вместе

	var errors []string

	// Пример проверки:
	// if cfg.AppName == "" {
	//     errors = append(errors, "app_name обязателен")
	// }

	if len(errors) > 0 {
		return fmt.Errorf("ошибки валидации: %s", strings.Join(errors, "; "))
	}
	return nil
}

// ApplyDefaults устанавливает значения по умолчанию для пустых полей
func ApplyDefaults(cfg *Config) {
	// TODO: установи значения по умолчанию:
	// - Server.ReadTimeout: 30 (если 0)
	// - Server.WriteTimeout: 30 (если 0)
	// - Database.Port: 5432 (если 0)
	// - Database.SSLMode: "disable" (если пустой)
	// - Logging.Level: "info" (если пустой)
	// - Logging.Format: "json" (если пустой)
}

// NewDefaultConfig возвращает конфигурацию со значениями по умолчанию
func NewDefaultConfig() *Config {
	// TODO: создай конфигурацию с разумными дефолтами
	// Это полезно когда файл конфигурации не найден
	return &Config{
		AppName: "DefaultApp",
		Version: "0.0.1",
		Debug:   false,
		Server: ServerConfig{
			Host:         "localhost",
			Port:         8080,
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		Database: DatabaseConfig{
			Host:    "localhost",
			Port:    5432,
			SSLMode: "disable",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}
}

// PrintConfig выводит конфигурацию в читаемом виде
func PrintConfig(cfg *Config) {
	fmt.Printf("Application: %s v%s\n", cfg.AppName, cfg.Version)
	fmt.Printf("Debug mode: %v\n", cfg.Debug)
	fmt.Println()

	fmt.Println("Server:")
	fmt.Printf("  Address: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("  Timeouts: read=%ds, write=%ds\n",
		cfg.Server.ReadTimeout, cfg.Server.WriteTimeout)
	fmt.Println()

	fmt.Println("Database:")
	fmt.Printf("  Host: %s:%d\n", cfg.Database.Host, cfg.Database.Port)
	fmt.Printf("  Database: %s\n", cfg.Database.DBName)
	fmt.Printf("  User: %s\n", cfg.Database.User)
	fmt.Printf("  SSL Mode: %s\n", cfg.Database.SSLMode)
	fmt.Println()

	fmt.Println("Logging:")
	fmt.Printf("  Level: %s\n", cfg.Logging.Level)
	fmt.Printf("  Format: %s\n", cfg.Logging.Format)

	if len(cfg.Features) > 0 {
		fmt.Println()
		fmt.Println("Features:", strings.Join(cfg.Features, ", "))
	}
}

// ConnectionString формирует строку подключения к PostgreSQL
func (db *DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.DBName, db.SSLMode)
}

func main() {
	fmt.Println("=== Конфигурация приложения ===")

	configFile := "config.json"

	// Создаём тестовый конфиг-файл
	fmt.Println("\n--- 1. Создание тестового config.json ---")
	testConfig := &Config{
		AppName: "MyAwesomeApp",
		Version: "1.0.0",
		Debug:   true,
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
			// ReadTimeout и WriteTimeout не заданы — должны получить дефолты
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "appuser",
			Password: "secret123",
			DBName:   "appdb",
			// SSLMode не задан — должен получить дефолт
		},
		Logging: LoggingConfig{
			Level: "debug",
			// Format не задан — должен получить дефолт
		},
		Features: []string{"feature_flags", "metrics", "tracing"},
	}

	err := SaveConfig(configFile, testConfig)
	if err != nil {
		fmt.Printf("Ошибка сохранения: %v\n", err)
	} else {
		fmt.Printf("Конфигурация сохранена в '%s'\n", configFile)
	}

	// Загружаем конфигурацию
	fmt.Println("\n--- 2. Загрузка конфигурации ---")
	cfg, err := LoadConfig(configFile)
	if err != nil {
		fmt.Printf("Ошибка загрузки: %v\n", err)
		fmt.Println("Используем конфигурацию по умолчанию...")
		cfg = NewDefaultConfig()
	}

	// Применяем значения по умолчанию
	fmt.Println("\n--- 3. Применение значений по умолчанию ---")
	fmt.Printf("До: ReadTimeout=%d, SSLMode=%q, Format=%q\n",
		cfg.Server.ReadTimeout, cfg.Database.SSLMode, cfg.Logging.Format)

	ApplyDefaults(cfg)

	fmt.Printf("После: ReadTimeout=%d, SSLMode=%q, Format=%q\n",
		cfg.Server.ReadTimeout, cfg.Database.SSLMode, cfg.Logging.Format)

	// Валидация
	fmt.Println("\n--- 4. Валидация ---")
	if err := ValidateConfig(cfg); err != nil {
		fmt.Printf("Ошибка валидации: %v\n", err)
	} else {
		fmt.Println("Конфигурация валидна")
	}

	// Вывод конфигурации
	fmt.Println("\n--- 5. Текущая конфигурация ---")
	PrintConfig(cfg)

	// Строка подключения к БД
	fmt.Println("\n--- 6. Connection String ---")
	fmt.Println(cfg.Database.ConnectionString())

	// Тест невалидной конфигурации
	fmt.Println("\n--- 7. Тест валидации (невалидная конфигурация) ---")
	invalidCfg := &Config{
		// AppName пустой
		Server: ServerConfig{
			Port: 99999, // невалидный порт
		},
		Database: DatabaseConfig{
			// Host и User пустые
			Port: -1, // невалидный порт
		},
	}

	if err := ValidateConfig(invalidCfg); err != nil {
		fmt.Printf("Ожидаемые ошибки: %v\n", err)
	}

	// Загрузка несуществующего файла
	fmt.Println("\n--- 8. Загрузка несуществующего файла ---")
	_, err = LoadConfig("nonexistent.json")
	if err != nil {
		fmt.Printf("Ожидаемая ошибка: %v\n", err)
	}

	// Очистка тестового файла (опционально)
	// os.Remove(configFile)
}
