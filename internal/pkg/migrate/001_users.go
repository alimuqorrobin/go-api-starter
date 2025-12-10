package migrate

import (
    "database/sql"
    "time"
)

func init() {
    Register(&Migration{
        Version: 1,
        Name: "create_users_table",
        Created: time.Now(),
        Up: upCreateUsers,
        Down: downCreateUsers,
    })
}

func upCreateUsers(tx *sql.Tx) error {
    _, err := tx.Exec(`
CREATE TABLE IF NOT EXISTS users (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  username VARCHAR(100) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB;
`)
    return err
}

func downCreateUsers(tx *sql.Tx) error {
    _, err := tx.Exec(`DROP TABLE IF EXISTS users;`)
    return err
}
