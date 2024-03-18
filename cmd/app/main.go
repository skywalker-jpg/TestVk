package main

import (
	"TestVK/internal/config"
	"TestVK/internal/db"
	"TestVK/internal/filmoteka"
	"TestVK/internal/logger"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	config, err := config.NewConfig("config.yaml")
	if err != nil {
		return err
	}

	logger, err := logger.New(config.Logger)
	if err != nil {
		return err
	}

	db, err := db.Connection(config.DB)
	if err != nil {
		logger.Error("Error creating db", err.Error())
		return err
	}

	defer db.Close()

	filmoteka := filmoteka.NewFilmoteka(db, config, logger)
	filmoteka.Api()
	return nil
}
