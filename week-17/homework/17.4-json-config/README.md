# ДЗ 17.4: Загрузка конфигурации из JSON

## Цель

Научиться загружать конфигурацию приложения из JSON-файла — типичная задача в реальных проектах. Реализовать валидацию конфигурации и значения по умолчанию.

## Что нужно сделать

1. Определить структуру конфигурации `Config` с вложенными секциями
2. Реализовать загрузку конфигурации из JSON-файла
3. Реализовать валидацию обязательных полей
4. Применить значения по умолчанию для опциональных полей
5. Сохранить конфигурацию обратно в файл

## Структура конфигурации

```go
type ServerConfig struct {
    Host         string `json:"host"`
    Port         int    `json:"port"`
    ReadTimeout  int    `json:"read_timeout"`   // секунды, по умолчанию 30
    WriteTimeout int    `json:"write_timeout"`  // секунды, по умолчанию 30
}

type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    User     string `json:"user"`
    Password string `json:"password"`
    DBName   string `json:"db_name"`
    SSLMode  string `json:"ssl_mode"` // по умолчанию "disable"
}

type LoggingConfig struct {
    Level  string `json:"level"`  // debug, info, warn, error. По умолчанию "info"
    Format string `json:"format"` // json, text. По умолчанию "json"
}

type Config struct {
    AppName  string          `json:"app_name"`
    Version  string          `json:"version"`
    Debug    bool            `json:"debug"`
    Server   ServerConfig    `json:"server"`
    Database DatabaseConfig  `json:"database"`
    Logging  LoggingConfig   `json:"logging"`
    Features []string        `json:"features,omitempty"`
}
```

## Пример JSON-файла (config.json)

```json
{
  "app_name": "MyApp",
  "version": "1.0.0",
  "debug": false,
  "server": {
    "host": "0.0.0.0",
    "port": 8080,
    "read_timeout": 60,
    "write_timeout": 60
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "postgres",
    "password": "secret",
    "db_name": "myapp"
  },
  "logging": {
    "level": "debug",
    "format": "text"
  },
  "features": ["feature_a", "feature_b"]
}
```

## Функции для реализации

```go
// LoadConfig загружает конфигурацию из JSON-файла
func LoadConfig(filename string) (*Config, error)

// SaveConfig сохраняет конфигурацию в JSON-файл
func SaveConfig(filename string, cfg *Config) error

// ValidateConfig проверяет обязательные поля конфигурации
func ValidateConfig(cfg *Config) error

// ApplyDefaults устанавливает значения по умолчанию для пустых полей
func ApplyDefaults(cfg *Config)

// NewDefaultConfig возвращает конфигурацию со значениями по умолчанию
func NewDefaultConfig() *Config
```

## Правила валидации

Обязательные поля (должны быть непустыми):
- `AppName`
- `Server.Host`
- `Server.Port` (должен быть > 0 и < 65536)
- `Database.Host`
- `Database.Port`
- `Database.User`
- `Database.DBName`

## Значения по умолчанию

| Поле | Значение по умолчанию |
|------|----------------------|
| Server.ReadTimeout | 30 |
| Server.WriteTimeout | 30 |
| Database.Port | 5432 |
| Database.SSLMode | "disable" |
| Logging.Level | "info" |
| Logging.Format | "json" |

## Пример использования

```go
func main() {
    // Загрузка конфигурации
    cfg, err := LoadConfig("config.json")
    if err != nil {
        log.Fatalf("Ошибка загрузки: %v", err)
    }

    // Применение значений по умолчанию
    ApplyDefaults(cfg)

    // Валидация
    if err := ValidateConfig(cfg); err != nil {
        log.Fatalf("Невалидная конфигурация: %v", err)
    }

    // Использование
    fmt.Printf("Starting %s v%s\n", cfg.AppName, cfg.Version)
    fmt.Printf("Server: %s:%d\n", cfg.Server.Host, cfg.Server.Port)
    fmt.Printf("Database: %s@%s:%d/%s\n",
        cfg.Database.User, cfg.Database.Host,
        cfg.Database.Port, cfg.Database.DBName)

    // Сохранение (например, с обновлённой версией)
    cfg.Version = "1.0.1"
    SaveConfig("config_updated.json", cfg)
}
```

## Подсказки

### Загрузка конфигурации
```go
func LoadConfig(filename string) (*Config, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("чтение файла: %w", err)
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("парсинг JSON: %w", err)
    }

    return &cfg, nil
}
```

### Валидация с несколькими ошибками
```go
func ValidateConfig(cfg *Config) error {
    var errors []string

    if cfg.AppName == "" {
        errors = append(errors, "app_name обязателен")
    }
    if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
        errors = append(errors, "server.port должен быть 1-65535")
    }
    // ... другие проверки

    if len(errors) > 0 {
        return fmt.Errorf("ошибки валидации: %s", strings.Join(errors, "; "))
    }
    return nil
}
```

### Применение значений по умолчанию
```go
func ApplyDefaults(cfg *Config) {
    if cfg.Server.ReadTimeout == 0 {
        cfg.Server.ReadTimeout = 30
    }
    // ... другие значения
}
```

## Критерии выполнения

- [ ] Структуры `Config`, `ServerConfig`, `DatabaseConfig`, `LoggingConfig` определены с правильными тегами
- [ ] `LoadConfig` загружает и парсит JSON-файл
- [ ] `SaveConfig` сохраняет конфигурацию с форматированием
- [ ] `ValidateConfig` проверяет все обязательные поля
- [ ] `ValidateConfig` проверяет корректность порта (1-65535)
- [ ] `ApplyDefaults` устанавливает значения по умолчанию только для пустых полей
- [ ] `NewDefaultConfig` создаёт конфигурацию с разумными дефолтами
- [ ] Ошибки содержат понятные сообщения
- [ ] Код компилируется и работает корректно

## Бонус

Реализуй функцию `MergeConfigs(base, override *Config) *Config`, которая объединяет две конфигурации: берёт значения из `base` и переопределяет непустыми значениями из `override`. Это полезно для:
- Базовая конфигурация + переопределения для конкретного окружения
- Дефолтные значения + пользовательские настройки
