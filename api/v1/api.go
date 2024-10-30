package apiv1

import (
	"log/slog"

	"github.com/BariVakhidov/rssaggregator/api/v1/handlers/feed"
	"github.com/BariVakhidov/rssaggregator/api/v1/handlers/post"
	"github.com/BariVakhidov/rssaggregator/api/v1/handlers/user"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Api struct {
	Router *chi.Mux
}

func New(
	log *slog.Logger,
	feedService feed.FeedService,
	userService user.UserService,
	postService post.Service,
) *Api {
	v1Router := chi.NewRouter()

	feedHandler := feed.New(log, feedService)
	userHandler := user.New(log, userService)
	postHandler := post.New(log, postService)

	v1Router.Get("/swagger/*", httpSwagger.WrapHandler)
	v1Router.Post("/users", userHandler.CreateUser)
	v1Router.Post("/users/login", userHandler.Login)
	v1Router.Get("/users", userHandler.User)

	v1Router.Post("/feeds", feedHandler.CreateFeed)
	v1Router.Get("/feeds", feedHandler.Feeds)

	v1Router.Get("/posts", postHandler.Posts)

	v1Router.Post("/feed-follows", feedHandler.CreateFeedFollow)
	v1Router.Get("/feed-follows", feedHandler.FeedFollows)
	v1Router.Delete("/feed-follows/{feedFollowID}", feedHandler.DeleteFeedFollow)

	return &Api{
		Router: v1Router,
	}
}
