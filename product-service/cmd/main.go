package main

import (
	"github.com/Gergenus/commerce/user-service/internal/config"
	handlers "github.com/Gergenus/commerce/user-service/internal/handler"
	"github.com/Gergenus/commerce/user-service/internal/repository"
	"github.com/Gergenus/commerce/user-service/internal/service"
	dbpkg "github.com/Gergenus/commerce/user-service/pkg/db"
	"github.com/Gergenus/commerce/user-service/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.InitConfig()
	db := dbpkg.InitDB(cfg.PostgresURL)
	log := logger.SetupLogger(cfg.LogLevel)

	repo := repository.NewPostgresRepository(db, log)
	serv := service.NewProductService(log, &repo)
	hand := handlers.NewProductHandler(&serv)
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	group := e.Group("/api/v1/products")
	{
		group.POST("/", hand.AddCategory)
	}

	e.Start(":" + cfg.HTTPPort)
}
