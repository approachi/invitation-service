package repository

import (
	"context"
	"errors"

	"invitation-service/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const pgErrUniqueViolation = "23505"

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Use(ctx context.Context, code, email string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	invitationID, err := r.claimInvitation(ctx, tx, code)
	if err != nil {
		return err
	}

	if err := r.insertRegistration(ctx, tx, invitationID, email); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepository) claimInvitation(ctx context.Context, tx pgx.Tx, code string) (int64, error) {
	var id int64
	err := tx.QueryRow(ctx, `
		UPDATE invitations
		SET used_count = used_count + 1
		WHERE code = $1 AND used_count < max_uses
		RETURNING id
	`, code).Scan(&id)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, r.classifyMissedUpdate(ctx, tx, code)
	}
	return id, err
}

func (r *PostgresRepository) insertRegistration(ctx context.Context, tx pgx.Tx, invitationID int64, email string) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO registrations (invitation_id, email) VALUES ($1, $2)
	`, invitationID, email)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgErrUniqueViolation {
			return domain.ErrAlreadyRegistered
		}
	}
	return err
}

func (r *PostgresRepository) classifyMissedUpdate(ctx context.Context, tx pgx.Tx, code string) error {
	var exists bool
	err := tx.QueryRow(ctx, `
		SELECT EXISTS(SELECT 1 FROM invitations WHERE code = $1)
	`, code).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return domain.ErrInvitationNotFound
	}
	return domain.ErrInvitationExhausted
}
