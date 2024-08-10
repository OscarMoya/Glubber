package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

// Repository Interface, defines all the base ops for the repository
type Repository interface {
	CreateTable(ctx context.Context, createStatement string) error
	BeginTransaction(context.Context) (Transaction, error)
	Listen(ctx context.Context, channel string) error
	Notifications(context.Context) <-chan *pq.Notification
	CloseListener(context.Context) error
}

// Transaction defines the base ops for a transaction
type Transaction interface {
	Commit(context.Context) error
	Rollback(context.Context) error
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

// DBRepository represent a PostgreSQL repository.
type DBRepository struct {
	db       *sql.DB
	listener *pq.Listener
}

func NewDBRepository(connStr string, channel string) (Repository, error) {
	/*
		connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
			user, password, dbname, host, port, sslmode)
	*/

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	listener := pq.NewListener(connStr, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("Listener error: %v", err)
		}
	})

	err = listener.Listen(channel)
	if err != nil {
		return nil, fmt.Errorf("failed to listen to channel: %w", err)
	}

	if err := listener.Ping(); err != nil {
		return nil, fmt.Errorf("failed to initialize listener: %w", err)
	}

	return &DBRepository{db: db, listener: listener}, nil
}

func (repo *DBRepository) CreateTable(ctx context.Context, createStmt string) error {
	_, err := repo.db.ExecContext(ctx, createStmt)
	return err
}

// BeginDBTransaction return new DB transacion implementing the Transaction interface
func (repo *DBRepository) BeginTransaction(ctx context.Context) (Transaction, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error starting DBTransaction: %w", err)
	}
	return &DBTransaction{tx: tx}, nil
}

// Start a new Listen request to a channel
func (repo *DBRepository) Listen(ctx context.Context, channel string) error {
	err := repo.listener.Listen(channel)
	if err != nil {
		return fmt.Errorf("error listening to channel %s: %w", channel, err)
	}

	log.Printf("Listening on channel: %s", channel)
	return nil
}

// Creates new notification channel and returns it
func (repo *DBRepository) Notifications(ctx context.Context) <-chan *pq.Notification {
	notificationChan := make(chan *pq.Notification)

	go func() {
		defer close(notificationChan)
		for {
			select {
			case <-ctx.Done():
				return
			case notification := <-repo.listener.Notify:
				notificationChan <- notification
			case <-time.After(90 * time.Second):
				go repo.listener.Ping()
				log.Println("Ping to listener")
			}
		}
	}()

	return notificationChan
}

// CloseListener closes
func (repo *DBRepository) CloseListener(ctx context.Context) error {
	if repo.listener != nil {
		err := repo.listener.Close()
		if err != nil {
			return fmt.Errorf("error closing listener: %w", err)
		}
	}
	return nil
}

// DBTransaction represents a PostgreSQL transaction.
type DBTransaction struct {
	tx *sql.Tx
}

// Commit confirms the transaction.
func (t *DBTransaction) Commit(ctx context.Context) error {
	err := t.tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing DBTransaction: %w", err)
	}
	return nil
}

// Rollback rollback the transaction.
func (t *DBTransaction) Rollback(ctx context.Context) error {
	err := t.tx.Rollback()
	if err != nil {
		if err != sql.ErrTxDone { // ignore if transaction is already closed
			return fmt.Errorf("error rolling back DBTransaction: %w", err)
		}
	}
	return nil
}

// Exec executes a query within the transaction.
func (t *DBTransaction) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	result, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	return result, nil
}

// QueryRow executes a query that returns a single row within the transaction.
func (t *DBTransaction) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

// Query executes a query that returns multiple rows within the transaction.
func (t *DBTransaction) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	return rows, nil
}
