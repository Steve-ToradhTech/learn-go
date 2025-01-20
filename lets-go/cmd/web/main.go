package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// Create struct for dependency injection
type application struct {
	logger *slog.Logger
}

func main() {

	addr := flag.String("addr", ":4000", "[:PORT] : HTTP network address")
	debug := flag.Bool("debug", false, "[True/False] : See debug logging.")
	flag.Parse()

	// Checks and sets the logging level if debug is set to True.
	logLevel := setDebugLevel(debug)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: *debug,
		Level:     logLevel,
	}))

	app := &application{
		logger: logger,
	}

	logger.Info("starting server", slog.Any("addr", *addr))

	err := http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}
