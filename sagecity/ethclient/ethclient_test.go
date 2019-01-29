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

package ethclient

import "github.com/zdbrig/sagecity"

// Verify that Client implements the sagecity interfaces.
var (
	_ = sagecity.ChainReader(&Client{})
	_ = sagecity.TransactionReader(&Client{})
	_ = sagecity.ChainStateReader(&Client{})
	_ = sagecity.ChainSyncReader(&Client{})
	_ = sagecity.ContractCaller(&Client{})
	_ = sagecity.GasEstimator(&Client{})
	_ = sagecity.GasPricer(&Client{})
	_ = sagecity.LogFilterer(&Client{})
	_ = sagecity.PendingStateReader(&Client{})
	// _ = sagecity.PendingStateEventer(&Client{})
	_ = sagecity.PendingContractCaller(&Client{})
)
