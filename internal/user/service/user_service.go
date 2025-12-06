package service

import (
    "errors"

    "github.com/yourname/go-backend-enterprise/config"
    "github.com/yourname/go-backend-enterprise/internal/pkg/logger"
    "github.com/yourname/go-backend-enterprise/internal/user/repository"
    "golang.org/x/crypto/bcrypt"
)

type UserService struct {
    Repo   repository.UserRepository
    logger *logger.Logger
    cfg    *config.Config
}

func NewUserService(r repository.UserRepository, lg *logger.Logger, cfg *config.Config) *UserService {
    return &UserService{Repo: r, logger: lg, cfg: cfg}
}

func (s *UserService) Authenticate(username, password string) (*repository.User, error) {
    u, err := s.Repo.FindByUsername(username)
    if err != nil {
        return nil, err
    }
    if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
        return nil, errors.New("invalid credentials")
    }
    return u, nil
}

func (s *UserService) CreateUser(username, plain string) (int64, error) {
    if _, err := s.Repo.FindByUsername(username); err == nil {
        return 0, errors.New("user exists")
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
    if err != nil {
        return 0, err
    }
    u := &repository.User{Username: username, Password: string(hash)}
    return s.Repo.Create(u)
}

func (s *UserService) FindByID(id int64) (*repository.User, error) {
    return s.Repo.FindByID(id)
}
