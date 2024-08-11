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
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	IsAdmin      bool       `json:"is_admin"`
	Status       UserStatus `json:"user_status"`
	CreatedTS    time.Time  `json:"created_ts"`
	UpdatedTS    time.Time  `json:"updated_ts"`
}

func fromSqlcUser(usr dbsqlc.User) User {
	return User{
		ID:           usr.ID,
		Name:         usr.Name,
		Email:        usr.Email,
		PasswordHash: usr.PasswordHash,
		IsAdmin:      usr.IsAdmin,
		Status:       UserStatus(usr.Status),
		CreatedTS:    usr.CreatedTs.UTC(),
		UpdatedTS:    usr.UpdatedTs.UTC(),
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

	return fromSqlcUser(user), nil
}

func GetUserByID(ctx context.Context, e db.Executor, id uuid.UUID) (User, error) {
	user, err := dbsqlc.New().GetUserByID(ctx, e, id)
	if err != nil {
		return User{}, db.HandleError(err)
	}

	return fromSqlcUser(user), nil
}
