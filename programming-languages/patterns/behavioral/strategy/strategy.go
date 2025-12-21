package dbstrategy

import (
	"database/sql"
	"fmt"
)

type DBClient struct {
	master  *sql.DB
	replica *sql.DB
	strat   RouteStrategy
}

func NewDBClient(master, replica *sql.DB, s RouteStrategy) *DBClient {
	return &DBClient{
		master:  master,
		replica: replica,
		strat:   s,
	}
}

func (c *DBClient) SetStrategy(s RouteStrategy) {
	c.strat = s
}

// Query/Exec используют текущую стратегию выбора подключения.
func (c *DBClient) Query(query string, args ...any) (*sql.Rows, error) {
	db := c.strat.ChooseDB(OperationRead, c.master, c.replica)
	return db.Query(query, args...)
}

func (c *DBClient) Exec(query string, args ...any) (sql.Result, error) {
	db := c.strat.ChooseDB(OperationWrite, c.master, c.replica)
	return db.Exec(query, args...)
}

// Стратегии

type Operation int

const (
	OperationRead Operation = iota
	OperationWrite
)

type RouteStrategy interface {
	ChooseDB(op Operation, master, replica *sql.DB) *sql.DB
}

// Всегда ходим на мастер
type MasterOnlyStrategy struct{}

func (MasterOnlyStrategy) ChooseDB(_ Operation, master, _ *sql.DB) *sql.DB {
	return master
}

// Пишем в мастер, читаем из реплики
type MasterReplicaStrategy struct{}

func (MasterReplicaStrategy) ChooseDB(op Operation, master, replica *sql.DB) *sql.DB {
	if op == OperationRead && replica != nil {
		return replica
	}

	return master
}

// Можно добавить стратегию с фолбэком, health‑чеками и т.п.
type SafeReplicaStrategy struct{}

func (SafeReplicaStrategy) ChooseDB(op Operation, master, replica *sql.DB) *sql.DB {
	if op == OperationWrite || replica == nil {
		return master
	}

	return replica
}

func Example() error {
	// masterDB, replicaDB := sql.Open(...)

	var masterDB, replicaDB *sql.DB // заглушки

	client := NewDBClient(masterDB, replicaDB, MasterReplicaStrategy{})

	// read -> реплика
	if _, err := client.Query("SELECT * FROM users WHERE id = ?", 1); err != nil {
		return fmt.Errorf("read: %w", err)
	}

	// write -> мастер
	if _, err := client.Exec("UPDATE users SET name = ? WHERE id = ?", "Bob", 1); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	// В рантайме можно сменить стратегию
	client.SetStrategy(MasterOnlyStrategy{})

	return nil
}
