package main

import (
	"auth/config"
	"auth/internal/http"
	"auth/internal/repository"
	"auth/internal/repository/mongodb"
	"auth/internal/repository/postgres"
	"auth/internal/service"
	"auth/internal/usecase"
	"auth/pkg/logger"
	"os"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to load config")
		os.Exit(1)
	}

	logcfg := logger.Config{
		Level:  "info",
		Pretty: true,
	}
	logger.InitGlobalLogger(&logcfg)

	logger.Logger.Info().Msg("Config loaded")

	// PostgreSQL
	pool, err := repository.NewPgxPool(cfg.Database)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to connect to PostgreSQL")
	}

	userRepoPostgres := postgres.NewUserPostgres(pool)

	// MongoDB
	mongoClient, err := mongodb.NewMongoClient(cfg.Database)
	if err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to connect to MongoDB")
	}

	reportRepoMongo := mongodb.NewReportMongo(mongoClient.Database(cfg.Database.Mongo.Name))

	// Services & Usecases
	jwtService := service.NewJWTService(cfg.Auth.JWTSecret)
	userUC := usecase.NewUserUsecase(userRepoPostgres, reportRepoMongo, jwtService)
	reportUC := usecase.NewReportUsecase(reportRepoMongo, userRepoPostgres)

	// Start HTTP server
	if err := http.StartServer(cfg, userUC, reportUC); err != nil {
		logger.Logger.Fatal().Err(err).Msg("failed to start server")
	}
}
