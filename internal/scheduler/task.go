package scheduler

import (
	"context"
	"go.uber.org/zap"
	"go-api-starter/internal/repository"
)

// CleanupExpiredTokensTask implements Task interface
type CleanupExpiredTokensTask struct {
	tokenRepo repository.JWTTokenRepository
	logger    *zap.SugaredLogger
}

func NewCleanupExpiredTokensTask(
	tokenRepo repository.JWTTokenRepository,
	logger *zap.SugaredLogger,
) Task {
	return &CleanupExpiredTokensTask{
		tokenRepo: tokenRepo,
		logger:    logger,
	}
}

func (t *CleanupExpiredTokensTask) Name() string {
	return "cleanup-expired-tokens"
}

func (t *CleanupExpiredTokensTask) Execute() error {
	t.logger.Info("Starting cleanup of expired JWT tokens")

	err := t.tokenRepo.DeleteExpiredTokens(context.Background())
	if err != nil {
		t.logger.Errorf("Failed to cleanup expired tokens: %v", err)
		return err
	}

	t.logger.Info("Expired JWT tokens cleaned up successfully")
	return nil
}
