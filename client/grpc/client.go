package grpc

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/types/tx"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

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

	// nolint:staticcheck
	options := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithChainUnaryInterceptor(timeout.UnaryClientInterceptor(c.cfg.Timeout)), // request timeout
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(c.cfg.MaxReceiveMessageSize)),
	}

	if c.cfg.MetricsEnabled {
		cm := grpcprom.NewClientMetrics(
			grpcprom.WithClientHandlingTimeHistogram())

		prometheus.MustRegister(cm)

		options = append(
			options,
			grpc.WithChainUnaryInterceptor(cm.UnaryClientInterceptor()),
			grpc.WithChainStreamInterceptor(cm.StreamClientInterceptor()),
		)
	}

	// Add required secure grpc option based on config parameter
	if c.cfg.SecureConnection {
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))) // nolint:gosec
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Create a connection to the gRPC server.
	// nolint:staticcheck
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
