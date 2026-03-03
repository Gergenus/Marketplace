package main

import (
	"context"
	"fmt"

	"github.com/Gergenus/commerce/product-service/internal/config"
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
	// repo.AddCategory(context.Background(), "SVO")
	// repo.AddCategory(context.Background(), "GAY")
	// z, err := repo.CreateProduct(context.Background(), models.Product{ProductName: "SEX", Price: 300.5, SellerID: 1488, CategoryID: 1})
	// if err != nil {
	// 	log.Println(err)
	// }
	// _, err = repo.CreateProduct(context.Background(), models.Product{ProductName: "FISTING", Price: 300.5, SellerID: 1488, CategoryID: 1})
	// if err != nil {
	// 	log.Println(err)
	// }
	// _, err = repo.CreateProduct(context.Background(), models.Product{ProductName: "ASS", Price: 300.5, SellerID: 1488, CategoryID: 2})
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(z)
	// model, err := repo.GetProductsByCategory(context.Background(), "GAY")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(model)
	z, err := repo.GetProductByID(context.Background(), 9)
	fmt.Println(z)
	svo, err := repo.AddStockByID(context.Background(), 1488, 9, 150)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(svo)
}
