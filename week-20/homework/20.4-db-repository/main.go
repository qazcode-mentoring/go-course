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
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (name, email, age)
				VALUES ($1, $2, $3)
				RETURNING id, name, email, age, created_at, updated_at`,
		user.Name, user.Email, user.Age,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrEmailAlreadyExists
	}

	return nil
}

// GetByID возвращает пользователя по ID
func (r *PostgresUserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	// TODO: реализуй метод
	// 1. Выполни SELECT запрос
	// 2. Проверь на pgx.ErrNoRows -> ErrUserNotFound
	// 3. Верни указатель на пользователя
	var user User
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, email, age, created_at, updated_at FROM users WHERE id = $1`, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

// GetByEmail возвращает пользователя по email
func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	// TODO: реализуй метод
	var user User
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, email, age, created_at, updated_at FROM users WHERE email = $1`, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	return &user, nil
}

// GetAll возвращает всех пользователей
func (r *PostgresUserRepository) GetAll(ctx context.Context) ([]*User, error) {
	// TODO: реализуй метод
	// 1. Выполни SELECT запрос
	// 2. Итерируй по rows
	// 3. Создавай указатели на User и добавляй в срез
	rows, err := r.pool.Query(ctx, `SELECT id, name, email, age, created_at, updated_at FROM users`)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Age,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}

	return users, nil
}

// Update обновляет данные пользователя
func (r *PostgresUserRepository) Update(ctx context.Context, user *User) error {
	// TODO: реализуй метод
	// 1. Выполни UPDATE с RETURNING updated_at
	// 2. Проверь на pgx.ErrNoRows -> ErrUserNotFound
	// 3. Обработай ошибку уникальности email
	err := r.pool.QueryRow(ctx,
		`UPDATE users SET name = $1, email = $2, age = $3, updated_at = NOW() WHERE id = $4 RETURNING updated_at`,
		user.Name, user.Email, user.Age, user.ID).Scan(
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrUserNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrEmailAlreadyExists
	}

	return nil
}

// Delete удаляет пользователя по ID
func (r *PostgresUserRepository) Delete(ctx context.Context, id int) error {
	// TODO: реализуй метод
	// 1. Выполни DELETE
	// 2. Проверь RowsAffected() == 0 -> ErrUserNotFound
	result, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
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
	// TODO: реализуй метод
	// 1. Захвати Lock: r.mu.Lock() / defer r.mu.Unlock()
	// 2. Проверь уникальность email среди существующих пользователей
	// 3. Если email уже существует -> return ErrEmailAlreadyExists
	// 4. Установи ID, CreatedAt, UpdatedAt
	// 5. Создай копию user и сохрани в map
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.users {
		if u.Email == user.Email {
			return ErrEmailAlreadyExists
		}
	}

	var addUser User
	now := time.Now()
	addUser.ID = r.nextID
	r.nextID++
	addUser.Name = user.Name
	addUser.Email = user.Email
	addUser.Age = user.Age
	addUser.CreatedAt = now
	addUser.UpdatedAt = now

	r.users[addUser.ID] = &addUser

	user.ID = addUser.ID
	user.CreatedAt = now
	user.UpdatedAt = now

	return nil
}

// GetByID возвращает пользователя по ID
func (r *InMemoryUserRepository) GetByID(ctx context.Context, id int) (*User, error) {
	// TODO: реализуй метод
	// 1. Захвати RLock
	// 2. Найди пользователя в map
	// 3. Если не найден -> ErrUserNotFound
	// 4. Верни копию пользователя (чтобы избежать race condition)
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}

	user := User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Age:       u.Age,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	return &user, nil
}

// GetByEmail возвращает пользователя по email
func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	// TODO: реализуй метод
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Email == email {
			user := User{
				ID:        u.ID,
				Name:      u.Name,
				Email:     u.Email,
				Age:       u.Age,
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
			}

			return &user, nil
		}
	}

	return nil, ErrUserNotFound
}

// GetAll возвращает всех пользователей
func (r *InMemoryUserRepository) GetAll(ctx context.Context) ([]*User, error) {
	// TODO: реализуй метод
	// 1. Захвати RLock
	// 2. Создай срез и добавь копии всех пользователей
	r.mu.RLock()
	defer r.mu.RUnlock()
	var users []*User

	for _, u := range r.users {
		user := User{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Age:       u.Age,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		}

		users = append(users, &user)
	}

	return users, nil
}

// Update обновляет данные пользователя
func (r *InMemoryUserRepository) Update(ctx context.Context, user *User) error {
	// TODO: реализуй метод
	// 1. Захвати Lock
	// 2. Проверь, что пользователь существует
	// 3. Проверь уникальность email (кроме текущего пользователя)
	// 4. Обнови данные и UpdatedAt
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[user.ID]
	if !ok {
		return ErrUserNotFound
	}

	for _, us := range r.users {
		if us.ID == user.ID {
			continue
		}
		if us.Email == user.Email {
			return ErrEmailAlreadyExists
		}
	}

	u.Name = user.Name
	u.Email = user.Email
	u.Age = user.Age
	u.UpdatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return nil
}

// Delete удаляет пользователя по ID
func (r *InMemoryUserRepository) Delete(ctx context.Context, id int) error {
	// TODO: реализуй метод
	// 1. Захвати Lock
	// 2. Проверь, что пользователь существует
	// 3. Удали из map
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.users[id]
	if !ok {
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
	// TODO: реализуй метод
	// 1. Валидация имени: if strings.TrimSpace(name) == "" -> ErrNameRequired
	// 2. Валидация email: if strings.TrimSpace(email) == "" -> ErrEmailRequired
	// 3. Валидация возраста: if age < 0 || age > 150 -> ErrInvalidAge
	// 4. Создай User с нормализованными данными
	// 5. Вызови s.repo.Create(ctx, user)
	// 6. Верни пользователя
	if strings.TrimSpace(name) == "" {
		return nil, ErrNameRequired
	}
	if strings.TrimSpace(email) == "" {
		return nil, ErrEmailRequired
	}
	if age < 0 || age > 150 {
		return nil, ErrInvalidAge
	}

	user := User{
		Name:  name,
		Email: email,
		Age:   age,
	}

	err := s.repo.Create(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to add user: %w", err)
	}

	return &user, nil
}

// GetUser возвращает пользователя по ID
func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
	// TODO: реализуй метод
	// Просто делегируй в репозиторий
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// ListUsers возвращает всех пользователей
func (s *UserService) ListUsers(ctx context.Context) ([]*User, error) {
	// TODO: реализуй метод
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	return users, nil
}

// UpdateUser обновляет пользователя
func (s *UserService) UpdateUser(ctx context.Context, id int, name, email string, age int) (*User, error) {
	// TODO: реализуй метод
	// 1. Получи пользователя по ID
	// 2. Валидируй входные данные
	// 3. Обнови поля
	// 4. Вызови repo.Update
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if strings.TrimSpace(name) == "" {
		return nil, ErrNameRequired
	}
	if strings.TrimSpace(email) == "" {
		return nil, ErrEmailRequired
	}
	if age < 0 || age > 150 {
		return nil, ErrInvalidAge
	}

	user.Name = name
	user.Email = email
	user.Age = age

	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update: %w", err)
	}

	return user, nil
}

// DeleteUser удаляет пользователя
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	// TODO: реализуй метод
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	return nil
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

}
