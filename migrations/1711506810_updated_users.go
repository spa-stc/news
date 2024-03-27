package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		collection.CreateRule = types.Pointer("@request.data.role = 'student'")

		collection.UpdateRule = types.Pointer("id = @request.auth.id &&\n\n// Only allow students as rbac roles.\n@request.data.role = 'student' &&\n\n// Disallow email changes.\n@request.data.email:isset = false")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		collection.CreateRule = types.Pointer("@request.auth.role = 'student'")

		collection.UpdateRule = types.Pointer("id = @request.auth.id &&\n@request.auth.role = 'student'")

		return dao.SaveCollection(collection)
	})
}
