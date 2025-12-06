package repository

import (
    "database/sql"
)

type User struct {
    ID       int64  `json:"id"`
    Username string `json:"username"`
    Password string `json:"-"`
}

type UserRepository interface {
    FindByID(id int64) (*User, error)
    FindByUsername(username string) (*User, error)
    Create(u *User) (int64, error)
}

type userRepoSQL struct {
    db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepository {
    return &userRepoSQL{db: db}
}

func (r *userRepoSQL) FindByID(id int64) (*User, error) {
    u := &User{}
    err := r.db.QueryRow("SELECT id, username, password_hash FROM users WHERE id = ?", id).Scan(&u.ID, &u.Username, &u.Password)
    if err != nil {
        return nil, err
    }
    return u, nil
}

func (r *userRepoSQL) FindByUsername(username string) (*User, error) {
    u := &User{}
    err := r.db.QueryRow("SELECT id, username, password_hash FROM users WHERE username = ?", username).Scan(&u.ID, &u.Username, &u.Password)
    if err != nil {
        return nil, err
    }
    return u, nil
}

func (r *userRepoSQL) Create(u *User) (int64, error) {
    res, err := r.db.Exec("INSERT INTO users(username, password_hash) VALUES(?, ?)", u.Username, u.Password)
    if err != nil {
        return 0, err
    }
    return res.LastInsertId()
}
