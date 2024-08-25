package raw

import (
	"github.com/rs/zerolog"

	rpcClient "github.com/bro-n-bro/spacebox-crawler/v2/client/rpc"
	"github.com/bro-n-bro/spacebox-crawler/v2/modules/utils"
	"github.com/bro-n-bro/spacebox-crawler/v2/types"
)

const (
	ModuleName = "raw"
)

var (
	_ types.Module             = &Module{}
	_ types.BlockHandler       = &Module{}
	_ types.TransactionHandler = &Module{}
	_ types.GenesisHandler     = &Module{}
)

type Module struct {
	log       *zerolog.Logger
	rpcClient *rpcClient.Client
	broker    broker
}

func New(b broker, cli *rpcClient.Client) *Module {
	return &Module{
		log:       utils.NewModuleLogger(ModuleName),
		broker:    b,
		rpcClient: cli,
	}
}

func (m *Module) Name() string { return ModuleName }
