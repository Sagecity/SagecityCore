// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

// SagecitynetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ethereum network.
var SagecitynetBootnodes = []string{
	// Ethereum Foundation Go Bootnodes
	"enode://ea2ae1a46a8a875901c66b3aee9dc31a46a857382c2999f2805377355f7fcc0f6a9663130e045b921d3510eba72a57167425b3caa0b90a042742d6d0470e4ea7@45.76.34.21:60888",// Amsterdam
	"enode://51395afd5ec5893409cd77887be1f99bf3d1e129930189ba04d4bc7d5f7421c435e93d098df00a7d8d8c69916e7ab11960df2ba3cba6b2cfa1f70e6e9b2a6d43@45.63.57.21:60888",
	"enode://c3272b65c1ced3dbfc38eaa70c04e9ab15a7dff6ee5d6a5c285fff882b34ca4b73c3d9e7b5deb438d2efe53e369671f6363e7d984831134f1218b5290092c977@104.156.227.21:60888",// New Jersey
	"enode://80fcf52a2ab06f0264b2083cecb3c6eee499e417c14ed6d20ceda9ebcfa3306a3aebc5396281dd5887bf8c375a823895c9e2b2b46835fef724723e3e45c2df9c@45.77.131.29:60888",// Tokoyo
	"enode://59b4a143f74f7e0b1ea7a74c78d65294bbf777f6442c392be35bb9f87cf2807583defb7a9f7d884d99978d1a3887c96ba05e0764f136a6406526e09d6067bb50@107.191.57.205:60888",
	"enode://589fe6ab2d7e061e9eaa21cf51e30b60b236b33346c532dd869189d01935e6a086e50a93e915d8d22f5eb4e7c68b7063394b6c24a8b7b095803d919acb6de7a7@207.180.198.242:60888"}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Ropsten test network.
var TestnetBootnodes = []string{

}

// RinkebyBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network.
var RinkebyBootnodes = []string{

}

// RinkebyV5Bootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network for the experimental RLPx v5 topic-discovery network.
var RinkebyV5Bootnodes = []string{

}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{
}
