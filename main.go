package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"gymshark/db"
	"gymshark/handlers"
)

type Config struct {
	Host   string
	Port   string
	DbName string
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Getenv); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, out io.Writer, getenv func(string) string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Init app config
	cfg := Config{
		Host:   withDefault(getenv("APP_HOST"), "0.0.0.0"),
		Port:   withDefault(getenv("APP_PORT"), "8081"),
		DbName: withDefault(getenv("APP_DB_NAME"), "dev.db"),
	}

	// Init database
	database, err := db.InitDatabase(ctx, cfg.DbName)
	if err != nil {
		return fmt.Errorf("connecting to database: %s", err)
	}
	defer func() {
		log.Println("stopping database")
		sqlDB, _ := database.DB()
		sqlDB.Close()
	}()

	packSizeRepo := db.NewPackSizeRepository(database)

	// Init app logger
	logger := log.New(out, "", log.LstdFlags)

	// Start server
	srv := NewServer(logger, packSizeRepo)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: srv,
	}
	go func() {
		log.Printf("server started on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Println("got interrupt signal. shutdown gracefully")
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}

// withDefault returns want if not empty else will return default def
func withDefault(want, def string) string {
	if want != "" {
		return want
	}
	return def
}

func NewServer(logger *log.Logger, packSizeRepo *db.PackSizeRepository) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, packSizeRepo, logger)

	var handler http.Handler = mux
	return handler
}

func addRoutes(mux *http.ServeMux, packSizeRepo *db.PackSizeRepository, logger *log.Logger) {
	mux.Handle("/health", handlers.HealthHandler(logger))

	mux.Handle("POST /api/save-pack-sizes", handlers.SavePackSizesHandler(packSizeRepo))
	mux.Handle("GET /api/pack-sizes", handlers.GetPackSizes(packSizeRepo))
	mux.Handle("GET /api/packs/{itemsNo}", handlers.CalculatePacks(packSizeRepo, logger))

	mux.Handle("/", handlers.UiHandler())
}
