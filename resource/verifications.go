package resource

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/xlzd/gotp"
	"stpaulacademy.tech/newsletter/db"
	"stpaulacademy.tech/newsletter/db/dbsqlc"
)

const secretLength = 32

type EmailVerification struct {
	UserID uuid.UUID
	Secret string
	Used   bool
}

func CreateEmailVerification(ctx context.Context, e db.Executor, userID uuid.UUID) (EmailVerification, error) {
	secret := gotp.RandomSecret(secretLength)
	if secret == "" {
		return EmailVerification{}, errors.New("failed to generate secret")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return EmailVerification{}, fmt.Errorf("error generating uuid: %w", err)
	}

	dbver, err := dbsqlc.New().InsertEmailVerification(ctx, e, dbsqlc.InsertEmailVerificationParams{
		ID:     id,
		UserID: userID,
		Secret: secret,
	})
	if err != nil {
		return EmailVerification{}, db.HandleError(err)
	}

	v := EmailVerification{
		UserID: dbver.UserID,
		Secret: dbver.Secret,
		Used:   dbver.Used,
	}

	return v, nil
}

func GetEmailVerificationBySecret(ctx context.Context, e db.Executor, secret string) (EmailVerification, error) {
	dbver, err := dbsqlc.New().GetEmailVerificationBySecret(ctx, e, secret)
	if err != nil {
		return EmailVerification{}, db.HandleError(err)
	}

	v := EmailVerification{
		UserID: dbver.UserID,
		Secret: dbver.Secret,
		Used:   dbver.Used,
	}

	return v, nil
}
