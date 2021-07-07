package migrations

import (
	"fmt"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
)

func EnsureCurrentSchema(db *pg.DB) error {
	var oldVersion, newVersion int64

	oldVersion, newVersion, err := migrations.Run(db, "up")
	if err != nil {
		return err
	}

	if newVersion != oldVersion {
		fmt.Printf("db schema migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("db schema version is %d\n", oldVersion)
	}

	return nil
}
