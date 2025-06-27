package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-recipes-server/internal/db"
	"go-recipes-server/internal/handler"
	"go-recipes-server/internal/middleware"
	"gorm.io/gorm"
	"os"
)

func main() {
	env := os.Getenv("APP_ENV")

	r, dbConn := initApp()

	r.Use(cors.Default())

	initRoutes(r, dbConn)

	if env == "" {
		gin.SetMode(gin.DebugMode)
		r.Use(gin.Logger())
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r.Run(":8080")
}

func initApp() (*gin.Engine, *gorm.DB) {
	dbConn := db.Connect()

	r := gin.Default()

	return r, dbConn
}

func initRoutes(r *gin.Engine, dbConn *gorm.DB) {
	// static files
	r.Static("/public", "./public")
	// auth routes
	authRouter := r.Group("/auth")
	handler.NewAuthHandler(dbConn, authRouter)
	// resource routes
	apiRouter := r.Group("/api")
	apiRouter.Use(middleware.JWTMiddleware())
	handler.NewRecipeHandler(dbConn, apiRouter)
}
