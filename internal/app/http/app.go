package httpapp

import (
	"context"
	"log/slog"
	"net/http"

	apiv1 "github.com/BariVakhidov/rssaggregator/api/v1"
	feedv1 "github.com/BariVakhidov/rssaggregator/api/v1/handlers/feed"
	postv1 "github.com/BariVakhidov/rssaggregator/api/v1/handlers/post"
	userv1 "github.com/BariVakhidov/rssaggregator/api/v1/handlers/user"
	auth "github.com/BariVakhidov/rssaggregator/internal/delivery/http/middleware"
	httpserver "github.com/BariVakhidov/rssaggregator/internal/delivery/http/server"
	"github.com/BariVakhidov/rssaggregator/internal/lib/token"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type AppOpts struct {
	Log           *slog.Logger
	TokenVerifier token.Verifier
	Port          int
}

type App struct {
	AppOpts
	server *httpserver.Server
}

func New(opts AppOpts,
	feedService feedv1.FeedService,
	userService userv1.UserService,
	postService postv1.Service,
) *App {
	chiRouter := chi.NewRouter()
	authMiddleware := auth.New(opts.Log, opts.TokenVerifier)

	chiRouter.Use(cors.Handler(cors.Options{
		AllowCredentials: false,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowedOrigins:   []string{"https://*", "http://*"},
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,
	}))
	chiRouter.Use(middleware.RequestID)
	chiRouter.Use(middleware.RealIP)
	chiRouter.Use(middleware.Logger)
	chiRouter.Use(middleware.Recoverer)
	chiRouter.Use(authMiddleware)
	chiRouter.Use(middleware.Heartbeat("/ping"))

	apiV1 := apiv1.New(opts.Log, feedService, userService, postService)
	chiRouter.Mount("/v1", apiV1.Router)

	server := httpserver.New(opts.Log, opts.Port, chiRouter)

	return &App{
		server:  server,
		AppOpts: opts,
	}
}

func (a *App) MustRun() {
	a.server.MustRun()
}

func (a *App) Stop(ctx context.Context) {
	a.server.Stop(ctx)
}
