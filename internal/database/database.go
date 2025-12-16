// ============================================
// FILE: internal/database/database.go
// COMPLETE FIX - Remove ALL unused imports
// ============================================

package database

import (
	"fmt"
	"sync"
	"time"

	"go-api-starter/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database interface - untuk decoupling
type Database interface {
	GetDB() *gorm.DB
	Close() error
	IsConnected() bool
	Ping() error
}

type database struct {
	db *gorm.DB
	mu sync.RWMutex
}

var (
	instance Database
	once     sync.Once
	dbLock   sync.Mutex
)

// NewDatabase - singleton pattern dengan error recovery
func NewDatabase(cfg *config.Config) (Database, error) {
	var err error
	once.Do(func() {
		dbLock.Lock()
		defer dbLock.Unlock()

		// Tentukan driver
		var dialector gorm.Dialector

		switch cfg.DBDriver {
		case "mysql":
			dsn := config.GetDSN(cfg)
			dialector = mysql.Open(dsn)
		case "postgres":
			dsn := config.GetDSN(cfg)
			dialector = postgres.Open(dsn)
		default:
			err = fmt.Errorf("unsupported database driver: %s", cfg.DBDriver)
			return
		}

		db, err := gorm.Open(dialector, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			return
		}

		sqlDB, err := db.DB()
		if err != nil {
			return
		}

		// Connection pool
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(time.Hour)

		// Test connection
		if err = sqlDB.Ping(); err != nil {
			return
		}

		instance = &database{db: db}
	})

	return instance, err
}

func (d *database) GetDB() *gorm.DB {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.db
}

func (d *database) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.db != nil {
		sqlDB, err := d.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func (d *database) IsConnected() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.db == nil {
		return false
	}

	sqlDB, err := d.db.DB()
	if err != nil {
		return false
	}

	return sqlDB.Ping() == nil
}

func (d *database) Ping() error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.db == nil {
		return fmt.Errorf("database not initialized")
	}

	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}

// AutoMigrate helper
func (d *database) AutoMigrate(models ...interface{}) error {
	return d.db.AutoMigrate(models...)
}

// Untuk injeksi ke repository
type DatabaseProvider interface {
	GetDatabase() Database
}

type databaseProvider struct {
	db Database
}

func NewDatabaseProvider(db Database) DatabaseProvider {
	return &databaseProvider{db: db}
}

func (p *databaseProvider) GetDatabase() Database {
	return p.db
}

// Batch operation helper untuk database abstraction
type BatchOperation struct {
	db Database
}

func (b *BatchOperation) Exec(fn func(*gorm.DB) error) error {
	return fn(b.db.GetDB())
}
