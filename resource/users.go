package resource

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"stpaulacademy.tech/newsletter/db"
	"stpaulacademy.tech/newsletter/db/dbsqlc"
	"stpaulacademy.tech/newsletter/util/passwordutil"
)

type UserStatus string

const (
	UserStatusVerified   UserStatus = "verified"
	UserStatusUnverified UserStatus = "unverified"
	UserStatusBanned     UserStatus = "banned"
)

type NewUser struct {
	Name     string
	Email    string
	Password string
}

type User struct {
	ID                   uuid.UUID
	Name                 string
	Email                string
	PasswordHash         string
	IsAdmin              bool
	Status               UserStatus
	VerificationAttempts int64
	CreatedTS            time.Time
	UpdatedTS            time.Time
}

type UpdateUser struct {
	Password *string
	Status   *UserStatus
}

type TokenClaims struct {
	ID      uuid.UUID
	IsAdmin bool
	Status  UserStatus
}

func fromSqlcUser(usr dbsqlc.User) User {
	return User{
		ID:                   usr.ID,
		Name:                 usr.Name,
		Email:                usr.Email,
		PasswordHash:         usr.PasswordHash,
		IsAdmin:              usr.IsAdmin,
		Status:               UserStatus(usr.Status),
		CreatedTS:            usr.CreatedTs.UTC(),
		UpdatedTS:            usr.UpdatedTs.UTC(),
		VerificationAttempts: 0,
	}
}

func CreateUser(ctx context.Context, e db.Executor, nu NewUser) (User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return User{}, fmt.Errorf("error generating id: %w", err)
	}

	hash, err := passwordutil.HashPassword(nu.Password)
	if err != nil {
		return User{}, fmt.Errorf("error hashing password: %w", err)
	}

	user, err := dbsqlc.New().InsertUser(ctx, e, dbsqlc.InsertUserParams{
		ID:           id,
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
	})
	if err != nil {
		return User{}, db.HandleError(err)
	}

	return fromSqlcUser(user), nil
}

func GetUserByEmail(ctx context.Context, e db.Executor, email string) (User, error) {
	user, err := dbsqlc.New().GetUserByEmail(ctx, e, email)
	if err != nil {
		return User{}, db.HandleError(err)
	}

	usr := User{
		ID:                   user.ID,
		Name:                 user.Name,
		Email:                user.Email,
		PasswordHash:         user.PasswordHash,
		IsAdmin:              user.IsAdmin,
		Status:               UserStatus(user.Status),
		VerificationAttempts: user.VerificationAttempts,
		CreatedTS:            user.CreatedTs.UTC(),
		UpdatedTS:            user.UpdatedTs.UTC(),
	}

	return usr, nil
}

func GetUserByID(ctx context.Context, e db.Executor, id uuid.UUID) (User, error) {
	user, err := dbsqlc.New().GetUserByID(ctx, e, id)
	if err != nil {
		return User{}, db.HandleError(err)
	}

	usr := User{
		ID:                   user.ID,
		Name:                 user.Name,
		Email:                user.Email,
		PasswordHash:         user.PasswordHash,
		IsAdmin:              user.IsAdmin,
		Status:               UserStatus(user.Status),
		VerificationAttempts: user.VerificationAttempts,
		CreatedTS:            user.CreatedTs.UTC(),
		UpdatedTS:            user.UpdatedTs.UTC(),
	}

	return usr, nil
}

func UpdateUserByID(ctx context.Context, e db.Executor, id uuid.UUID, uu UpdateUser) error {
	sqlcusr := dbsqlc.UpdateUserByIDParams{
		ID: id,
	}

	if uu.Status != nil {
		sqlcusr.StatusDoUpdate = true
		sqlcusr.Status = dbsqlc.UserStatus(*uu.Status)
	} else {
		sqlcusr.Status = dbsqlc.UserStatusBanned
	}

	if uu.Password != nil {
		hash, err := passwordutil.HashPassword(*uu.Password)
		if err != nil {
			return fmt.Errorf("error hashing password: %w", err)
		}

		sqlcusr.PasswordHashDoUpdate = true
		sqlcusr.PasswordHash = hash
	}

	err := dbsqlc.New().UpdateUserByID(ctx, e, sqlcusr)
	if err != nil {
		return db.HandleError(err)
	}

	return nil
}

func GetTokenClaims(ctx context.Context, e db.Executor, id uuid.UUID) (TokenClaims, error) {
	claims, err := dbsqlc.New().GetTokenClaimsByUserID(ctx, e, id)
	if err != nil {
		return TokenClaims{}, db.HandleError(err)
	}

	c := TokenClaims{
		ID:      claims.ID,
		IsAdmin: claims.IsAdmin,
		Status:  UserStatus(claims.Status),
	}

	return c, nil
}
