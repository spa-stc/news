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

		collection.CreateRule = types.Pointer("// Force Author to be creating user.\n@request.auth.id != \"\" &&\n@request.auth.id = @request.data.author.id &&\n\n// Make sure that we are always unapproved upon creation, unless the creator is an admin.\n(@request.data.approved = false || @request.auth.role = 'admin') &&\n\n// Make sure that we cannot set the start date to be less than the end date.\n(@request.data.start_showing_at < @request.data.finish_showing_at) ")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("21r5c1mha7urwzo")
		if err != nil {
			return err
		}

		collection.CreateRule = types.Pointer("// Force Author to be creating user.\n@request.auth.id != \"\" &&\n@request.auth.id = @request.data.author &&\n\n// Make sure that we are always unapproved upon creation, unless the creator is an admin.\n(@request.data.approved = false || @request.auth.role = 'admin') &&\n\n// Make sure that we cannot set the start date to be less than the end date.\n(@request.data.start_showing_at < @request.data.finish_showing_at) ")

		return dao.SaveCollection(collection)
	})
}
