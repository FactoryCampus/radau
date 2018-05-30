package main

import (
	"fmt"
	"os"

	"./internal/pkg"
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

	for _, model := range []interface{}{(*pkg.User)(nil)} {
		err := db.CreateTable(model, tableOpts)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:5432", os.Getenv("DB_HOST")),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
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

	userRoutes := r.Group("/user", pkg.HandleApiKeyAuth(authManagementKey))
	pkg.InitUserHandler(userRoutes, db)

	tokenRoutes := r.Group("/token", pkg.HandleApiKeyAuth(authManagementKey))
	pkg.InitTokenHandler(tokenRoutes, db)

	radiusRoutes := r.Group("/radius", gin.BasicAuth(gin.Accounts{
		"Radius": authRadiusKey,
	}))
	pkg.InitRadiusHandler(radiusRoutes, db)

	r.Run()
}
