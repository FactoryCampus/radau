package main

import (
	"fmt"
	"os"

	"github.com/factorycampus/radau/api"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

func createSchema(db *pg.DB) error {
	tableOpts := &orm.CreateTableOptions{
		Temp:        false,
		IfNotExists: true,
	}
	if gin.IsDebugging() {
		tableOpts = &orm.CreateTableOptions{
			Temp: true,
		}
	}

	for _, model := range []interface{}{(*api.User)(nil)} {
		err := db.CreateTable(model, tableOpts)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	db := pg.Connect(&pg.Options{
		Addr:       fmt.Sprintf("%s:5432", os.Getenv("DB_HOST")),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		Database:   os.Getenv("DB_DATABASE"),
		MaxRetries: 4,
	})
	defer db.Close()

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	authManagementKey := os.Getenv("API_KEY_MANAGEMENT")
	authRadiusKey := os.Getenv("API_KEY_RADIUS")
	if !gin.IsDebugging() && (authManagementKey == "" || authRadiusKey == "") {
		panic("Please secure the service using environment variables: API_KEY_MANAGEMENT & API_KEY_RADIUS")
	}

	r := gin.Default()

	managementAuthed := r.Group("", api.HandleApiKeyAuth(authManagementKey))
	api.InitUserHandler(managementAuthed, db)
	api.InitTokenHandler(managementAuthed, db)

	radiusRoutes := r.Group("/radius", gin.BasicAuth(gin.Accounts{
		"Radius": authRadiusKey,
		"radius": authRadiusKey,
	}))
	api.InitRadiusHandler(radiusRoutes, db)

	r.Run()
}
