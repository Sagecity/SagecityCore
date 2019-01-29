// Copyright 2015 The sagecity Authors
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

package core

import (
	"fmt"
	"math/big"

	"github.com/zdbrig/sagecity/common"
	"github.com/zdbrig/sagecity/consensus"
	"github.com/zdbrig/sagecity/consensus/misc"
	"github.com/zdbrig/sagecity/core/state"
	"github.com/zdbrig/sagecity/core/types"
	"github.com/zdbrig/sagecity/core/vm"
	"github.com/zdbrig/sagecity/crypto"
	"github.com/zdbrig/sagecity/params"
)

// StateProcessor is a basic Processor, which takes care of transitioning
// state from one point to another.
//
// StateProcessor implements Processor.
type StateProcessor struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for block rewards
}

// NewStateProcessor initialises a new StateProcessor.
func NewStateProcessor(config *params.ChainConfig, bc *BlockChain, engine consensus.Engine) *StateProcessor {
	return &StateProcessor{
		config: config,
		bc:     bc,
		engine: engine,
	}
}

// Process processes the state changes according to the Sagecity rules by running
// the transaction messages using the statedb and applying any rewards to both
// the processor (coinbase) and any included uncles.
//
// Process returns the receipts and logs accumulated during the process and
// returns the amount of gas that was used in the process. If any of the
// transactions failed to execute due to insufficient gas it will return an error.
func (p *StateProcessor) Process(block *types.Block, statedb *state.StateDB, cfg vm.Config) (types.Receipts, []*types.Log, *big.Int, error) {
	var (
		receipts     types.Receipts
		totalUsedGas = big.NewInt(0)
		header       = block.Header()
		allLogs      []*types.Log
		gp           = new(GasPool).AddGas(block.GasLimit())
	)
	// Mutate the the block and state according to any hard-fork specs
	fmt.Print("Process if p.config.DAOForkSupport && p.config.DAOForkBlock != nil && p.config.DAOForkBlock.Cmp(block.Number()) == 0")
	if p.config.DAOForkSupport && p.config.DAOForkBlock != nil && p.config.DAOForkBlock.Cmp(block.Number()) == 0 {
		fmt.Print("Process misc.ApplyDAOHardFork(statedb)")
		misc.ApplyDAOHardFork(statedb)
	}
	// Iterate over and process the individual transactions
	fmt.Print("Process for i, tx := range block.Transactions() ")
	for i, tx := range block.Transactions() {
		fmt.Print("Process statedb.Prepare(tx.Hash(), block.Hash(), i)")
		statedb.Prepare(tx.Hash(), block.Hash(), i)
		fmt.Print("Process receipt, _, err := ApplyTransaction(p.config, p.bc, nil, gp, statedb, header, tx, totalUsedGas, cfg)")
		receipt, _, err := ApplyTransaction(p.config, p.bc, nil, gp, statedb, header, tx, totalUsedGas, cfg)
		fmt.Print("Process if err != nil")
		if err != nil {
			fmt.Print("Process return nil, nil, nil, err")
			return nil, nil, nil, err
		}
		fmt.Print("Process receipts = append(receipts, receipt)")
		receipts = append(receipts, receipt)
		fmt.Print("Process allLogs = append(allLogs, receipt.Logs...)")
		allLogs = append(allLogs, receipt.Logs...)
	}
	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	fmt.Print("Process p.engine.Finalize(p.bc, header, statedb, block.Transactions(), block.Uncles(), receipts)")
	p.engine.Finalize(p.bc, header, statedb, block.Transactions(), block.Uncles(), receipts)
	fmt.Print("Process return receipts, allLogs, totalUsedGas, nil")
	return receipts, allLogs, totalUsedGas, nil
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func ApplyTransaction(config *params.ChainConfig, bc *BlockChain, author *common.Address, gp *GasPool, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *big.Int, cfg vm.Config) (*types.Receipt, *big.Int, error) {
	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number))
	if err != nil {
		fmt.Print("########## erreur ligne 103 ################  " , *tx.To()  , " --- \n" );
		return nil, nil, err
	}
	// Create a new context to be used in the EVM environment
	fmt.Print("----108----> no error for ", *tx.To()  , " --- \n" )
	context := NewEVMContext(msg, header, bc, author)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	fmt.Print("ApplyTransaction vmenv := vm.NewEVM(context, statedb, config, cfg) return error")
	vmenv := vm.NewEVM(context, statedb, config, cfg)
	// Apply the transaction to the current state (included in the env)
	fmt.Print("ApplyTransaction _, gas, failed, err := ApplyMessage(vmenv, msg, gp) return error")
	_, gas, failed, err := ApplyMessage(vmenv, msg, gp)
	if err != nil {
		fmt.Print("########## erreur ligne 117 ################  " , *tx.To()  , " --- \n" );
		return nil, nil, err
	}
	fmt.Print("----121----> no error for ", *tx.To()  , " --- \n" )
	// Update the state with pending changes
	var root []byte
	fmt.Print("if config.IsByzantium(header.Number) ")
	if config.IsByzantium(header.Number) {
		fmt.Print("statedb.Finalise(true)")
		statedb.Finalise(true)
	} else {
		fmt.Print("root = statedb.IntermediateRoot(config.IsEIP158(header.Number)).Bytes()")
		root = statedb.IntermediateRoot(config.IsEIP158(header.Number)).Bytes()
	}
	fmt.Print("usedGas.Add(usedGas, gas)")
	usedGas.Add(usedGas, gas)
	// Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
	// based on the eip phase, we're passing wether the root touch-delete accounts.
	fmt.Print("receipt := types.NewReceipt(root, failed, usedGas)")
	receipt := types.NewReceipt(root, failed, usedGas)
	fmt.Print("receipt.TxHash = tx.Hash()")
	receipt.TxHash = tx.Hash()
	fmt.Print("receipt.GasUsed = new(big.Int).Set(gas)")
	receipt.GasUsed = new(big.Int).Set(gas)
	fmt.Print("if msg.To() == nil")
	// if the transaction created a contract, store the creation address in the receipt.
	if msg.To() == nil {
		fmt.Print("receipt.ContractAddress = crypto.CreateAddress(vmenv.Context.Origin, tx.Nonce())")
		receipt.ContractAddress = crypto.CreateAddress(vmenv.Context.Origin, tx.Nonce())
	}

	// Set the receipt logs and create a bloom for filtering
	fmt.Print("receipt.Logs = statedb.GetLogs(tx.Hash())")
	receipt.Logs = statedb.GetLogs(tx.Hash())
	fmt.Print("receipt.Bloom = types.CreateBloom(types.Receipts{receipt})")
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	fmt.Print("return receipt, gas, err")
	return receipt, gas, err
}
