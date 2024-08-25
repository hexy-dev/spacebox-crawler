package rpc

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	tmHttp "github.com/tendermint/tendermint/rpc/client/http"
	jsonrpcclient "github.com/tendermint/tendermint/rpc/jsonrpc/client"
)

type Client struct {
	*jsonrpcclient.WSClient

	RPCClient *tmHttp.HTTP

	cfg Config
}

func New(cfg Config) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Start(ctx context.Context) error {
	httpClient, err := jsonrpcclient.DefaultHTTPClient(c.cfg.Host)
	if err != nil {
		return err
	}

	httpClient.Timeout = c.cfg.Timeout

	if c.cfg.MetricsEnabled {
		httpClient.Transport = promhttp.InstrumentRoundTripperInFlight(inFlightGauge,
			promhttp.InstrumentRoundTripperCounter(counter,
				promhttp.InstrumentRoundTripperDuration(histVec, http.DefaultTransport)),
		)
	}

	c.RPCClient, err = tmHttp.NewWithClient(c.cfg.Host, httpClient)
	if err != nil {
		return err
	}

	if err = c.RPCClient.Start(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Client) Stop(_ context.Context) error {
	if err := c.RPCClient.Stop(); err != nil {
		return err
	}

	return nil
}
