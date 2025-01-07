package rpc

import (
	"context"
	"net/http"

	httpCli "github.com/cometbft/cometbft/rpc/client/http"
	jsonrpcclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Client struct {
	*jsonrpcclient.WSClient
	*httpCli.WSEvents

	RPCClient *httpCli.HTTP

	cfg Config
}

func New(cfg Config) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Start(_ context.Context) error {
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

	c.RPCClient, err = httpCli.NewWithClient(c.cfg.Host, "/websocket", httpClient)
	if err != nil {
		return err
	}

	if err = c.RPCClient.Start(); err != nil {
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
