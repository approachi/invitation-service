package repository_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"invitation-service/internal/domain"
	"invitation-service/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUse_ShouldSucceed_WhenInvitationIsAvailable(t *testing.T) {
	truncate(t)
	seedInvitation(t, "test-code", 10, 0)

	repo := repository.NewPostgresRepository(pool)
	err := repo.Use(context.Background(), "test-code", "user@example.com")

	require.NoError(t, err)
}

func TestUse_ShouldReturnNotFound_WhenCodeIsUnknown(t *testing.T) {
	truncate(t)

	repo := repository.NewPostgresRepository(pool)
	err := repo.Use(context.Background(), "unknown-code", "user@example.com")

	assert.ErrorIs(t, err, domain.ErrInvitationNotFound)
}

func TestUse_ShouldReturnExhausted_WhenMaxUsesReached(t *testing.T) {
	truncate(t)
	seedInvitation(t, "test-code", 5, 5)

	repo := repository.NewPostgresRepository(pool)
	err := repo.Use(context.Background(), "test-code", "user@example.com")

	assert.ErrorIs(t, err, domain.ErrInvitationExhausted)
}

func TestUse_ShouldReturnAlreadyRegistered_WhenEmailAlreadyUsed(t *testing.T) {
	truncate(t)
	seedInvitation(t, "test-code", 10, 0)

	repo := repository.NewPostgresRepository(pool)
	require.NoError(t, repo.Use(context.Background(), "test-code", "user@example.com"))

	err := repo.Use(context.Background(), "test-code", "user@example.com")
	assert.ErrorIs(t, err, domain.ErrAlreadyRegistered)
}

func TestUse_ShouldNotExceedMaxUses_WhenCalledConcurrently(t *testing.T) {
	truncate(t)
	maxUses := 5
	seedInvitation(t, "test-code", maxUses, 0)

	repo := repository.NewPostgresRepository(pool)

	var wg sync.WaitGroup
	goroutines := 20
	wg.Add(goroutines)

	for i := range goroutines {
		go func(i int) {
			defer wg.Done()
			_ = repo.Use(context.Background(), "test-code", fmt.Sprintf("user%d@example.com", i))
		}(i)
	}
	wg.Wait()

	var usedCount int
	err := pool.QueryRow(context.Background(),
		"SELECT used_count FROM invitations WHERE code = $1", "test-code",
	).Scan(&usedCount)
	require.NoError(t, err)
	assert.Equal(t, maxUses, usedCount)
}
