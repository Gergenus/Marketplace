package main

import (
	"github.com/Gergenus/commerce/order-service/internal/config"
	"github.com/Gergenus/commerce/order-service/internal/handlers"
	"github.com/Gergenus/commerce/order-service/internal/repository"
	"github.com/Gergenus/commerce/order-service/internal/service"
	"github.com/Gergenus/commerce/order-service/pkg/db"
	"github.com/Gergenus/commerce/order-service/pkg/logger"
	"github.com/Gergenus/commerce/order-service/proto"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.InitConfig()
	log := logger.SetUp(cfg.LogLevel)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	productConn, err := grpc.NewClient(cfg.GRPCProductServiceAddress, opts...)
	if err != nil {
		panic(err)
	}
	defer productConn.Close()
	cartConn, err := grpc.NewClient(cfg.GRPCCartServiceAddress, opts...)
	if err != nil {
		panic(err)
	}
	defer cartConn.Close()

	db := db.InitDB(cfg.POSTGRES_URL)
	cartClient := proto.NewOrderServiceClient(cartConn)
	productClient := proto.NewOrderServiceClient(productConn)
	repo := repository.NewOrderRepository(db)
	srv := service.NewOrderService(&repo, log, cartClient, productClient)
	hnd := handlers.NewOrderHandler(srv)

	e := echo.New()
	e.Use(middleware.Recover())

	e.POST("/api/v1/order/", hnd.CreateOrder, handlers.OrderAuth)
	e.GET("/api/v1/order/", hnd.Orders, handlers.SellerAuth)

	e.Start(":" + cfg.HTTPPort)
}
