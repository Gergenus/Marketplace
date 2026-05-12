package main

import (
	"context"
	"net"

	"github.com/Gergenus/commerce/product-service/internal/config"
	"github.com/Gergenus/commerce/product-service/internal/handlers"
	"github.com/Gergenus/commerce/product-service/internal/repository"
	"github.com/Gergenus/commerce/product-service/internal/service"
	dbpkg "github.com/Gergenus/commerce/product-service/pkg/db"
	"github.com/Gergenus/commerce/product-service/pkg/elastic"
	"github.com/Gergenus/commerce/product-service/pkg/jwtpkg"
	"github.com/Gergenus/commerce/product-service/pkg/logger"
	"github.com/Gergenus/commerce/product-service/proto"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.InitConfig()
	db := dbpkg.InitDB(cfg.PostgresURL)
	defer db.DB.Close(context.Background())
	log := logger.SetupLogger(cfg.LogLevel)

	repo := repository.NewPostgresRepository(db)
	eClient := elastic.NewElasticClient([]string{cfg.ElasticAddress}, cfg.ElasticUser, cfg.ElasticPassword, cfg.ElasticCrt, log, &repo)
	serv := service.NewProductService(log, &repo, &eClient)
	hand := handlers.NewProductHandler(&serv)
	jwtPkg := jwtpkg.NewJWTpkg(cfg.JWTSecret, log)
	middleWare := handlers.NewProductMiddleware(jwtPkg)

	// use faktory

	isCreated := eClient.InitIndexation(context.Background())
	if isCreated {
		err := eClient.IndexAllProducts(context.Background())
		if err != nil {
			panic(err)
		}

	}

	lis, err := net.Listen("tcp", cfg.GRPCProductServerAddress)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	proto.RegisterAvailablilityServiceServer(s, &hand)
	proto.RegisterOrderServiceServer(s, &hand)
	go func() {
		log.Info("starting gRPC server")
		if err := s.Serve(lis); err != nil {
			panic(err)
		}

	}()
	e := echo.New()

	e.Use(middleware.Recover())
	group := e.Group("/api/v1/products", middleWare.RateLimiting)
	{
		group.POST("/", hand.AddCategory, middleWare.Auth)
		group.POST("/create", hand.CreateProduct, middleWare.CreateProductAuth) // create product
		group.GET("/", hand.GetProductByID)                                     // get product by id
		group.POST("/stock/add", hand.AddStockByID, middleWare.Auth)            // add stock by id
		group.GET("/stock", hand.GetStockByID, middleWare.Auth)                 // get stock by id
		group.GET("", hand.Products)
	}

	e.Start(":" + cfg.HTTPPort)

}
