package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		options := map[string]any{}
		if err := json.Unmarshal([]byte(`{
			"allowEmailAuth": true,
			"allowOAuth2Auth": false,
			"allowUsernameAuth": false,
			"exceptEmailDomains": null,
			"manageRule": null,
			"minPasswordLength": 12,
			"onlyEmailDomains": [
				"students.spa.edu",
				"spa.edu"
			],
			"onlyVerified": true,
			"requireEmail": true
		}`), &options); err != nil {
			return err
		}
		collection.SetOptions(options)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		options := map[string]any{}
		if err := json.Unmarshal([]byte(`{
			"allowEmailAuth": true,
			"allowOAuth2Auth": false,
			"allowUsernameAuth": false,
			"exceptEmailDomains": null,
			"manageRule": null,
			"minPasswordLength": 8,
			"onlyEmailDomains": [
				"students.spa.edu",
				"spa.edu"
			],
			"onlyVerified": true,
			"requireEmail": true
		}`), &options); err != nil {
			return err
		}
		collection.SetOptions(options)

		return dao.SaveCollection(collection)
	})
}
