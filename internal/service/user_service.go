package service

import (
    "context"
    "fmt"
    "sync"

    "go-api-starter/internal/domain"
    "go-api-starter/internal/repository"
    "golang.org/x/crypto/bcrypt"
)

type UserService struct {
    repo       repository.UserRepository
    workerPool *WorkerPool
}

func NewUserService(repo repository.UserRepository, workerPool *WorkerPool) *UserService {
    return &UserService{
        repo:       repo,
        workerPool: workerPool,
    }
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.UserResponse, error) {
    // Check if user already exists
    existingUser, _ := s.repo.FindByEmail(ctx, req.Email)
    if existingUser != nil {
        return nil, domain.ErrUserAlreadyExists
    }

    existingUser, _ = s.repo.FindByUsername(ctx, req.Username)
    if existingUser != nil {
        return nil, domain.ErrUserAlreadyExists
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, fmt.Errorf("failed to hash password: %w", err)
    }

    // Create user
    user := &domain.User{
        Username:     req.Username,
        Email:        req.Email,
        PasswordHash: string(hashedPassword),
        FullName:     req.FullName,
        IsActive:     true,
    }

    createdUser, err := s.repo.Create(ctx, user)
    if err != nil {
        return nil, err
    }

    return createdUser.ToResponse(), nil
}

// GetUserByID retrieves user by ID
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*domain.UserResponse, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    return user.ToResponse(), nil
}

// GetUserByEmail retrieves user by email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
    return s.repo.FindByEmail(ctx, email)
}

// ListUsers retrieves users with pagination
func (s *UserService) ListUsers(ctx context.Context, page, limit int) ([]*domain.UserResponse, int64, error) {
    offset := (page - 1) * limit

    users, err := s.repo.FindAll(ctx, limit, offset)
    if err != nil {
        return nil, 0, err
    }

    total, err := s.repo.Count(ctx)
    if err != nil {
        return nil, 0, err
    }

    responses := make([]*domain.UserResponse, len(users))
    for i, user := range users {
        responses[i] = user.ToResponse()
    }

    return responses, total, nil
}

// UpdateUser updates user data
func (s *UserService) UpdateUser(ctx context.Context, id int64, req *domain.UpdateUserRequest) (*domain.UserResponse, error) {
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Update fields if provided
    if req.Username != "" {
        user.Username = req.Username
    }
    if req.Email != "" {
        user.Email = req.Email
    }
    if req.FullName != "" {
        user.FullName = req.FullName
    }
    if req.IsActive != nil {
        user.IsActive = *req.IsActive
    }

    updatedUser, err := s.repo.Update(ctx, user)
    if err != nil {
        return nil, err
    }

    return updatedUser.ToResponse(), nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
    return s.repo.Delete(ctx, id)
}

// BulkCreateUsers creates multiple users concurrently
func (s *UserService) BulkCreateUsers(ctx context.Context, requests []domain.CreateUserRequest) map[string]interface{} {
    var (
        successCount int
        failCount    int
        mu           sync.Mutex
        wg           sync.WaitGroup
        errors       []string
    )

    // Semaphore untuk membatasi concurrent goroutines
    semaphore := make(chan struct{}, 10)

    for i, req := range requests {
        wg.Add(1)

        go func(index int, request domain.CreateUserRequest) {
            defer wg.Done()

            // Acquire semaphore
            semaphore <- struct{}{}
            defer func() { <-semaphore }()

            // Create user
            _, err := s.CreateUser(ctx, &request)

            mu.Lock()
            defer mu.Unlock()

            if err != nil {
                failCount++
                errors = append(errors, fmt.Sprintf("User %d (%s): %v", index, request.Email, err))
            } else {
                successCount++
            }
        }(i, req)
    }

    wg.Wait()

    return map[string]interface{}{
        "total":   len(requests),
        "success": successCount,
        "failed":  failCount,
        "errors":  errors,
    }
}

// VerifyPassword verifies user password
func (s *UserService) VerifyPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
```

**Penjelasan:**
- Business logic layer
- Validasi bisnis (cek duplikat email/username)
- Password hashing dengan bcrypt
- Bulk create dengan **concurrency** (goroutines + semaphore)
- Pagination support

**Concurrency Flow di BulkCreateUsers:**
```
// 100 users to create
//        ↓
// Semaphore (max 10 concurrent)
//        ↓
// ┌──────┬──────┬──────┬──────┬──────┐
// │ Go 1 │ Go 2 │ Go 3 │ ... │ Go 10│
// └──────┴──────┴──────┴──────┴──────┘
//        ↓
// WaitGroup wait semua selesai
//        ↓
// Return result: {success: 95, failed: 5}