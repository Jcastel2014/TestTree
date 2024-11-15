package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/Jcastel2014/test3/internal/data"
	_ "github.com/lib/pq"
)

const appVersion = "1.0.0"

type serverConfig struct {
	port int
	env  string
	db   struct {
		dsn string
	}

	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
}

type appDependencies struct {
	config   serverConfig
	logger   *slog.Logger
	bookclub data.BookClub
}

func openDB(settings serverConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", settings.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func main() {
	var settings serverConfig

	flag.IntVar(&settings.port, "port", 4000, "Server Port")
	flag.Float64Var(&settings.limiter.rps, "limiter-rps", 2, "Rate Limiter maximum requests per second")

	flag.IntVar(&settings.limiter.burst, "limiter-burst", 5, "Rate Limiter maximum burst")

	flag.BoolVar(&settings.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&settings.env, "env", "development", "Environment(Development|Staging|Production)")
	flag.StringVar(&settings.db.dsn, "db-dsn", "postgres://comments:fishsticks@localhost/comments?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(settings)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("database connection pool established")

	appInstance := &appDependencies{
		config:   settings,
		logger:   logger,
		bookclub: data.BookClub{DB: db},
	}

	// apiServer := &http.Server{
	// 	Addr:         fmt.Sprintf(":%d", settings.port),
	// 	Handler:      appInstance.routes(),
	// 	IdleTimeout:  time.Minute,
	// 	ReadTimeout:  5 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// 	ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	// }

	// logger.Info("starting server", "address", apiServer.Addr, "env", settings.env)
	// err = apiServer.ListenAndServe()
	// logger.Error(err.Error())
	// os.Exit(1)

	err = appInstance.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
