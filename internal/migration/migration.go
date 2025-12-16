package migration

import (
	"fmt"

	"go-api-starter/internal/migration/migrations"
	"gorm.io/gorm"
)

type Migration interface {
	Name() string
	Up(db *gorm.DB) error
}

type Migrator struct {
	db         *gorm.DB
	migrations []Migration
}

func NewMigrator(db *gorm.DB) *Migrator {
	m := &Migrator{
		db:         db,
		migrations: make([]Migration, 0),
	}

	// ============ REGISTER MIGRATIONS ============
	// PENTING: Urutan harus dari yang paling awal
	// Jika tambah migration baru, append di akhir

	m.migrations = append(m.migrations, migrations.NewCreateUsersTable())
	m.migrations = append(m.migrations, migrations.NewCreateProductsTable())
	m.migrations = append(m.migrations, migrations.NewCreateJWTTokensTable())

	// TODO: Untuk menambah table baru:
	// 1. Create file: internal/migration/migrations/004_create_<table_name>.go
	// 2. Implement Migration interface
	// 3. Append di sini: m.migrations = append(m.migrations, migrations.NewCreate<TableName>Table())

	return m
}

func (m *Migrator) Migrate() error {
	// Create migrations table jika belum ada
	type MigrationRecord struct {
		Name string `gorm:"primaryKey"`
	}

	if err := m.db.AutoMigrate(&MigrationRecord{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Run each migration
	for _, migration := range m.migrations {
		var record MigrationRecord

		// Check if migration sudah pernah dijalankan
		result := m.db.Where("name = ?", migration.Name()).First(&record)

		if result.Error == gorm.ErrRecordNotFound {
			// Migration belum dijalankan, jalankan sekarang
			fmt.Printf("Running migration: %s\n", migration.Name())

			if err := migration.Up(m.db); err != nil {
				return fmt.Errorf("migration %s failed: %v", migration.Name(), err)
			}

			// Record bahwa migration sudah dijalankan
			m.db.Create(&MigrationRecord{Name: migration.Name()})

			fmt.Printf("✅ Migration completed: %s\n", migration.Name())
		}
	}

	return nil
}
