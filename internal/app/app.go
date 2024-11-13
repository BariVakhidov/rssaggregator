package app

import (
	"context"
	"log/slog"
	"time"

	kafkaconsumerapp "github.com/BariVakhidov/rssaggregator/internal/app/event_consumer"
	httpapp "github.com/BariVakhidov/rssaggregator/internal/app/http"
	rssfetcherapp "github.com/BariVakhidov/rssaggregator/internal/app/rss_fetcher"
	storageapp "github.com/BariVakhidov/rssaggregator/internal/app/storage"
	jaegerapp "github.com/BariVakhidov/rssaggregator/internal/app/tracing/jaeger"
	"github.com/BariVakhidov/rssaggregator/internal/clients/sso/grpc"
	"github.com/BariVakhidov/rssaggregator/internal/config"
	"github.com/BariVakhidov/rssaggregator/internal/lib/logger/sl"
	jwtverifier "github.com/BariVakhidov/rssaggregator/internal/lib/token/jwt"
	uuidgenerator "github.com/BariVakhidov/rssaggregator/internal/lib/uuid/generator"
	"github.com/BariVakhidov/rssaggregator/internal/service/feed"
	postservice "github.com/BariVakhidov/rssaggregator/internal/service/post"
	userservice "github.com/BariVakhidov/rssaggregator/internal/service/user"
	"go.opentelemetry.io/otel"
)

type App struct {
	log            *slog.Logger
	httpApp        *httpapp.App
	rssFetcher     *rssfetcherapp.App
	userConsumer   *kafkaconsumerapp.App
	storage        *storageapp.App
	tracerProvider *jaegerapp.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	brokers := []string{"host.docker.internal:29092"}
	topic := "user_created"

	storage := storageapp.MustCreateApp("postgresql://postgres:password@db:5432/rssaggregator?sslmode=disable", log)

	grpcAuth, err := grpc.New(log, "host.docker.internal:8080", 5, time.Millisecond*100)
	if err != nil {
		log.Error("failed to create auth client", sl.Err(err))
	}

	tracerProvider, err := jaegerapp.New(log, "host.docker.internal:4318")
	if err != nil {
		log.Error("failed to create tracerApp", sl.Err(err))
	}
	tracer := otel.Tracer("rssaggregator")

	uuidGenerator := uuidgenerator.New()
	userService := userservice.New(log, storage.Storage, grpcAuth, uuidGenerator, tracer)
	feedService := feed.New(log, storage.Storage)
	postService := postservice.New(log, storage.Storage)

	rssFetcher := rssfetcherapp.New(log, storage.Storage, storage.Storage, time.Second*10)

	userConsumer := kafkaconsumerapp.New(log, storage.Storage, brokers, topic)

	tokenVerifier := jwtverifier.New("rss_secret")

	httpAppOpts := httpapp.AppOpts{
		Log:           log,
		TokenVerifier: tokenVerifier,
		Port:          cfg.Port,
	}
	httpApp := httpapp.New(httpAppOpts, feedService, userService, postService)

	return &App{
		log:            log,
		httpApp:        httpApp,
		rssFetcher:     rssFetcher,
		userConsumer:   userConsumer,
		storage:        storage,
		tracerProvider: tracerProvider,
	}
}

func (a *App) MustRun() {
	go a.httpApp.MustRun()
	go a.rssFetcher.RunScrapper(100, time.Second*5)
	go a.userConsumer.MustRun(context.TODO())
}

func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	a.httpApp.Stop(ctx)
	a.userConsumer.Stop()
	a.storage.Stop()
	a.tracerProvider.Stop(ctx)
}
