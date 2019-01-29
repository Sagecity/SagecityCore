// Copyright 2016 The sagecity Authors
// This file is part of the sagecity library.
//
// The sagecity library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The sagecity library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the sagecity library. If not, see <http://www.gnu.org/licenses/>.

// Package les implements the Light Sagecity Subprotocol.
package les

import (
	"fmt"
	"sync"
	"time"

	"github.com/zdbrig/sagecity/accounts"
	"github.com/zdbrig/sagecity/common"
	"github.com/zdbrig/sagecity/common/hexutil"
	"github.com/zdbrig/sagecity/consensus"
	"github.com/zdbrig/sagecity/core"
	"github.com/zdbrig/sagecity/core/types"
	"github.com/zdbrig/sagecity/eth"
	"github.com/zdbrig/sagecity/eth/downloader"
	"github.com/zdbrig/sagecity/eth/filters"
	"github.com/zdbrig/sagecity/eth/gasprice"
	"github.com/zdbrig/sagecity/ethdb"
	"github.com/zdbrig/sagecity/event"
	"github.com/zdbrig/sagecity/internal/ethapi"
	"github.com/zdbrig/sagecity/light"
	"github.com/zdbrig/sagecity/log"
	"github.com/zdbrig/sagecity/node"
	"github.com/zdbrig/sagecity/p2p"
	"github.com/zdbrig/sagecity/p2p/discv5"
	"github.com/zdbrig/sagecity/params"
	rpc "github.com/zdbrig/sagecity/rpc"
)

type LightSagecity struct {
	odr         *LesOdr
	relay       *LesTxRelay
	chainConfig *params.ChainConfig
	// Channel for shutting down the service
	shutdownChan chan bool
	// Handlers
	peers           *peerSet
	txPool          *light.TxPool
	blockchain      *light.LightChain
	protocolManager *ProtocolManager
	serverPool      *serverPool
	reqDist         *requestDistributor
	retriever       *retrieveManager
	// DB interfaces
	chainDb ethdb.Database // Block chain database

	ApiBackend *LesApiBackend

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	networkId     uint64
	netRPCService *ethapi.PublicNetAPI

	wg sync.WaitGroup
}

func New(ctx *node.ServiceContext, config *eth.Config) (*LightSagecity, error) {
	chainDb, err := eth.CreateDB(ctx, config, "lightchaindata")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, isCompat := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !isCompat {
		return nil, genesisErr
	}
	log.Info("Initialised chain configuration", "config", chainConfig)

	peers := newPeerSet()
	quitSync := make(chan struct{})

	eth := &LightSagecity{
		chainConfig:    chainConfig,
		chainDb:        chainDb,
		eventMux:       ctx.EventMux,
		peers:          peers,
		reqDist:        newRequestDistributor(peers, quitSync),
		accountManager: ctx.AccountManager,
		engine:         eth.CreateConsensusEngine(ctx, config, chainConfig, chainDb),
		shutdownChan:   make(chan bool),
		networkId:      config.NetworkId,
	}

	eth.relay = NewLesTxRelay(peers, eth.reqDist)
	eth.serverPool = newServerPool(chainDb, quitSync, &eth.wg)
	eth.retriever = newRetrieveManager(peers, eth.reqDist, eth.serverPool)
	eth.odr = NewLesOdr(chainDb, eth.retriever)
	if eth.blockchain, err = light.NewLightChain(eth.odr, eth.chainConfig, eth.engine); err != nil {
		return nil, err
	}
	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		eth.blockchain.SetHead(compat.RewindTo)
		core.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}

	eth.txPool = light.NewTxPool(eth.chainConfig, eth.blockchain, eth.relay)
	if eth.protocolManager, err = NewProtocolManager(eth.chainConfig, true, config.NetworkId, eth.eventMux, eth.engine, eth.peers, eth.blockchain, nil, chainDb, eth.odr, eth.relay, quitSync, &eth.wg); err != nil {
		return nil, err
	}
	eth.ApiBackend = &LesApiBackend{eth, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.GasPrice
	}
	eth.ApiBackend.gpo = gasprice.NewOracle(eth.ApiBackend, gpoParams)
	return eth, nil
}

func lesTopic(genesisHash common.Hash) discv5.Topic {
	return discv5.Topic("LES@" + common.Bytes2Hex(genesisHash.Bytes()[0:8]))
}

type LightDummyAPI struct{}

// Etherbase is the address that mining rewards will be send to
func (s *LightDummyAPI) Etherbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("not supported")
}

// Coinbase is the address that mining rewards will be send to (alias for Etherbase)
func (s *LightDummyAPI) Coinbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("not supported")
}

// Hashrate returns the POW hashrate
func (s *LightDummyAPI) Hashrate() hexutil.Uint {
	return 0
}

// Mining returns an indication if this node is currently mining.
func (s *LightDummyAPI) Mining() bool {
	return false
}

// APIs returns the collection of RPC services the sagecity package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *LightSagecity) APIs() []rpc.API {
	return append(ethapi.GetAPIs(s.ApiBackend), []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   &LightDummyAPI{},
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.ApiBackend, true),
			Public:    true,
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *LightSagecity) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *LightSagecity) BlockChain() *light.LightChain      { return s.blockchain }
func (s *LightSagecity) TxPool() *light.TxPool              { return s.txPool }
func (s *LightSagecity) Engine() consensus.Engine           { return s.engine }
func (s *LightSagecity) LesVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *LightSagecity) Downloader() *downloader.Downloader { return s.protocolManager.downloader }
func (s *LightSagecity) EventMux() *event.TypeMux           { return s.eventMux }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *LightSagecity) Protocols() []p2p.Protocol {
	return s.protocolManager.SubProtocols
}

// Start implements node.Service, starting all internal goroutines needed by the
// Sagecity protocol implementation.
func (s *LightSagecity) Start(srvr *p2p.Server) error {
	log.Warn("Light client mode is an experimental feature")
	s.netRPCService = ethapi.NewPublicNetAPI(srvr, s.networkId)
	s.serverPool.start(srvr, lesTopic(s.blockchain.Genesis().Hash()))
	s.protocolManager.Start()
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Sagecity protocol.
func (s *LightSagecity) Stop() error {
	s.odr.Stop()
	s.blockchain.Stop()
	s.protocolManager.Stop()
	s.txPool.Stop()

	s.eventMux.Stop()

	time.Sleep(time.Millisecond * 200)
	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
