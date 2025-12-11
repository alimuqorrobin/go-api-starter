.package database

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    "golang-api-starter/config"
)

// Database interface untuk abstraksi - tidak tergantung ORM tertentu
type Database interface {
    GetDB() *sql.DB
    Close() error
    Ping() error
}

// Connection struct
type Connection struct {
    db *sql.DB
}

// NewConnection membuat koneksi database berdasarkan driver
func NewConnection(cfg *config.Config) (Database, error) {
    var (
        db  *sql.DB
        err error
    )

    switch cfg.DBDriver {
    case "postgres":
        db, err = NewPostgresConnection(cfg)
    case "mysql":
        db, err = NewMySQLConnection(cfg)
    default:
        return nil, fmt.Errorf("unsupported database driver: %s", cfg.DBDriver)
    }

    if err != nil {
        return nil, err
    }

    // Set connection pool settings
    db.SetMaxOpenConns(cfg.DBMaxOpenConns)
    db.SetMaxIdleConns(cfg.DBMaxIdleConns)
    db.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

    // Verify connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := db.PingContext(ctx); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    return &Connection{db: db}, nil
}

func (c *Connection) GetDB() *sql.DB {
    return c.db
}

func (c *Connection) Close() error {
    return c.db.Close()
}

func (c *Connection) Ping() error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    return c.db.PingContext(ctx)
}
```

**Penjelasan:**
- `Database interface` = Contract untuk semua database driver
- `NewConnection()` = Factory pattern, pilih driver berdasarkan config
- Connection pooling = Optimasi performa database
- Ping verification = Pastikan koneksi sukses

**Keuntungan Abstraksi:**
```
`Ganti dari PostgreSQL ke MySQL?
``1. Ubah DB_DRIVER=mysql di .env
``2. Done! Tidak perlu ubah code lain