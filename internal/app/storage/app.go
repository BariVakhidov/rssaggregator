package storageapp

import (
	"log/slog"

	"github.com/BariVakhidov/rssaggregator/internal/storage/postgres"
)

type App struct {
	Storage *postgres.Storage
	log     *slog.Logger
	addr    string
}

func MustCreateApp(addr string, log *slog.Logger) *App {
	postgres, err := postgres.New(log, addr)
	if err != nil {
		panic(err)
	}

	return &App{
		log:     log,
		Storage: postgres,
		addr:    addr,
	}
}

func (a *App) Stop() {
	const op = "storageapp.Stop"
	a.log.With(slog.String("op", op)).Info("stopping storage app")
	a.Storage.ClosePool()
}
