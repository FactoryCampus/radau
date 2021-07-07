package migrations

import (
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg/orm"
)

type user struct {
	ID              int
	Username        string            `sql:"username,notnull,unique"`
	Token           string            `sql:"token"`
	ExtraProperties map[string]string `sql:"extraProperties" pg:"hstore"`
}

func init() {
	migrations.RegisterTx(func(db migrations.DB) error {
		err := orm.CreateTable(db, &user{}, &orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		return err
	}, func(db migrations.DB) error {
		err := orm.DropTable(db, &user{}, &orm.DropTableOptions{
			IfExists: true,
		})
		return err
	})
}
