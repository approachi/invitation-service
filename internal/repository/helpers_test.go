package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func truncate(t *testing.T) {
	t.Helper()

	exec := func() {
		_, err := pool.Exec(context.Background(), "TRUNCATE registrations, invitations RESTART IDENTITY CASCADE")
		require.NoError(t, err)
	}

	exec()
	t.Cleanup(exec)
}

func seedInvitation(t *testing.T, code string, maxUses, usedCount int) {
	t.Helper()
	_, err := pool.Exec(context.Background(),
		"INSERT INTO invitations (code, max_uses, used_count) VALUES ($1, $2, $3)",
		code, maxUses, usedCount,
	)
	require.NoError(t, err)
}
