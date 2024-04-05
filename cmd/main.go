package main

import (
	"fmt"
	"github.com/IkoAfianando/mispress/db"
	"github.com/IkoAfianando/mispress/handler"
	"github.com/elastic/go-elasticsearch"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	var dbPort int
	var err error

	port := os.Getenv("POSTGRES_PORT")
	fmt.Println(port, "ini port")
	if dbPort, err = strconv.Atoi(port); err != nil {
		log.Err(err).Msg("failed to parse database port")
		os.Exit(1)
	}

	//encoderConfig := ecszap.NewDefaultEncoderConfig()
	//core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger() // Initializing zerolog logger

	dbConfig := db.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     dbPort,
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DbName:   os.Getenv("POSTGRES_DB"),
		Logger:   logger,
	}
	logger.Info().Interface("config", &dbConfig).Msg("config:")
	dbInstance, err := db.Init(dbConfig)
	if err != nil {
		logger.Err(err).Msg("Connection failed")
		os.Exit(1)
	}
	logger.Info().Msg("Database connection established")

	esClient, err := elasticsearch.NewDefaultClient()
	if err != nil {
		logger.Err(err).Msg("Connection failed")
		os.Exit(1)
	}

	h := handler.New(dbInstance, esClient, logger)
	router := gin.Default()
	rg := router.Group("/v1")
	h.Register(rg)
	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}

	count := 0
	for {
		if rand.Float32() > 0.8 {
			logger.Error().Int("count", count).Msg("oops...something is wrong")
		} else {
			logger.Info().Int("count", count).Msg("everything is fine")
		}
		count++
		time.Sleep(time.Second * 2)
	}
}
