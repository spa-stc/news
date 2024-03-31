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

		collection.UpdateRule = types.Pointer("// Allow the author or admins to update. \n(@request.auth.id = author || @request.auth.role = 'admin') &&\n\n// Only allow admins to make changes.\n(@request.data.approved:isset = false || @request.auth.role = 'admin') &&\n\n(@request.data.author.id:isset = false)")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("21r5c1mha7urwzo")
		if err != nil {
			return err
		}

		collection.UpdateRule = types.Pointer("// Allow the author or admins to update. \n(@request.auth.id = author || @request.auth.role = 'admin') &&\n\n// Only allow admins to make changes.\n(@request.data.approved:isset = false || @request.auth.role = 'admin')")

		return dao.SaveCollection(collection)
	})
}
