package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "21r5c1mha7urwzo",
			"created": "2024-03-27 02:22:29.358Z",
			"updated": "2024-03-27 02:22:29.358Z",
			"name": "announcements",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "7yjma9ks",
					"name": "title",
					"type": "text",
					"required": true,
					"presentable": true,
					"unique": false,
					"options": {
						"min": null,
						"max": 255,
						"pattern": ""
					}
				},
				{
					"system": false,
					"id": "2d49oxsk",
					"name": "content",
					"type": "editor",
					"required": true,
					"presentable": false,
					"unique": false,
					"options": {
						"convertUrls": false
					}
				},
				{
					"system": false,
					"id": "bgddqe3i",
					"name": "author",
					"type": "relation",
					"required": true,
					"presentable": false,
					"unique": false,
					"options": {
						"collectionId": "_pb_users_auth_",
						"cascadeDelete": true,
						"minSelect": null,
						"maxSelect": 1,
						"displayFields": null
					}
				},
				{
					"system": false,
					"id": "cpjxwo3l",
					"name": "approved",
					"type": "bool",
					"required": false,
					"presentable": true,
					"unique": false,
					"options": {}
				}
			],
			"indexes": [],
			"listRule": "",
			"viewRule": "",
			"createRule": "// Force Author to be creating user.\n@request.auth.id != \"\" &&\n@request.auth.id = @request.data.author &&\n\n// Make sure that we are always unapproved upon creation, unless the creator is an admin.\n(@request.data.approved = false || @request.auth.role = 'admin')",
			"updateRule": "// Disallow author from being changed. \n@request.data.author:isset = false &&\n\n// Only allow admins to change approved field.\n(@request.data.approved:isset = false || @request.auth.role = 'admin') &&\n\n// If the announcement is already approved, then only allow admins to make updates.\napproved = false || @request.auth.role = 'admin'",
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("21r5c1mha7urwzo")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
