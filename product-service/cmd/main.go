package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Gergenus/commerce/product-service/internal/config"
	"github.com/Gergenus/commerce/product-service/internal/models"
	"github.com/Gergenus/commerce/product-service/internal/repository"
	dbpkg "github.com/Gergenus/commerce/product-service/pkg/db"
)

func main() {
	cfg := config.InitConfig()
	db := dbpkg.InitDB(cfg.PostgresURL)
	defer db.DB.Close(context.Background())
	// log := logger.SetupLogger(cfg.LogLevel)

	repo := repository.NewPostgresRepository(db)
	// serv := service.NewProductService(log, &repo)
	// hand := handlers.NewProductHandler(&serv)
	// e := echo.New()

	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	// group := e.Group("/api/v1/products")
	// {
	// 	group.POST("/", hand.AddCategory)
	// }

	// e.Start(":" + cfg.HTTPPort)
	z, err := repo.AddCategory(context.Background(), "SVO")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(z)
	id, err := repo.CreateProduct(context.Background(), models.Product{ProductName: "ТАНК", Price: 5000.5, SellerID: 5, CategoryID: 1})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}
