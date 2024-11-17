package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	ssov1 "github.com/BariVakhidov/ssoprotos/gen/go/sso"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	authClient ssov1.AuthClient
	log        *slog.Logger
}

func New(log *slog.Logger, addr string, maxRetries int, timeout time.Duration) (*Client, error) {
	const op = "sso.grpcClient.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(maxRetries)),
		grpcretry.WithPerRetryTimeout(timeout),
	}
	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.PayloadSent, grpclog.PayloadReceived),
	}

	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return &Client{
		authClient: ssov1.NewAuthClient(cc),
		log:        log,
	}, nil
}

func (c *Client) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	const op = "sso.grpcClient.Register"

	resp, err := c.authClient.Register(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return resp, nil
}

func (c *Client) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	const op = "sso.grpcClient.Login"

	resp, err := c.authClient.Login(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return resp, nil
}

func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, level grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(level), msg, fields...)
	})
}
