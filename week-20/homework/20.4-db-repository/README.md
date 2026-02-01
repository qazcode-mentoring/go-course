# ДЗ 20.4: Repository Pattern

## Цель

Научиться применять паттерн Repository для абстракции работы с базой данных. Понять, как этот паттерн помогает в тестировании и поддержке кода.

## Зачем нужен Repository Pattern?

1. **Абстракция** — бизнес-логика не знает о деталях хранения данных
2. **Тестирование** — легко подменить реальную БД на mock в тестах
3. **Гибкость** — можно сменить хранилище (PostgreSQL -> MongoDB) без изменения бизнес-логики
4. **Единая точка доступа** — все операции с сущностью в одном месте

## Подготовка

1. Убедись, что PostgreSQL запущен:

```bash
cd week-20
docker-compose up -d
```

## Что нужно сделать

### 1. Определить интерфейс UserRepository

```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id int) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    GetAll(ctx context.Context) ([]*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int) error
}
```

### 2. Реализовать PostgresUserRepository

```go
type PostgresUserRepository struct {
    pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository
```

Реализовать все методы интерфейса UserRepository:

- **Create** — вставляет пользователя, заполняет ID и CreatedAt
- **GetByID** — возвращает пользователя или ErrUserNotFound
- **GetByEmail** — возвращает пользователя или ErrUserNotFound
- **GetAll** — возвращает всех пользователей
- **Update** — обновляет пользователя, обновляет UpdatedAt
- **Delete** — удаляет пользователя

### 3. Реализовать InMemoryUserRepository (для тестов)

```go
type InMemoryUserRepository struct {
    mu     sync.RWMutex
    users  map[int]*User
    nextID int
}

func NewInMemoryUserRepository() *InMemoryUserRepository
```

Реализовать те же методы, но с хранением в памяти.

### 4. Создать UserService

```go
type UserService struct {
    repo UserRepository
}

func NewUserService(repo UserRepository) *UserService

// Бизнес-логика
func (s *UserService) RegisterUser(ctx context.Context, name, email string, age int) (*User, error)
func (s *UserService) GetUser(ctx context.Context, id int) (*User, error)
func (s *UserService) ListUsers(ctx context.Context) ([]*User, error)
func (s *UserService) UpdateUser(ctx context.Context, id int, name, email string, age int) (*User, error)
func (s *UserService) DeleteUser(ctx context.Context, id int) error
```

**Важно:** UserService зависит от интерфейса UserRepository, а не от конкретной реализации.

## Структуры данных

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

## Пример использования

```go
func main() {
    ctx := context.Background()

    // Можно использовать PostgreSQL...
    pool, _ := pgxpool.New(ctx, dsn)
    pgRepo := NewPostgresUserRepository(pool)
    service := NewUserService(pgRepo)

    // ...или in-memory для тестов
    memRepo := NewInMemoryUserRepository()
    testService := NewUserService(memRepo)

    // Бизнес-логика одинаковая!
    user, err := service.RegisterUser(ctx, "John", "john@example.com", 28)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Registered: %+v\n", user)

    testUser, _ := testService.RegisterUser(ctx, "Test", "test@example.com", 25)
    fmt.Printf("Test user: %+v\n", testUser)
}
```

Ожидаемый вывод:
```
=== Repository Pattern Demo ===

--- Using PostgreSQL Repository ---
Registered user: User{ID: 4, Name: John Doe, Email: john.doe@example.com}
Found user: User{ID: 4, Name: John Doe, Email: john.doe@example.com}
All users: 4

--- Using In-Memory Repository ---
Registered user: User{ID: 1, Name: Alice, Email: alice@test.com}
Registered user: User{ID: 2, Name: Bob, Email: bob@test.com}
All users: 2

--- Demonstrating Interface Flexibility ---
Both repositories implement the same interface!
PostgreSQL repo type: *main.PostgresUserRepository
InMemory repo type: *main.InMemoryUserRepository
```

## Подсказки

### PostgresUserRepository.Create

```go
func (r *PostgresUserRepository) Create(ctx context.Context, user *User) error {
    err := r.pool.QueryRow(ctx,
        `INSERT INTO users (name, email, age)
         VALUES ($1, $2, $3)
         RETURNING id, created_at, updated_at`,
        user.Name, user.Email, user.Age,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            return ErrEmailAlreadyExists
        }
        return fmt.Errorf("create user: %w", err)
    }
    return nil
}
```

### InMemoryUserRepository.Create

```go
func (r *InMemoryUserRepository) Create(ctx context.Context, user *User) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    // Проверяем уникальность email
    for _, u := range r.users {
        if u.Email == user.Email {
            return ErrEmailAlreadyExists
        }
    }

    user.ID = r.nextID
    r.nextID++
    user.CreatedAt = time.Now()
    user.UpdatedAt = user.CreatedAt

    // Копируем, чтобы избежать побочных эффектов
    stored := *user
    r.users[user.ID] = &stored

    return nil
}
```

### UserService.RegisterUser

```go
func (s *UserService) RegisterUser(ctx context.Context, name, email string, age int) (*User, error) {
    // Валидация
    if strings.TrimSpace(name) == "" {
        return nil, ErrNameRequired
    }
    if strings.TrimSpace(email) == "" {
        return nil, ErrEmailRequired
    }
    if age < 0 || age > 150 {
        return nil, ErrInvalidAge
    }

    user := &User{
        Name:  strings.TrimSpace(name),
        Email: strings.ToLower(strings.TrimSpace(email)),
        Age:   age,
    }

    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}
```

## Преимущества для тестирования

```go
func TestUserService_RegisterUser(t *testing.T) {
    // Используем in-memory репозиторий
    repo := NewInMemoryUserRepository()
    service := NewUserService(repo)

    ctx := context.Background()

    // Тест создания пользователя
    user, err := service.RegisterUser(ctx, "Test User", "test@example.com", 25)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if user.ID == 0 {
        t.Error("expected non-zero ID")
    }

    // Тест дублирования email
    _, err = service.RegisterUser(ctx, "Another", "test@example.com", 30)
    if !errors.Is(err, ErrEmailAlreadyExists) {
        t.Errorf("expected ErrEmailAlreadyExists, got %v", err)
    }
}
```

## Критерии выполнения

- [ ] Интерфейс UserRepository определён корректно
- [ ] PostgresUserRepository реализует все методы интерфейса
- [ ] InMemoryUserRepository реализует все методы интерфейса
- [ ] UserService зависит только от интерфейса (не от конкретной реализации)
- [ ] UserService содержит валидацию входных данных
- [ ] Оба репозитория возвращают одинаковые ошибки (ErrUserNotFound, ErrEmailAlreadyExists)
- [ ] InMemoryUserRepository потокобезопасен (использует sync.RWMutex)
- [ ] Код компилируется без ошибок
