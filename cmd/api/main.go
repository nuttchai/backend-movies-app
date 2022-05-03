package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

// NOTE: we create a new type of application where we'll put the info we need to share with our handler and other components in the application
type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// NOTE: Flag lets you create options when running command-line programs
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production)")
	flag.Parse()

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server on port", cfg.port)

	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
