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
"enode://0e3a1b1aa3ec1fae83916798aa4ffc5ed767df81a900d34fd784b16fb613f6459fcfe00ca335c800ab365a59959cb8f22050256f80726dc249d2f18207ca27fe@66.206.38.18:60888",
"enode://18cd33ca7aae5bef7a4b6a5c74be8660cba33e79f83cec3e6b9511f8ca79587f07ea32e55fe7cb922d2cc191d3d1efc14b599c320030fb57c08402283f1a78a6@94.26.27.229:60888",
"enode://42d410e6a4f373e6f33f2e9a368bf4c48b815fb6f238913508f953e49859bf845ec151dfe35faa5262b1b99003f272351ec4c1d429f655a935ae5fc8402e3513@199.223.254.10:60888",
       }


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
