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

		collection.ListRule = types.Pointer("// Only Grant Individuals Access to Unnaproved Announcements If They Are The Owner Or an Admin. \napproved = true || @request.auth.role = 'admin' || @request.auth.id = author.id")

		collection.ViewRule = types.Pointer("// Only Grant Individuals Access to Unnaproved Announcements If They Are The Owner Or an Admin. \napproved = true || @request.auth.role = 'admin' || @request.auth.id = author.id")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("21r5c1mha7urwzo")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("")

		collection.ViewRule = types.Pointer("")

		return dao.SaveCollection(collection)
	})
}
