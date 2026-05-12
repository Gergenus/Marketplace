package main

import (
	"net"

	"github.com/Gergenus/commerce/cart-service/internal/config"
	"github.com/Gergenus/commerce/cart-service/internal/handlers"
	"github.com/Gergenus/commerce/cart-service/internal/repository"
	"github.com/Gergenus/commerce/cart-service/internal/service"
	"github.com/Gergenus/commerce/cart-service/pkg/jwtpkg"
	"github.com/Gergenus/commerce/cart-service/pkg/logger"
	"github.com/Gergenus/commerce/cart-service/pkg/redispkg"
	"github.com/Gergenus/commerce/cart-service/proto"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.InitConfig()
	redisClient := redispkg.InitRedisDB(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB)
	log := logger.SetupLogger(cfg.LogLevel)
	jwtPkg := jwtpkg.NewCartJWTpkg(cfg.JWTSecret)
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(cfg.GRPCProductServiceAddress, opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := proto.NewAvailablilityServiceClient(conn)

	cartMiddleware := handlers.NewCartMiddleware(jwtPkg)

	repo := repository.NewRedisRepository(redisClient, cfg.CartTTL)
	srv := service.NewCartService(log, repo, client)
	hnd := handlers.NewCartHandler(srv)

	lis, err := net.Listen("tcp", cfg.GRPCCartServerAddress)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	proto.RegisterOrderServiceServer(s, hnd)
	go func() {
		log.Info("starting gRPC server")
		if err := s.Serve(lis); err != nil {
			panic(err)
		}

	}()

	e := echo.New()
	e.Use(middleware.Recover())
	group := e.Group("/api/v1/cart", cartMiddleware.CartMiddleware)
	{
		group.POST("/add", hnd.AddToCart)
		group.POST("/delete", hnd.DeleteFromCart)
		group.GET("/", hnd.Cart)
	}

	e.Start(":" + cfg.HTTPPort)
}
