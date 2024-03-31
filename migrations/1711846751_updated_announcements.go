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

		collection, err := dao.FindCollectionByNameOrId("21r5c1mha7urwzo")
		if err != nil {
			return err
		}

		collection.CreateRule = types.Pointer("// Only allow logged in users.\n@request.auth.id != \"\" && \n\n// Ensure approved is false.\n@request.data.approved = false && \n\n@request.data.author.id = @request.auth.id")

		collection.UpdateRule = types.Pointer("// Allow the author or admins to update. \n(@request.auth.id = author || @request.auth.role = 'admin') &&\n\n// Only allow admins to make changes.\n(@request.data.approved:isset = false || @request.auth.role = 'admin')")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("21r5c1mha7urwzo")
		if err != nil {
			return err
		}

		collection.CreateRule = types.Pointer("// Only allow logged in users.\n@request.auth.id != \"\" && \n\n// Ensure approved is unset.\n@request.data.approved = false")

		collection.UpdateRule = types.Pointer("")

		return dao.SaveCollection(collection)
	})
}
