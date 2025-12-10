package migrate

import (
    "database/sql"
    "fmt"
    "sort"
    "time"
)

type Migration struct {
    Version int
    Name string
    Up func(tx *sql.Tx) error
    Down func(tx *sql.Tx) error
    Created time.Time
}

var registry = map[int]*Migration{}

func Register(m *Migration) {
    if _, exists := registry[m.Version]; exists {
        panic(fmt.Sprintf("migration version %d already registered", m.Version))
    }
    registry[m.Version] = m
}

func List() []*Migration {
    vers := make([]int, 0, len(registry))
    for v := range registry {
        vers = append(vers, v)
    }
    sort.Ints(vers)
    out := make([]*Migration, 0, len(vers))
    for _, v := range vers {
        out = append(out, registry[v])
    }
    return out
}

func EnsureSchema(db *sql.DB) error {
    const q = `
CREATE TABLE IF NOT EXISTS schema_migrations (
  version INT NOT NULL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  applied_at TIMESTAMP NOT NULL
);`
    _, err := db.Exec(q)
    return err
}

func AppliedVersions(db *sql.DB) (map[int]time.Time, error) {
    rows, err := db.Query("SELECT version, applied_at FROM schema_migrations")
    if err != nil {
        return map[int]time.Time{}, nil
    }
    defer rows.Close()
    out := map[int]time.Time{}
    for rows.Next() {
        var v int
        var t time.Time
        if err := rows.Scan(&v, &t); err != nil {
            return nil, err
        }
        out[v] = t
    }
    return out, nil
}

func ApplyPending(db *sql.DB) error {
    if err := EnsureSchema(db); err != nil { return err }
    applied, err := AppliedVersions(db)
    if err != nil { return err }
    migs := List()
    for _, m := range migs {
        if _, ok := applied[m.Version]; ok { continue }
        tx, err := db.Begin()
        if err != nil { return err }
        if err := m.Up(tx); err != nil {
            _ = tx.Rollback()
            return fmt.Errorf("migration %d up failed: %w", m.Version, err)
        }
        if _, err := tx.Exec("INSERT INTO schema_migrations(version,name,applied_at) VALUES(?,?,?)", m.Version, m.Name, time.Now().UTC()); err != nil {
            _ = tx.Rollback()
            return fmt.Errorf("failed to write schema_migrations: %w", err)
        }
        if err := tx.Commit(); err != nil {
            return fmt.Errorf("commit failed for migration %d: %w", m.Version, err)
        }
    }
    return nil
}

func RollbackLast(db *sql.DB) error {
    if err := EnsureSchema(db); err != nil { return err }
    applied, err := AppliedVersions(db)
    if err != nil { return err }
    if len(applied) == 0 { return fmt.Errorf("no applied migrations") }
    max := -1
    for v := range applied { if v > max { max = v } }
    m, ok := registry[max]
    if !ok { return fmt.Errorf("migration %d not found in registry", max) }
    tx, err := db.Begin()
    if err != nil { return err }
    if err := m.Down(tx); err != nil {
        _ = tx.Rollback()
        return fmt.Errorf("migration %d down failed: %w", m.Version, err)
    }
    if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = ?", m.Version); err != nil {
        _ = tx.Rollback()
        return err
    }
    if err := tx.Commit(); err != nil { return err }
    return nil
}
