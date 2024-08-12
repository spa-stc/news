package resource_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"stpaulacademy.tech/newsletter/db"
	"stpaulacademy.tech/newsletter/resource"
	"stpaulacademy.tech/newsletter/util/testutil"
)

func TestVerification(t *testing.T) {
	t.Parallel()

	ctx := testutil.Setup(t)

	t.Run("test_create", func(t *testing.T) {
		tx := testutil.TestTx(ctx, t)

		v, err := resource.CreateEmailVerification(ctx, tx, uuid.MustParse("aea38951-ca26-4e76-ad65-d5296a0095e6"))
		require.NoError(t, err)

		found, err := resource.GetEmailVerificationBySecret(ctx, tx, v.Secret)
		require.NoError(t, err)

		require.Equal(t, v, found)
	})

	t.Run("test_get_failure", func(t *testing.T) {
		tx := testutil.TestTx(ctx, t)

		_, err := resource.GetEmailVerificationBySecret(ctx, tx, "eersdsdsrerss")
		require.ErrorIs(t, err, db.ErrNotFound)
	})

	t.Run("test_get", func(t *testing.T) {
		tx := testutil.TestTx(ctx, t)

		expected := resource.EmailVerification{
			UserID: uuid.MustParse("aea38951-ca26-4e76-ad65-d5296a0095e6"),
			Secret: "FSPWAPYLS5ZOZM34BARF45DYDY",
			Used:   false,
		}

		v, err := resource.GetEmailVerificationBySecret(ctx, tx, expected.Secret)
		require.NoError(t, err)

		require.Equal(t, expected, v)
	})
}
