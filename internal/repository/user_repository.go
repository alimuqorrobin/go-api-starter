package repository

import (
    "context"
    "database/sql"
    "fmt"
    "time"

    "go-api-starter/internal/domain"
)

type userRepository struct {
    db *sql.DB
}

// NewUserRepository creates new user repository
func NewUserRepository(db *sql.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
    query := `
        INSERT INTO users (username, email, password_hash, full_name, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at
    `

    now := time.Now()
    err := r.db.QueryRowContext(
        ctx,
        query,
        user.Username,
        user.Email,
        user.PasswordHash,
        user.FullName,
        user.IsActive,
        now,
        now,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    return user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
    query := `
        SELECT id, username, email, password_hash, full_name, is_active, created_at, updated_at
        FROM users
        WHERE id = $1
    `

    user := &domain.User{}
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.FullName,
        &user.IsActive,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, domain.ErrUserNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }

    return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
    query := `
        SELECT id, username, email, password_hash, full_name, is_active, created_at, updated_at
        FROM users
        WHERE email = $1
    `

    user := &domain.User{}
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.FullName,
        &user.IsActive,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, domain.ErrUserNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }

    return user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
    query := `
        SELECT id, username, email, password_hash, full_name, is_active, created_at, updated_at
        FROM users
        WHERE username = $1
    `

    user := &domain.User{}
    err := r.db.QueryRowContext(ctx, query, username).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
        &user.FullName,
        &user.IsActive,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    if err == sql.ErrNoRows {
        return nil, domain.ErrUserNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }

    return user, nil
}

func (r *userRepository) FindAll(ctx context.Context, limit, offset int) ([]*domain.User, error) {
    query := `
        SELECT id, username, email, password_hash, full_name, is_active, created_at, updated_at
        FROM users
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

    rows, err := r.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("failed to query users: %w", err)
    }
    defer rows.Close()

    var users []*domain.User
    for rows.Next() {
        user := &domain.User{}
        err := rows.Scan(
            &user.ID,
            &user.Username,
            &user.Email,
            &user.PasswordHash,
            &user.FullName,
            &user.IsActive,
            &user.CreatedAt,
            &user.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan user: %w", err)
        }
        users = append(users, user)
    }

    return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
    query := `
        UPDATE users
        SET username = $1, email = $2, full_name = $3, is_active = $4, updated_at = $5
        WHERE id = $6
        RETURNING updated_at
    `

    now := time.Now()
    err := r.db.QueryRowContext(
        ctx,
        query,
        user.Username,
        user.Email,
        user.FullName,
        user.IsActive,
        now,
        user.ID,
    ).Scan(&user.UpdatedAt)

    if err != nil {
        return nil, fmt.Errorf("failed to update user: %w", err)
    }

    return user, nil
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
    query := `DELETE FROM users WHERE id = $1`

    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }

    if rowsAffected == 0 {
        return domain.ErrUserNotFound
    }

    return nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
    query := `SELECT COUNT(*) FROM users`

    var count int64
    err := r.db.QueryRowContext(ctx, query).Scan(&count)
    if err != nil {
        return 0, fmt.Errorf("failed to count users: %w", err)
    }

    return count, nil
}
```

**Penjelasan:**
- Implementasi interface `UserRepository`
- Semua query menggunakan **raw SQL** (tidak pakai ORM)
- Context-aware untuk cancellation
- Error handling yang proper
- `$1, $2` = PostgreSQL placeholders (untuk MySQL pakai `?`)

**Kenapa Raw SQL?**
```
// ✅ Performance maksimal
// ✅ Full control atas query
// ✅ Tidak tergantung library ORM
// ✅ Easy to optimize (add index, etc)
// ✅ Ganti database mudah (tinggal ganti query syntax)