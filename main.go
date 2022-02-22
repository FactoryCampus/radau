package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/FactoryCampus/radau/api"
	"github.com/FactoryCampus/radau/migrations"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

func initCORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()

	originConfig, hasCorsOrigins := os.LookupEnv("CORS_ORIGINS")
	if !hasCorsOrigins {
		corsConfig.AllowAllOrigins = true
	} else {
		corsConfig.AllowOrigins = strings.Split(originConfig, ",")
	}

	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AllowCredentials = true

	return corsConfig
}

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
	r.Use(cors.New(initCORSConfig()))

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
