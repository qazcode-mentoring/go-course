package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DSN для подключения к локальному PostgreSQL в Docker
const defaultDSN = "postgresql://gouser:gopassword@localhost:5432/gocourse?sslmode=disable"

// Ошибки приложения
var (
	ErrAccountNotFound   = errors.New("account not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrSameAccount       = errors.New("cannot transfer to same account")
)

// Account представляет банковский счёт
type Account struct {
	ID        int
	UserID    int
	Balance   float64
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// String возвращает строковое представление счёта
func (a Account) String() string {
	return fmt.Sprintf("Account{ID: %d, UserID: %d, Balance: %.2f %s}", a.ID, a.UserID, a.Balance, a.Currency)
}

// Transaction представляет денежную транзакцию
type Transaction struct {
	ID            int
	FromAccountID int
	ToAccountID   int
	Amount        float64
	Description   string
	CreatedAt     time.Time
}

// String возвращает строковое представление транзакции
func (t Transaction) String() string {
	return fmt.Sprintf("Transaction{ID: %d, From: %d, To: %d, Amount: %.2f}", t.ID, t.FromAccountID, t.ToAccountID, t.Amount)
}

// GetAccountByID возвращает счёт по ID
func GetAccountByID(ctx context.Context, pool *pgxpool.Pool, id int) (Account, error) {
	// TODO: реализуй функцию
	// 1. Выполни SELECT запрос:
	//    query := `SELECT id, user_id, balance, currency, created_at, updated_at
	//              FROM accounts WHERE id = $1`
	// 2. Используй QueryRow и Scan
	// 3. Проверь на pgx.ErrNoRows
	// 4. Верни счёт или ошибку
	query := `SELECT id, user_id, balance, currency, created_at, updated_at FROM accounts WHERE id = $1`
	var account Account
	err := pool.QueryRow(ctx, query, id).Scan(
		&account.ID,
		&account.UserID,
		&account.Balance,
		&account.Currency,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Account{}, ErrAccountNotFound
		}
		return Account{}, err
	}

	return account, nil
}

// TransferMoney выполняет перевод денег между счетами в транзакции
func TransferMoney(ctx context.Context, pool *pgxpool.Pool, fromID, toID int, amount float64, description string) error {
	// TODO: реализуй функцию
	//
	// 1. Проверь, что fromID != toID (нельзя переводить на тот же счёт):
	//    if fromID == toID {
	//        return ErrSameAccount
	//    }
	//
	// 2. Начни транзакцию:
	//    tx, err := pool.Begin(ctx)
	//    if err != nil {
	//        return fmt.Errorf("begin transaction: %w", err)
	//    }
	//    defer tx.Rollback(ctx) // безопасный откат
	//
	// 3. Списание с исходного счёта (с проверкой баланса):
	//    result, err := tx.Exec(ctx,
	//        `UPDATE accounts SET balance = balance - $1, updated_at = NOW()
	//         WHERE id = $2 AND balance >= $1`,
	//        amount, fromID,
	//    )
	//    if err != nil {
	//        return fmt.Errorf("debit account: %w", err)
	//    }
	//    if result.RowsAffected() == 0 {
	//        return ErrInsufficientFunds // или ErrAccountNotFound
	//    }
	//
	// 4. Зачисление на целевой счёт:
	//    result, err = tx.Exec(ctx,
	//        `UPDATE accounts SET balance = balance + $1, updated_at = NOW()
	//         WHERE id = $2`,
	//        amount, toID,
	//    )
	//    if err != nil {
	//        return fmt.Errorf("credit account: %w", err)
	//    }
	//    if result.RowsAffected() == 0 {
	//        return ErrAccountNotFound
	//    }
	//
	// 5. Запись транзакции в историю:
	//    _, err = tx.Exec(ctx,
	//        `INSERT INTO transactions (from_account_id, to_account_id, amount, description)
	//         VALUES ($1, $2, $3, $4)`,
	//        fromID, toID, amount, description,
	//    )
	//    if err != nil {
	//        return fmt.Errorf("record transaction: %w", err)
	//    }
	//
	// 6. Зафиксируй транзакцию:
	//    if err := tx.Commit(ctx); err != nil {
	//        return fmt.Errorf("commit transaction: %w", err)
	//    }
	//
	// 7. Верни nil при успехе
	if toID == fromID {
		return ErrSameAccount
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx,
		`UPDATE accounts SET balance = balance - $1, updated_at = NOW() 
                WHERE id = $2 AND balance >= $1`, amount, fromID)
	if err != nil {
		return fmt.Errorf("debit account: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrInsufficientFunds
	}

	result, err = tx.Exec(ctx,
		`UPDATE accounts SET balance = balance + $1, updated_at = NOW()
                WHERE id = $2`, amount, toID)
	if err != nil {
		return fmt.Errorf("credit account: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrAccountNotFound
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO transactions (from_account_id, to_account_id, amount, description)
			VALUES ($1, $2, $3, $4)`, fromID, toID, amount, description)

	if err != nil {
		return fmt.Errorf("record transaction: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// GetAccountTransactions возвращает историю транзакций для счёта
func GetAccountTransactions(ctx context.Context, pool *pgxpool.Pool, accountID int) ([]Transaction, error) {
	// TODO: реализуй функцию
	// 1. Выполни SELECT запрос:
	//    query := `SELECT id, from_account_id, to_account_id, amount, description, created_at
	//              FROM transactions
	//              WHERE from_account_id = $1 OR to_account_id = $1
	//              ORDER BY created_at DESC`
	// 2. Используй Query и итерируй по rows
	// 3. Не забудь rows.Close() и проверку rows.Err()
	// 4. Верни список транзакций
	query := `SELECT id, from_account_id, to_account_id, amount, description, created_at
				FROM transactions
				WHERE from_account_id = $1 OR to_account_id = $1
				ORDER BY created_at DESC`
	rows, err := pool.Query(ctx, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		err = rows.Scan(
			&transaction.ID,
			&transaction.FromAccountID,
			&transaction.ToAccountID,
			&transaction.Amount,
			&transaction.Description,
			&transaction.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		transactions = append(transactions, transaction)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("fail: %w", rows.Err())
	}

	return transactions, nil
}

// WithTransaction выполняет функцию внутри транзакции
func WithTransaction(ctx context.Context, pool *pgxpool.Pool, fn func(tx pgx.Tx) error) error {
	// TODO: реализуй функцию
	// 1. Начни транзакцию: tx, err := pool.Begin(ctx)
	// 2. Добавь defer tx.Rollback(ctx)
	// 3. Выполни функцию: if err := fn(tx); err != nil { return err }
	// 4. Зафиксируй: return tx.Commit(ctx)
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err = fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func main() {
	fmt.Println("=== PostgreSQL Transactions Demo ===")

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

	// --- Начальные балансы ---
	fmt.Println("\n--- Initial Balances ---")
	acc1, err := GetAccountByID(ctx, pool, 1)
	if err != nil {
		log.Printf("Failed to get account 1: %v", err)
	} else {
		fmt.Printf("Account 1: %.2f %s\n", acc1.Balance, acc1.Currency)
	}

	acc2, err := GetAccountByID(ctx, pool, 2)
	if err != nil {
		log.Printf("Failed to get account 2: %v", err)
	} else {
		fmt.Printf("Account 2: %.2f %s\n", acc2.Balance, acc2.Currency)
	}

	// --- Перевод денег ---
	fmt.Println("\n--- Transfer Money ---")
	fmt.Println("Transferring 100.00 from account 1 to account 2...")

	err = TransferMoney(ctx, pool, 1, 2, 100.00, "Payment for services")
	if err != nil {
		log.Printf("Transfer failed: %v", err)
	} else {
		fmt.Println("Transfer successful!")
	}

	// --- Балансы после перевода ---
	fmt.Println("\n--- Balances After Transfer ---")
	acc1, err = GetAccountByID(ctx, pool, 1)
	if err != nil {
		log.Printf("Failed to get account 1: %v", err)
	} else {
		fmt.Printf("Account 1: %.2f %s\n", acc1.Balance, acc1.Currency)
	}

	acc2, err = GetAccountByID(ctx, pool, 2)
	if err != nil {
		log.Printf("Failed to get account 2: %v", err)
	} else {
		fmt.Printf("Account 2: %.2f %s\n", acc2.Balance, acc2.Currency)
	}

	// --- История транзакций ---
	fmt.Println("\n--- Transaction History ---")
	transactions, err := GetAccountTransactions(ctx, pool, 1)
	if err != nil {
		log.Printf("Failed to get transactions: %v", err)
	} else {
		fmt.Printf("Account 1 transactions: %d\n", len(transactions))
		for _, t := range transactions {
			if t.FromAccountID == 1 {
				fmt.Printf("  - Sent %.2f to account %d: %s\n", t.Amount, t.ToAccountID, t.Description)
			} else {
				fmt.Printf("  - Received %.2f from account %d: %s\n", t.Amount, t.FromAccountID, t.Description)
			}
		}
	}

	// --- Тест: Недостаточно средств ---
	fmt.Println("\n--- Test Insufficient Funds ---")
	fmt.Println("Trying to transfer 10000.00 (more than balance)...")

	err = TransferMoney(ctx, pool, 1, 2, 10000.00, "Large payment")
	if errors.Is(err, ErrInsufficientFunds) {
		fmt.Println("Error: insufficient funds (expected)")
	} else if err != nil {
		log.Printf("Unexpected error: %v", err)
	} else {
		fmt.Println("Transfer succeeded (unexpected)")
	}

	// --- Тест: Перевод на тот же счёт ---
	fmt.Println("\n--- Test Same Account ---")
	fmt.Println("Trying to transfer to same account...")

	err = TransferMoney(ctx, pool, 1, 1, 50.00, "Self transfer")
	if errors.Is(err, ErrSameAccount) {
		fmt.Println("Error: cannot transfer to same account (expected)")
	} else if err != nil {
		log.Printf("Unexpected error: %v", err)
	} else {
		fmt.Println("Transfer succeeded (unexpected)")
	}

}
