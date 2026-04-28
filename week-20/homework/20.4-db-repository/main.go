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
	query := `INSERT INTO users (name, email, age) 
              VALUES ($1, $2, $3) 
              RETURNING id, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query, user.Name, user.Email, user.Age).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return ErrEmailAlreadyExists
			}
		}
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// GetByID возвращает пользователя по ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	user := &User{}
	query := `SELECT id, name, email, age, created_at, updated_at 
              FROM users WHERE id = $1`

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	user := &User{}
	query := `SELECT id, name, email, age, created_at, updated_at 
              FROM users WHERE email = $1`

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

// GetAll возвращает всех пользователей
func (r *PostgresUserRepository) GetAll(ctx context.Context) ([]*User, error) {
	query := `SELECT id, name, email, age, created_at, updated_at FROM users`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query all users: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{} // Создаем новый объект на каждой итерации
		err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows loop: %w", err)
	}

	return users, nil
}

// Update обновляет данные пользователя
func (r *PostgresUserRepository) Update(ctx context.Context, user *User) error {
	query := `UPDATE users 
              SET name = $1, email = $2, age = $3, updated_at = NOW() 
              WHERE id = $4 
              RETURNING updated_at`

	err := r.pool.QueryRow(ctx, query, user.Name, user.Email, user.Age, user.ID).Scan(&user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrEmailAlreadyExists
		}
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
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
	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверка уникальности email
	for _, u := range r.users {
		if u.Email == user.Email {
			return ErrEmailAlreadyExists
		}
	}

	user.ID = r.nextID
	r.nextID++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	userCopy := *user
	r.users[user.ID] = &userCopy

	return nil
}

// GetByID возвращает пользователя по ID
func (r *InMemoryUserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}

	userCopy := *user
	return &userCopy, nil
}

// GetByEmail возвращает пользователя по email
func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Email == email {
			userCopy := *u
			return &userCopy, nil
		}
	}
	return nil, ErrUserNotFound
}

// GetAll возвращает всех пользователей
func (r *InMemoryUserRepository) GetAll(ctx context.Context) ([]*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*User, 0, len(r.users))
	for _, u := range r.users {
		userCopy := *u
		users = append(users, &userCopy)
	}
	return users, nil
}

// Update обновляет данные пользователя
func (r *InMemoryUserRepository) Update(ctx context.Context, user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.users[user.ID]
	if !ok {
		return ErrUserNotFound
	}

	for _, u := range r.users {
		if u.Email == user.Email && u.ID != user.ID {
			return ErrEmailAlreadyExists
		}
	}

	user.UpdatedAt = time.Now()
	userCopy := *user
	r.users[user.ID] = &userCopy

	return nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return ErrUserNotFound
	}

	delete(r.users, id)
	return nil
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
	// 1. Валидация
	if strings.TrimSpace(name) == "" {
		return nil, ErrNameRequired
	}
	if strings.TrimSpace(email) == "" {
		return nil, ErrEmailRequired
	}
	if age < 0 || age > 150 {
		return nil, ErrInvalidAge
	}

	// 2. Нормализация и создание объекта
	user := &User{
		Name:  strings.TrimSpace(name),
		Email: strings.ToLower(strings.TrimSpace(email)),
		Age:   age,
	}

	// 3. Сохранение через репозиторий
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err // Здесь может вернуться ErrEmailAlreadyExists из репозитория
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]*User, error) {
	return s.repo.GetAll(ctx)
}

// UpdateUser обновляет пользователя
func (s *UserService) UpdateUser(ctx context.Context, id int, name, email string, age int) (*User, error) {
	// 1. Сначала проверяем, существует ли пользователь
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err // Вернет ErrUserNotFound
	}

	// 2. Валидация новых данных
	if strings.TrimSpace(name) == "" {
		return nil, ErrNameRequired
	}
	if strings.TrimSpace(email) == "" {
		return nil, ErrEmailRequired
	}
	if age < 0 || age > 150 {
		return nil, ErrInvalidAge
	}

	// 3. Обновляем поля
	user.Name = strings.TrimSpace(name)
	user.Email = strings.ToLower(strings.TrimSpace(email))
	user.Age = age

	// 4. Сохраняем
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser удаляет пользователя
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
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
