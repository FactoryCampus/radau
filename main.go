package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/factorycampus/radau/api"
	"github.com/factorycampus/radau/migrations"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

func main() {
	isDebug := os.Getenv("DEBUG")
	if strings.ToLower(isDebug) != "true" && isDebug != "1" {
		gin.SetMode(gin.ReleaseMode)
	}

	db := pg.Connect(&pg.Options{
		Addr:       fmt.Sprintf("%s:5432", os.Getenv("DB_HOST")),
		User:       os.Getenv("DB_USER"),
		Password:   os.Getenv("DB_PASSWORD"),
		Database:   os.Getenv("DB_DATABASE"),
		MaxRetries: 4,
	})
	defer db.Close()

	err := migrations.EnsureCurrentSchema(db)
	if err != nil {
		panic(err)
	}

	authManagementKey := os.Getenv("API_KEY_MANAGEMENT")
	authRadiusKey := os.Getenv("API_KEY_RADIUS")
	if !gin.IsDebugging() && (authManagementKey == "" || authRadiusKey == "") {
		panic("Please secure the service using environment variables: API_KEY_MANAGEMENT & API_KEY_RADIUS")
	}

	r := gin.Default()

	// Auth is handled individually per route
	api.InitUserHandler(r, db)

	// Auth is handled by TokenHandler
	api.InitTokenHandler(r.Group("/token"), db)

	radiusRoutes := r.Group("/radius", gin.BasicAuth(gin.Accounts{
		"Radius": authRadiusKey,
		"radius": authRadiusKey,
	}))
	api.InitRadiusHandler(radiusRoutes, db)

	r.Run()
}
