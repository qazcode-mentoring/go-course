package main

import (
	"context"
	"errors"
	"fmt"
	"log"
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
)

// User представляет пользователя в базе данных
type User struct {
	ID        int
	Name      string
	Email     string
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// String возвращает строковое представление пользователя
func (u User) String() string {
	return fmt.Sprintf("User{ID: %d, Name: %s, Email: %s, Age: %d}", u.ID, u.Name, u.Email, u.Age)
}

// CreateUser создаёт нового пользователя в базе данных
func CreateUser(ctx context.Context, pool *pgxpool.Pool, name, email string, age int) (User, error) {
	// TODO: реализуй функцию
	// 1. Выполни INSERT запрос с RETURNING:
	//    query := `INSERT INTO users (name, email, age)
	//              VALUES ($1, $2, $3)
	//              RETURNING id, name, email, age, created_at, updated_at`
	// 2. Используй QueryRow для выполнения и Scan для получения результата
	// 3. Обработай ошибку уникальности email (код 23505):
	//    var pgErr *pgconn.PgError
	//    if errors.As(err, &pgErr) && pgErr.Code == "23505" {
	//        return User{}, ErrEmailAlreadyExists
	//    }
	// 4. Верни созданного пользователя
	query := `INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id, name, email, age, created_at, updated_at`
	var user User
	err := pool.QueryRow(ctx, query, name, email, age).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return User{}, ErrEmailAlreadyExists
	}

	return user, nil
}

