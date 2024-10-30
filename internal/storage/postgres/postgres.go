package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/BariVakhidov/rssaggregator/internal/converter"
	"github.com/BariVakhidov/rssaggregator/internal/database"
	"github.com/BariVakhidov/rssaggregator/internal/lib/logger/sl"
	"github.com/BariVakhidov/rssaggregator/internal/model"
	"github.com/BariVakhidov/rssaggregator/internal/storage"
	dbConverter "github.com/BariVakhidov/rssaggregator/internal/storage/converter"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var pgOnce sync.Once

type Storage struct {
	queries *database.Queries
	dbPool  *pgxpool.Pool
	log     *slog.Logger
}

func New(log *slog.Logger, dbAddr string) (*Storage, error) {

	const op = "storage.postgres.New"

	var (
		dbPool *pgxpool.Pool
		err    error
	)

	//single instance of the db
	pgOnce.Do(func() {
		dbPool, err = pgxpool.New(context.Background(), dbAddr)
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	queries := database.New(dbPool)

	return &Storage{queries: queries, dbPool: dbPool, log: log}, nil
}

func (s *Storage) CreateUser(ctx context.Context, userID uuid.UUID, email string) (userCreated model.User, err error) {
	const op = "storage.postgres.CreateUser"
	log := s.log.With(slog.String("op", op))

	tx, err := s.dbPool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				log.Error("rollback failed", sl.Err(rollbackErr))
			}
			return
		}

		if commitErr := tx.Commit(ctx); commitErr != nil {
			log.Error("commit failed", sl.Err(commitErr))
			err = fmt.Errorf("%s: %w", op, commitErr)
		}
	}()

	pendingUser, err := s.queries.WithTx(tx).PendingUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	userParams := database.CreateUserParams{
		ID:        userID,
		Email:     email,
		Name:      pendingUser.Name,
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	user, err := s.queries.WithTx(tx).CreateUser(ctx, userParams)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == storage.UniqueViolationCode {
			return model.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.queries.WithTx(tx).DeletePendingUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return converter.DatabaseUserToUser(user), nil
}

func (s *Storage) CreateFeed(ctx context.Context, feedInfo model.FeedInfo) (*model.Feed, error) {
	const op = "storage.postgres.CreateFeed"

	params := database.CreateFeedParams{
		ID:   uuid.New(),
		Name: feedInfo.Name,
		Url:  feedInfo.Url,
	}

	feed, err := s.queries.CreateFeed(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == storage.UniqueViolationCode {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrFeedExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	convertedFeed := converter.DatabaseFeedToFeed(feed)
	return &convertedFeed, nil
}

func (s *Storage) Feeds(ctx context.Context) ([]model.Feed, error) {
	const op = "storage.postgres.Feeds"

	feeds, err := s.queries.GetFeeds(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.DatabaseFeedsToFeeds(feeds), nil
}

func (s *Storage) FeedsByUser(ctx context.Context, userId uuid.UUID) ([]model.Feed, error) {
	const op = "storage.postgres.FeedsByUser"

	feeds, err := s.queries.GetFeedsByUser(ctx, userId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == storage.ForeignKeyViolationCode {
				return nil, fmt.Errorf("%s: %w", op, storage.ErrForeignKeyViolation)
			}
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.DatabaseFeedsToFeeds(feeds), nil
}

func (s *Storage) Posts(ctx context.Context, userId uuid.UUID, limit int) ([]model.Post, error) {
	const op = "storage.postgres.Posts"

	params := database.GetPostsForUserParams{
		UserID: userId,
		Limit:  int32(limit),
	}

	posts, err := s.queries.GetPostsForUser(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == storage.ForeignKeyViolationCode {
				return nil, fmt.Errorf("%s: %w", op, storage.ErrForeignKeyViolation)
			}

			if pgErr.Code == storage.InvalidTextRepresentation {
				return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidParams)
			}

		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.DatabasePostsToPosts(posts), nil
}

func (s *Storage) CreateFeedFollow(ctx context.Context, userId uuid.UUID, feedId uuid.UUID) (*model.FeedFollow, error) {
	const op = "storage.postgres.CreateFeedFollow"

	params := database.CreateFeedFollowParams{
		ID:     uuid.New(),
		FeedID: feedId,
		UserID: userId,
	}

	feedFollow, err := s.queries.CreateFeedFollow(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == storage.UniqueViolationCode {
				return nil, fmt.Errorf("%s: %w", op, storage.ErrFeedFollowExists)
			}

			if pgErr.Code == storage.ForeignKeyViolationCode {
				return nil, fmt.Errorf("%s: %w", op, storage.ErrForeignKeyViolation)
			}
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	convertedFeedFollow := converter.DatabaseFeedFollowToFeedFollow(feedFollow)
	return &convertedFeedFollow, nil
}

func (s *Storage) DeleteFeedFollow(ctx context.Context, userId uuid.UUID, followId uuid.UUID) error {
	const op = "storage.postgres.DeleteFeedFollow"

	params := database.DeleteFeedFollowParams{
		ID:     followId,
		UserID: userId,
	}

	if err := s.queries.DeleteFeedFollow(ctx, params); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == storage.ForeignKeyViolationCode {
			return fmt.Errorf("%s: %w", op, storage.ErrForeignKeyViolation)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) FeedFollows(ctx context.Context, info model.FeedFollowsInfo) ([]model.FeedFollow, error) {
	const op = "storage.postgres.FeedFollows"

	follows, err := s.queries.GetFeedFollows(ctx, dbConverter.FeedFollowsInfoToDBFeedFollowsInfo(info))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == storage.ForeignKeyViolationCode {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrForeignKeyViolation)
		}

		if pgErr.Code == storage.InvalidTextRepresentation {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidParams)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.DatabaseFeedFollowsToFeedFollows(follows), nil
}

func (s *Storage) MarkFeedAsFetched(ctx context.Context, feedId uuid.UUID) (*model.Feed, error) {
	const op = "storage.postgres.MarkFeedAsFetched"

	feed, err := s.queries.MarkFeedAsFetched(ctx, feedId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	serviceFeed := converter.DatabaseFeedToFeed(feed)
	return &serviceFeed, nil
}

func (s *Storage) NextFeedsToFetch(ctx context.Context, limit int) ([]model.Feed, error) {
	const op = "storage.postgres.NextFeedsToFetch"

	feeds, err := s.queries.GetNextFeedsToFetch(ctx, int32(limit))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == storage.InvalidTextRepresentation {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrInvalidParams)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return converter.DatabaseFeedsToFeeds(feeds), nil
}

func (s *Storage) SavePendingUser(ctx context.Context, pendingUserInfo model.PendingUserInfo) (model.PendingUser, error) {
	const op = "storage.postgres.SavePendingUser"

	params := dbConverter.ServicePendingUserInfoToDBParams(pendingUserInfo)
	dbUser, err := s.queries.SavePendingUser(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == storage.UniqueViolationCode {
			return model.PendingUser{}, fmt.Errorf("%s: %w", op, storage.ErrFeedFollowExists)
		}

		return model.PendingUser{}, fmt.Errorf("%s: %w", op, err)
	}

	return converter.DatabasePendingUserToPendingUser(dbUser), nil
}

func (s *Storage) CreatePost(ctx context.Context, postInfo model.CreatePostInfo) (post *model.Post, err error) {
	const op = "storage.postgres.CreatePost"
	log := s.log.With(slog.String("op", op))

	params := dbConverter.CreatePostInfoToDBCreatePostInfo(postInfo)

	tx, err := s.dbPool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback(ctx)
			if rollbackErr != nil {
				log.Error("rollback failed", sl.Err(rollbackErr))
			}
			return
		}

		if commitErr := tx.Commit(ctx); commitErr != nil {
			log.Error("commit failed", sl.Err(commitErr))
			err = fmt.Errorf("%s: %w", op, commitErr)
		}
	}()

	_, err = s.queries.WithTx(tx).GetPost(ctx, postInfo.Url)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err == nil {
		return nil, fmt.Errorf("%s: %w", op, storage.ErrPostExists)
	}

	dbPost, err := s.queries.WithTx(tx).CreatePost(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == storage.UniqueViolationCode {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrPostExists)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	convertedPost := converter.DatabasePostToPost(dbPost)
	return &convertedPost, nil
}

func (s *Storage) ChangePendingUserStatus(ctx context.Context, userId uuid.UUID, status string) (model.PendingUser, error) {
	const op = "storage.postgres.ChangePendingUserStatus"

	params := database.ChangeStatusOfPendingUserParams{
		Status: status,
		ID:     userId,
	}
	dbUser, err := s.queries.ChangeStatusOfPendingUser(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PendingUser{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}

		return model.PendingUser{}, fmt.Errorf("%s: %w", op, err)
	}

	return converter.DatabasePendingUserToPendingUser(dbUser), nil
}

func (s *Storage) PendingUserByEmail(ctx context.Context, userEmail string) (model.PendingUser, error) {
	const op = "storage.postgres.PendingUserByEmail"

	dbUser, err := s.queries.PendingUserByEmail(ctx, userEmail)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.PendingUser{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}

		return model.PendingUser{}, fmt.Errorf("%s: %w", op, err)
	}

	return converter.DatabasePendingUserToPendingUser(dbUser), nil
}

func (s *Storage) GetUser(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	const op = "storage.postgres.GetUser"

	dbUser, err := s.queries.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	convertedUser := converter.DatabaseUserToUser(dbUser)

	return &convertedUser, nil
}

func (s *Storage) ClosePool() {
	s.dbPool.Close()
}
