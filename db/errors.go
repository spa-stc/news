package db

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrNotFound        = errors.New("record not found")
	ErrUniqueViolation = errors.New("unique violation")
)

const codeUniqueViolation = "23505"

func HandleError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}

	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		if pgError.Code == codeUniqueViolation {
			return ErrUniqueViolation
		}
	}

	return err
}