// GetUserByID возвращает пользователя по ID
func GetUserByID(ctx context.Context, pool *pgxpool.Pool, id int) (User, error) {
	// TODO: реализуй функцию
	// 1. Выполни SELECT запрос:
	//    query := `SELECT id, name, email, age, created_at, updated_at
	//              FROM users WHERE id = $1`
	// 2. Используй QueryRow и Scan
	// 3. Проверь на pgx.ErrNoRows:
	//    if errors.Is(err, pgx.ErrNoRows) {
	//        return User{}, ErrUserNotFound
	//    }
	// 4. Верни пользователя
	query := `SELECT id, name, email, age, created_at, updated_at FROM users WHERE id = $1`
	var user User
	err := pool.QueryRow(ctx, query, id).Scan(
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

	return user, nil
}

// GetUserByEmail возвращает пользователя по email
func GetUserByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (User, error) {
	// TODO: реализуй функцию
	// Аналогично GetUserByID, но поиск по email
	query := `SELECT id, name, email, age, created_at, updated_at FROM users WHERE email = $1`
	var user User
	err := pool.QueryRow(ctx, query, email).Scan(
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

	return user, nil
}

// GetAllUsers возвращает всех пользователей
func GetAllUsers(ctx context.Context, pool *pgxpool.Pool) ([]User, error) {
	// TODO: реализуй функцию
	// 1. Выполни SELECT запрос:
	//    query := `SELECT id, name, email, age, created_at, updated_at
	//              FROM users ORDER BY id`
	// 2. Используй Query для получения rows:
	//    rows, err := pool.Query(ctx, query)
	// 3. Не забудь defer rows.Close()
	// 4. Итерируй по rows с помощью rows.Next():
	//    var users []User
	//    for rows.Next() {
	//        var user User
	//        err := rows.Scan(&user.ID, &user.Name, ...)
	//        if err != nil {
	//            return nil, err
	//        }
	//        users = append(users, user)
	//    }
	// 5. Проверь rows.Err() после цикла
	// 6. Верни users
	query := `SELECT id, name, email, age, created_at, updated_at FROM users ORDER BY id`
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
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
		return nil, fmt.Errorf("fail: %w", rows.Err())
	}

	return users, nil
}

// UpdateUser обновляет данные пользователя
func UpdateUser(ctx context.Context, pool *pgxpool.Pool, id int, name, email string, age int) (User, error) {
	// TODO: реализуй функцию
	// 1. Выполни UPDATE запрос с RETURNING:
	//    query := `UPDATE users SET name = $1, email = $2, age = $3, updated_at = NOW()
	//              WHERE id = $4
	//              RETURNING id, name, email, age, created_at, updated_at`
	// 2. Используй QueryRow и Scan
	// 3. Проверь на pgx.ErrNoRows (пользователь не найден)
	// 4. Проверь на дублирование email (код 23505)
	// 5. Верни обновлённого пользователя
	query := `UPDATE users SET name = $1, email = $2, age = $3, updated_at = now() WHERE id = $4
				RETURNING id, name, email, age, created_at, updated_at`
	var updatedUser User
	err := pool.QueryRow(ctx, query, name, email, age, id).Scan(
		&updatedUser.ID,
		&updatedUser.Name,
		&updatedUser.Email,
		&updatedUser.Age,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrUserNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return User{}, ErrEmailAlreadyExists
	}

	return updatedUser, nil
}

// DeleteUser удаляет пользователя по ID
func DeleteUser(ctx context.Context, pool *pgxpool.Pool, id int) error {
	// TODO: реализуй функцию
	// 1. Выполни DELETE запрос:
	//    result, err := pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	// 2. Проверь ошибку
	// 3. Проверь количество затронутых строк:
	//    if result.RowsAffected() == 0 {
	//        return ErrUserNotFound
	//    }
	// 4. Верни nil при успехе
	result, err := pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

func main() {
	fmt.Println("=== PostgreSQL CRUD Demo ===")

	ctx := context.Background()

	// Подключаемся к базе данных
	pool, err := pgxpool.New(ctx, defaultDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Проверяем подключение
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("Connected to PostgreSQL!")

	// --- Create User ---
	fmt.Println("\n--- Create User ---")
	user, err := CreateUser(ctx, pool, "John Doe", "john@example.com", 28)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
	} else {
		fmt.Printf("Created user: %s\n", user)
	}

	// --- Get User by ID ---
	fmt.Println("\n--- Get User by ID ---")
	if user.ID > 0 {
		foundUser, err := GetUserByID(ctx, pool, user.ID)
		if err != nil {
			log.Printf("Failed to get user: %v", err)
		} else {
			fmt.Printf("Found user: %s\n", foundUser)
		}
	}

	// --- Get User by Email ---
	fmt.Println("\n--- Get User by Email ---")
	foundByEmail, err := GetUserByEmail(ctx, pool, "john@example.com")
	if err != nil {
		log.Printf("Failed to get user by email: %v", err)
	} else {
		fmt.Printf("Found user: %s\n", foundByEmail)
	}

	// --- Update User ---
	fmt.Println("\n--- Update User ---")
	if user.ID > 0 {
		updatedUser, err := UpdateUser(ctx, pool, user.ID, "John Updated", "john.updated@example.com", 29)
		if err != nil {
			log.Printf("Failed to update user: %v", err)
		} else {
			fmt.Printf("Updated user: %s\n", updatedUser)
		}
	}

	// --- Get All Users ---
	fmt.Println("\n--- Get All Users ---")
	users, err := GetAllUsers(ctx, pool)
	if err != nil {
		log.Printf("Failed to get all users: %v", err)
	} else {
		fmt.Printf("Total users: %d\n", len(users))
		for _, u := range users {
			fmt.Printf("  - %s (%s)\n", u.Name, u.Email)
		}
	}

	// --- Delete User ---
	fmt.Println("\n--- Delete User ---")
	if user.ID > 0 {
		err := DeleteUser(ctx, pool, user.ID)
		if err != nil {
			log.Printf("Failed to delete user: %v", err)
		} else {
			fmt.Println("User deleted successfully")
		}
	}

	// --- Verify Deletion ---
	fmt.Println("\n--- Verify Deletion ---")
	if user.ID > 0 {
		_, err := GetUserByID(ctx, pool, user.ID)
		if errors.Is(err, ErrUserNotFound) {
			fmt.Println("GetUserByID: user not found (expected)")
		} else if err != nil {
			log.Printf("Unexpected error: %v", err)
		} else {
			fmt.Println("User still exists (unexpected)")
		}
	}

}
