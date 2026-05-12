package main

import (
	"github.com/Gergenus/commerce/user-service/internal/brocker"
	"github.com/Gergenus/commerce/user-service/internal/config"
	"github.com/Gergenus/commerce/user-service/internal/handlers"
	"github.com/Gergenus/commerce/user-service/internal/repository"
	"github.com/Gergenus/commerce/user-service/internal/service"
	dbpkg "github.com/Gergenus/commerce/user-service/pkg/db"
	"github.com/Gergenus/commerce/user-service/pkg/jwtpkg"
	"github.com/Gergenus/commerce/user-service/pkg/kafka"
	"github.com/Gergenus/commerce/user-service/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.InitConfig()
	db := dbpkg.InitDB(cfg.PostgresURL)
	defer db.DB.Close()
	log := logger.SetUp(cfg.LogLevel)

	jwtstr := jwtpkg.NewUserJWTpkg(cfg.JWTSecret, cfg.AccessTTL, cfg.EmailTTL, cfg.JWTMailSecret)
	middleWare := handlers.NewUserMiddleware(jwtstr)
	kafkaConn := kafka.ConnectProducer([]string{cfg.KafkaURL})
	defer kafkaConn.Close()

	kafka := brocker.NewKafkaBrocker(kafkaConn)
	repo := repository.NewPostgresRepository(db)
	srv := service.NewUserService(log, &repo, jwtstr, cfg.RefreshTTl, kafka)
	handler := handlers.NewUserHandler(&srv)

	e := echo.New()
	e.Use(middleware.Recover())
	group := e.Group("/api/v1/users/auth")
	{
		group.POST("/register", handler.Register)
		group.POST("/login", handler.Login)
		group.POST("/refresh", handler.Refresh)
		group.POST("/logout", handler.Logout)            //fdgdfgd
		group.GET("/verification", handler.Verification) // автоматика по ссылке из почты
		group.POST("/confirmation", handler.RegistrationConfirmation, middleWare.ConfirmationMiddleware)
	}
	e.Start(":" + cfg.HTTPPort)
}
