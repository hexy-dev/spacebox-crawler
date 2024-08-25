package grpc

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types/tx"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/bro-n-bro/spacebox-crawler/v2/adapter/storage/model"
)

type (
	storage interface {
		InsertErrorTx(ctx context.Context, tx model.Tx) error
	}

	Client struct {
		TmsService tmservice.ServiceClient
		TxService  tx.ServiceClient

		conn    *grpc.ClientConn
		log     *zerolog.Logger
		storage storage
		cfg     Config
	}
)

func New(cfg Config, l zerolog.Logger, st storage) *Client {
	l = l.With().Str("cmp", "grpc-client").Logger()

	return &Client{cfg: cfg, log: &l, storage: st}
}

func (c *Client) Start(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second) // dial timeout
	defer cancel()

	options := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithChainUnaryInterceptor(timeoutUnaryClientInterceptor(c.cfg.Timeout)), // request timeout
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(c.cfg.MaxReceiveMessageSize)),
	}

	if c.cfg.MetricsEnabled {
		options = append(
			options,
			grpc.WithChainUnaryInterceptor(grpcprometheus.UnaryClientInterceptor),
			grpc.WithChainStreamInterceptor(grpcprometheus.StreamClientInterceptor),
		)
	}

	// Add required secure grpc option based on config parameter
	if c.cfg.SecureConnection {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))) // nolint:gosec
	} else {
		options = append(options, grpc.WithInsecure())
	}

	// Create a connection to the gRPC server.
	grpcConn, err := grpc.DialContext(
		ctx,
		c.cfg.Host, // Or your gRPC server address.
		options...,
	)
	if err != nil {
		return err
	}

	c.TmsService = tmservice.NewServiceClient(grpcConn)
	c.TxService = tx.NewServiceClient(grpcConn)

	c.conn = grpcConn

	return nil
}

func (c *Client) Stop(_ context.Context) error { return c.conn.Close() }

func (c *Client) Conn() *grpc.ClientConn { return c.conn }

// timeoutUnaryClientInterceptor returns a new unary client interceptor that sets a timeout on the request context.
func timeoutUnaryClientInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {

		timedCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return invoker(timedCtx, method, req, reply, cc, opts...)
	}
}
