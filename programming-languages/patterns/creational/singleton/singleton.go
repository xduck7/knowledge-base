package singleton

import "sync"

type DB struct {
	DSN string
}

var (
	instance *DB
	once     sync.Once
)

func Instance() *DB {
	once.Do(func() {
		instance = &DB{
			DSN: "postgres://user:pass@localhost:5432/db",
		}
	})

	return instance
}
