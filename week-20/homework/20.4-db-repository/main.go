package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DSN для подключения к локальному PostgreSQL в Docker
const defaultDSN = "postgresql://gouser:gopassword@localhost:5432/gocourse?sslmode=disable"

// Ошибки приложения
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrNameRequired       = errors.New("name is required")
	ErrEmailRequired      = errors.New("email is required")
	ErrInvalidAge         = errors.New("age must be between 0 and 150")
)

// User представляет пользователя
type User struct {
	ID        int
	Name      string
	Email     string
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// String возвращает строковое представление пользователя
func (u *User) String() string {
	return fmt.Sprintf("User{ID: %d, Name: %s, Email: %s, Age: %d}", u.ID, u.Name, u.Email, u.Age)
}

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int) error
}

// ============================================================================
// PostgresUserRepository — реализация для PostgreSQL
// ============================================================================

// PostgresUserRepository реализует UserRepository для PostgreSQL
type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresUserRepository создаёт новый PostgreSQL репозиторий
func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

// Create создаёт нового пользователя в базе данных
func (r *PostgresUserRepository) Create(ctx context.Context, user *User) error {
	// TODO: реализуй метод
	// 1. Выполни INSERT с RETURNING:
	//    err := r.pool.QueryRow(ctx,
	//        `INSERT INTO users (name, email, age)
	//         VALUES ($1, $2, $3)
	//         RETURNING id, created_at, updated_at`,
	//        user.Name, user.Email, user.Age,
	//    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	// 2. Обработай ошибку уникальности email (код 23505)
	// 3. Верни nil при успехе
	_ = ctx
	_ = user
	return fmt.Errorf("not implemented")
}

// GetByID возвращает пользователя по ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	// TODO: реализуй метод
	// 1. Выполни SELECT запрос
	// 2. Проверь на pgx.ErrNoRows -> ErrUserNotFound
	// 3. Верни указатель на пользователя
	_ = ctx
	_ = id
	return nil, fmt.Errorf("not implemented")
}

// GetByEmail возвращает пользователя по email
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	// TODO: реализуй метод
	_ = ctx
	_ = email
	return nil, fmt.Errorf("not implemented")
}

// GetAll возвращает всех пользователей
func (r *PostgresUserRepository) GetAll(ctx context.Context) ([]*User, error) {
	// TODO: реализуй метод
	// 1. Выполни SELECT запрос
	// 2. Итерируй по rows
	// 3. Создавай указатели на User и добавляй в срез
	_ = ctx
	return nil, fmt.Errorf("not implemented")
}

// Update обновляет данные пользователя
func (r *PostgresUserRepository) Update(ctx context.Context, user *User) error {
	// TODO: реализуй метод
	// 1. Выполни UPDATE с RETURNING updated_at
	// 2. Проверь на pgx.ErrNoRows -> ErrUserNotFound
	// 3. Обработай ошибку уникальности email
	_ = ctx
	_ = user
	return fmt.Errorf("not implemented")
}

// Delete удаляет пользователя по ID
func (r *PostgresUserRepository) Delete(ctx context.Context, id int) error {
	// TODO: реализуй метод
	// 1. Выполни DELETE
	// 2. Проверь RowsAffected() == 0 -> ErrUserNotFound
	_ = ctx
	_ = id
	return fmt.Errorf("not implemented")
}

// ============================================================================
// InMemoryUserRepository — реализация для тестирования
// ============================================================================

// InMemoryUserRepository реализует UserRepository с хранением в памяти
type InMemoryUserRepository struct {
	mu     sync.RWMutex
	users  map[int]*User
	nextID int
}

// NewInMemoryUserRepository создаёт новый in-memory репозиторий
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

// Create создаёт нового пользователя в памяти
func (r *InMemoryUserRepository) Create(ctx context.Context, user *User) error {
	// TODO: реализуй метод
	// 1. Захвати Lock: r.mu.Lock() / defer r.mu.Unlock()
	// 2. Проверь уникальность email среди существующих пользователей
	// 3. Если email уже существует -> return ErrEmailAlreadyExists
	// 4. Установи ID, CreatedAt, UpdatedAt
	// 5. Создай копию user и сохрани в map
	_ = ctx
	_ = user
	return fmt.Errorf("not implemented")
}

// GetByID возвращает пользователя по ID
func (r *InMemoryUserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	// TODO: реализуй метод
	// 1. Захвати RLock
	// 2. Найди пользователя в map
	// 3. Если не найден -> ErrUserNotFound
	// 4. Верни копию пользователя (чтобы избежать race condition)
	_ = ctx
	_ = id
	return nil, fmt.Errorf("not implemented")
}

// GetByEmail возвращает пользователя по email
func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	// TODO: реализуй метод
	_ = ctx
	_ = email
	return nil, fmt.Errorf("not implemented")
}

// GetAll возвращает всех пользователей
func (r *InMemoryUserRepository) GetAll(ctx context.Context) ([]*User, error) {
	// TODO: реализуй метод
	// 1. Захвати RLock
	// 2. Создай срез и добавь копии всех пользователей
	_ = ctx
	return nil, fmt.Errorf("not implemented")
}

// Update обновляет данные пользователя
func (r *InMemoryUserRepository) Update(ctx context.Context, user *User) error {
	// TODO: реализуй метод
	// 1. Захвати Lock
	// 2. Проверь, что пользователь существует
	// 3. Проверь уникальность email (кроме текущего пользователя)
	// 4. Обнови данные и UpdatedAt
	_ = ctx
	_ = user
	return fmt.Errorf("not implemented")
}

