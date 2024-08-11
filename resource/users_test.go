package resource_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/db"
	"stpaulacademy.tech/newsletter/resource"
	"stpaulacademy.tech/newsletter/util/passwordutil"
	"stpaulacademy.tech/newsletter/util/testutil"
)

func TestUsersResource(t *testing.T) {
	t.Parallel()
	ctx := testutil.Setup(t)

	seededUser := resource.User{
		ID:                   uuid.MustParse("aea38951-ca26-4e76-ad65-d5296a0095e6"),
		Name:                 "name",
		Email:                "email",
		PasswordHash:         "password_hash",
		IsAdmin:              false,
		Status:               resource.UserStatusUnverified,
		CreatedTS:            time.UnixMicro(0),
		UpdatedTS:            time.UnixMicro(0),
		VerificationAttempts: 3,
	}

	t.Run("test_get_id_not_found", func(t *testing.T) {
		tx := testutil.TestTx(ctx, t)
		_, err := resource.GetUserByID(ctx, tx, uuid.MustParse("b082b4b2-845f-4512-8034-506cad6ff097"))

		require.ErrorIs(t, err, db.ErrNotFound)
	})

	t.Run("test_get_id_found", func(t *testing.T) {
		tx := testutil.TestTx(ctx, t)
		usr, err := resource.GetUserByID(ctx, tx, uuid.MustParse("aea38951-ca26-4e76-ad65-d5296a0095e6"))
		require.NoError(t, err)

		usr.CreatedTS = time.UnixMicro(0)
		usr.UpdatedTS = time.UnixMicro(0)

		require.Equal(t, seededUser, usr)
	})

	t.Run("test_get_email_not_found", func(t *testing.T) {
		tx := testutil.TestTx(ctx, t)
		_, err := resource.GetUserByEmail(ctx, tx, "hi@me.com")

		require.ErrorIs(t, err, db.ErrNotFound)
	})

	t.Run("test_get_email_found", func(t *testing.T) {
		tx := testutil.TestTx(ctx, t)
		usr, err := resource.GetUserByEmail(ctx, tx, "email")
		require.NoError(t, err)

		usr.CreatedTS = time.UnixMicro(0)
		usr.UpdatedTS = time.UnixMicro(0)

		require.Equal(t, seededUser, usr)
	})

	t.Run("test_insert_unique_violation", func(t *testing.T) {
		tx := testutil.TestTx(ctx, t)
		_, err := resource.CreateUser(ctx, tx, resource.NewUser{
			Name:     "Test Testy",
			Email:    "email",
			Password: "1234",
		})

		require.Error(t, err)
	})

	t.Run("test_insert_success", func(t *testing.T) {
		tx := testutil.TestTx(ctx, t)
		user, err := resource.CreateUser(ctx, tx, resource.NewUser{
			Name:     "Test Testy",
			Email:    "test@test.com",
			Password: "1234",
		})
		require.NoError(t, err)

		found, err := resource.GetUserByID(ctx, tx, user.ID)
		require.NoError(t, err)

		require.Equal(t, user, found)
	})

	t.Run("test_users_update", func(t *testing.T) {
		tx := testutil.TestTx(ctx, t)

		err := resource.UpdateUserByID(ctx, tx, seededUser.ID, resource.UpdateUser{
			Password: testutil.Pointer("hi"),
			Status:   testutil.Pointer(resource.UserStatusVerified),
		})
		require.NoError(t, err)

		usr, err := resource.GetUserByID(ctx, tx, seededUser.ID)
		require.NoError(t, err)

		ok, err := passwordutil.VerifyPassword(usr.PasswordHash, "hi")
		require.NoError(t, err)
		require.True(t, ok)

		require.Equal(t, resource.UserStatusVerified, usr.Status)
	})
}
