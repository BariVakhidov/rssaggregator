package routes

import (
	"github.com/BariVakhidov/rssaggregator/api/handlers"
	"github.com/go-chi/chi"
)

func InitV1Router(apiCfg *handlers.ApiConfig) *chi.Mux {
	v1Router := chi.NewRouter()
	v1Router.Get("/health", handlers.HandlerReadiness)
	v1Router.Get("/err", handlers.HandlerErr)

	v1Router.Post("/users", apiCfg.HandlerCreateUser)
	v1Router.Get("/users", apiCfg.MiddlewareAuth(apiCfg.HandlerGetUserByAPIKey))

	v1Router.Post("/feeds", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.HandlerGetFeeds)

	v1Router.Get("/posts", apiCfg.MiddlewareAuth(apiCfg.HandlerGetPostsForUser))

	v1Router.Post("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.HandlerFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.MiddlewareAuth(apiCfg.HandlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.MiddlewareAuth(apiCfg.HandlerDeleteFeedFollow))

	return v1Router
}
