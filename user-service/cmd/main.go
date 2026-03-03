package main

import (
	"github.com/Gergenus/commerce/user-service/internal/config"
	dbpkg "github.com/Gergenus/commerce/user-service/pkg/db"
	"github.com/Gergenus/commerce/user-service/pkg/logger"
)

func main() {
	cfg := config.InitConfig()
	db := dbpkg.InitDB(cfg.PostgresURL)
	log := logger.SetupLogger(cfg.LogLevel)

}
