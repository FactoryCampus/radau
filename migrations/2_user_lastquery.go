package migrations

import (
	"github.com/go-pg/migrations"
)

func init() {
	migrations.RegisterTx(func(db migrations.DB) error {
		_, err := db.Exec(`ALTER TABLE users ADD lastquery timestamptz`)
		return err
	}, func(db migrations.DB) error {
		_, err := db.Exec(`ALTER TABLE users DROP lastquery`)
		return err
	})
}
