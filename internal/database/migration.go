package database

import (
    "database/sql"
    "fmt"
    "log"

    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationManager struct {
    db           *sql.DB
    migrationURL string
}

func NewMigrationManager(db *sql.DB, migrationPath string) *MigrationManager {
    return &MigrationManager{
        db:           db,
        migrationURL: migrationPath,
    }
}

func (m *MigrationManager) RunMigrations() error {
    driver, err := postgres.WithInstance(m.db, &postgres.Config{})
    if err != nil {
        return fmt.Errorf("could not create database driver: %w", err)
    }

    migration, err := migrate.NewWithDatabaseInstance(
        m.migrationURL,
        "postgres",
        driver,
    )
    if err != nil {
        return fmt.Errorf("could not create migrate instance: %w", err)
    }

    if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("migration failed: %w", err)
    }

    version, dirty, err := migration.Version()
    if err != nil && err != migrate.ErrNilVersion {
        return fmt.Errorf("could not get migration version: %w", err)
    }

    log.Printf("✅ Migration completed. Version: %d (dirty: %v)", version, dirty)
    return nil
}

func (m *MigrationManager) Rollback() error {
    driver, err := postgres.WithInstance(m.db, &postgres.Config{})
    if err != nil {
        return err
    }

    migration, err := migrate.NewWithDatabaseInstance(
        m.migrationURL,
        "postgres",
        driver,
    )
    if err != nil {
        return err
    }

    if err := migration.Steps(-1); err != nil {
        return err
    }

    log.Println("✅ Rollback completed successfully")
    return nil
}

func (m *MigrationManager) GetVersion() (uint, bool, error) {
    driver, err := postgres.WithInstance(m.db, &postgres.Config{})
    if err != nil {
        return 0, false, err
    }

    migration, err := migrate.NewWithDatabaseInstance(
        m.migrationURL,
        "postgres",
        driver,
    )
    if err != nil {
        return 0, false, err
    }

    return migration.Version()
}
```

**Penjelasan:**
- `RunMigrations()` = Jalankan semua migration yang belum dijalankan
- `Rollback()` = Rollback 1 migration terakhir
- `GetVersion()` = Cek versi migration saat ini
- Auto-run saat aplikasi start (di main.go)

**Cara Kerja Migration:**
```
// migrations/
// ├── 000001_create_users_table.up.sql    ← CREATE TABLE
// ├── 000001_create_users_table.down.sql  ← DROP TABLE
// ├── 000002_create_posts_table.up.sql
// └── 000002_create_posts_table.down.sql

// Schema migrations table:
// version | dirty
// --------|------
//    2    | false  ← Currently at version 2

// migrate up → Run all .up.sql yang belum dijalankan
// migrate down → Run .down.sql untuk rollback
```

---

### ✅ Part 2 Selesai!

Sudah selesai:
- ✅ internal/database/connection.go (Factory pattern)
- ✅ internal/database/postgres.go (PostgreSQL driver)
- ✅ internal/database/mysql.go (MySQL driver)
- ✅ internal/database/migration.go (Migration manager)

**Struktur sekarang:**
```
// golang-api-starter/
// ├── cmd/server/main.go
// ├── config/config.go
// ├── internal/
// │   └── database/
// │       ├── connection.go
// │       ├── postgres.go
// │       ├── mysql.go
// │       └── migration.go
// ├── go.mod
// ├── .env.example
// └── .gitignore