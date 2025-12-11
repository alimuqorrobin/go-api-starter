package repository

import (
    "context"
    "golang-api-starter/internal/domain"
)

// UserRepository interface - tidak tergantung pada ORM tertentu
type UserRepository interface {
    Create(ctx context.Context, user *domain.User) (*domain.User, error)
    FindByID(ctx context.Context, id int64) (*domain.User, error)
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
    FindByUsername(ctx context.Context, username string) (*domain.User, error)
    FindAll(ctx context.Context, limit, offset int) ([]*domain.User, error)
    Update(ctx context.Context, user *domain.User) (*domain.User, error)
    Delete(ctx context.Context, id int64) error
    Count(ctx context.Context) (int64, error)
}
```

**Penjelasan:**
- Interface untuk repository layer
- **PENTING**: Ini yang membuat kita tidak tergantung ORM!
- Implementasi bisa pakai raw SQL, GORM, atau apapun
- Ganti database? Implementasi ulang interface ini saja

**Contoh flexibilitas:**
```
// UserRepository Interface
//          │
//          ├─→ PostgresUserRepository (raw SQL)
//          ├─→ MySQLUserRepository (raw SQL)
//          ├─→ GORMUserRepository (pakai GORM)
//          └─→ MongoUserRepository (NoSQL)

// Semua implementasi interface yang sama!