// Delete удаляет пользователя по ID
func (r *InMemoryUserRepository) Delete(ctx context.Context, id int) error {
	// TODO: реализуй метод
	// 1. Захвати Lock
	// 2. Проверь, что пользователь существует
	// 3. Удали из map
	_ = ctx
	_ = id
	return fmt.Errorf("not implemented")
}

// ============================================================================
// UserService — бизнес-логика
// ============================================================================

// UserService предоставляет бизнес-логику для работы с пользователями
type UserService struct {
	repo UserRepository
}

// NewUserService создаёт новый сервис пользователей
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// RegisterUser регистрирует нового пользователя
func (s *UserService) RegisterUser(ctx context.Context, name, email string, age int) (*User, error) {
	// TODO: реализуй метод
	// 1. Валидация имени: if strings.TrimSpace(name) == "" -> ErrNameRequired
	// 2. Валидация email: if strings.TrimSpace(email) == "" -> ErrEmailRequired
	// 3. Валидация возраста: if age < 0 || age > 150 -> ErrInvalidAge
	// 4. Создай User с нормализованными данными
	// 5. Вызови s.repo.Create(ctx, user)
	// 6. Верни пользователя
	_ = ctx
	_ = name
	_ = email
	_ = age
	return nil, fmt.Errorf("not implemented")
}

// GetUser возвращает пользователя по ID
func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
	// TODO: реализуй метод
	// Просто делегируй в репозиторий
	_ = ctx
	_ = id
	return nil, fmt.Errorf("not implemented")
}

// ListUsers возвращает всех пользователей
func (s *UserService) ListUsers(ctx context.Context) ([]*User, error) {
	// TODO: реализуй метод
	_ = ctx
	return nil, fmt.Errorf("not implemented")
}

// UpdateUser обновляет пользователя
func (s *UserService) UpdateUser(ctx context.Context, id int, name, email string, age int) (*User, error) {
	// TODO: реализуй метод
	// 1. Получи пользователя по ID
	// 2. Валидируй входные данные
	// 3. Обнови поля
	// 4. Вызови repo.Update
	_ = ctx
	_ = id
	_ = name
	_ = email
	_ = age
	return nil, fmt.Errorf("not implemented")
}

// DeleteUser удаляет пользователя
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	// TODO: реализуй метод
	_ = ctx
	_ = id
	return fmt.Errorf("not implemented")
}

func main() {
	fmt.Println("=== Repository Pattern Demo ===")

	ctx := context.Background()

	// --- PostgreSQL Repository ---
	fmt.Println("\n--- Using PostgreSQL Repository ---")

	pool, err := pgxpool.New(ctx, defaultDSN)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		log.Println("Continuing with in-memory repository only...")
	} else {
		defer pool.Close()

		if err := pool.Ping(ctx); err != nil {
			log.Printf("Failed to ping database: %v", err)
		} else {
			pgRepo := NewPostgresUserRepository(pool)
			pgService := NewUserService(pgRepo)

			user, err := pgService.RegisterUser(ctx, "John Doe", "john.doe@example.com", 30)
			if err != nil {
				log.Printf("Failed to register user: %v", err)
			} else {
				fmt.Printf("Registered user: %s\n", user)

				foundUser, err := pgService.GetUser(ctx, user.ID)
				if err != nil {
					log.Printf("Failed to get user: %v", err)
				} else {
					fmt.Printf("Found user: %s\n", foundUser)
				}

				users, err := pgService.ListUsers(ctx)
				if err != nil {
					log.Printf("Failed to list users: %v", err)
				} else {
					fmt.Printf("All users: %d\n", len(users))
				}

				// Удаляем тестового пользователя
				_ = pgService.DeleteUser(ctx, user.ID)
			}
		}
	}

	// --- In-Memory Repository ---
	fmt.Println("\n--- Using In-Memory Repository ---")

	memRepo := NewInMemoryUserRepository()
	memService := NewUserService(memRepo)

	// Создаём пользователей
	user1, err := memService.RegisterUser(ctx, "Alice", "alice@test.com", 25)
	if err != nil {
		log.Printf("Failed to register user: %v", err)
	} else {
		fmt.Printf("Registered user: %s\n", user1)
	}

	user2, err := memService.RegisterUser(ctx, "Bob", "bob@test.com", 30)
	if err != nil {
		log.Printf("Failed to register user: %v", err)
	} else {
		fmt.Printf("Registered user: %s\n", user2)
	}

	// Список пользователей
	users, err := memService.ListUsers(ctx)
	if err != nil {
		log.Printf("Failed to list users: %v", err)
	} else {
		fmt.Printf("All users: %d\n", len(users))
	}

	// --- Демонстрация полиморфизма ---
	fmt.Println("\n--- Demonstrating Interface Flexibility ---")
	fmt.Println("Both repositories implement the same interface!")

	// Функция, работающая с любым репозиторием
	printRepoType := func(repo UserRepository) {
		fmt.Printf("Repository type: %T\n", repo)
	}

	if pool != nil {
		pgRepo := NewPostgresUserRepository(pool)
		printRepoType(pgRepo)
	}
	printRepoType(memRepo)

	// Используем импорты
	_ = pgx.ErrNoRows
	_ = pgconn.PgError{}
	_ = strings.TrimSpace
}
