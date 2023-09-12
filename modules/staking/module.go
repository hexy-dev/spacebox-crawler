package staking

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/rs/zerolog"

	grpcClient "github.com/bro-n-bro/spacebox-crawler/client/grpc"
	tb "github.com/bro-n-bro/spacebox-crawler/pkg/mapper/to_broker"
	"github.com/bro-n-bro/spacebox-crawler/types"
)

const (
	moduleName = "staking"
)

var (
	_ types.Module            = &Module{}
	_ types.GenesisHandler    = &Module{}
	_ types.BlockHandler      = &Module{}
	_ types.MessageHandler    = &Module{}
	_ types.ValidatorsHandler = &Module{}
)

type (
	Module struct {
		log            *zerolog.Logger
		client         *grpcClient.Client
		broker         broker
		tbM            tb.ToBroker
		cdc            codec.Codec
		defaultDenom   string
		enabledModules []string // xxx fixme
	}
)

func New(b broker, cli *grpcClient.Client, tbM tb.ToBroker, cdc codec.Codec, modules []string, defaultDenom string) *Module {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("module", moduleName).Logger()

	return &Module{
		log:            &l,
		broker:         b,
		client:         cli,
		tbM:            tbM,
		cdc:            cdc,
		enabledModules: modules,
		defaultDenom:   defaultDenom,
	}
}

func (m *Module) Name() string { return moduleName }
