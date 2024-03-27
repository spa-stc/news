package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("21r5c1mha7urwzo")
		if err != nil {
			return err
		}

		collection.CreateRule = types.Pointer("// Force Author to be creating user.\n@request.auth.id != \"\" &&\n@request.auth.id = @request.data.author &&\n\n// Make sure that we are always unapproved upon creation, unless the creator is an admin.\n(@request.data.approved = false || @request.auth.role = 'admin') &&\n\n// Make sure that we cannot set the start date to be less than the end date.\n(@request.data.start_showing_at < @request.data.finish_showing_at) ")

		collection.UpdateRule = types.Pointer("// Disallow author from being changed. \n@request.data.author:isset = false &&\n\n// Only allow admins to change approved field.\n(@request.data.approved:isset = false || @request.auth.role = 'admin') &&\n\n// If the announcement is already approved, then only allow admins to make updates.\n(approved = false || @request.auth.role = 'admin') &&\n\n// Make sure that we cannot set the start date to be less than the end date.\n(@request.data.start_showing_at < @request.data.finish_showing_at)")

		// add
		new_start_showing_at := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "vpw0pwsf",
			"name": "start_showing_at",
			"type": "date",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), new_start_showing_at); err != nil {
			return err
		}
		collection.Schema.AddField(new_start_showing_at)

		// add
		new_finish_showing_at := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "hno4op54",
			"name": "finish_showing_at",
			"type": "date",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"min": "",
				"max": ""
			}
		}`), new_finish_showing_at); err != nil {
			return err
		}
		collection.Schema.AddField(new_finish_showing_at)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("21r5c1mha7urwzo")
		if err != nil {
			return err
		}

		collection.CreateRule = types.Pointer("// Force Author to be creating user.\n@request.auth.id != \"\" &&\n@request.auth.id = @request.data.author &&\n\n// Make sure that we are always unapproved upon creation, unless the creator is an admin.\n(@request.data.approved = false || @request.auth.role = 'admin')")

		collection.UpdateRule = types.Pointer("// Disallow author from being changed. \n@request.data.author:isset = false &&\n\n// Only allow admins to change approved field.\n(@request.data.approved:isset = false || @request.auth.role = 'admin') &&\n\n// If the announcement is already approved, then only allow admins to make updates.\napproved = false || @request.auth.role = 'admin'")

		// remove
		collection.Schema.RemoveField("vpw0pwsf")

		// remove
		collection.Schema.RemoveField("hno4op54")

		return dao.SaveCollection(collection)
	})
}
