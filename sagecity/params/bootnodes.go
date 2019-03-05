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
	"enode://6c6d558cceb1cc7e6b695ced7b88e86e891df08d102b3a08e3b09c5327b774d6b02650ed367c17e201ef303fb1971fe22f6774a61088da42ce420d31db235655@66.206.38.18:60888",// Amsterdam
        "enode://b997f95a21bb7de6bf00fceb6e747daecbe558d82f1ff59e863c9b9a139a4f568e91c5c60384726130c7ad11327d012d9cf934ffbb9fbb15d800910c5914aed0@94.26.27.229:60888",

       "enode://5438df8fcd6b49f693bfa454afac2e36052dddd35fbce81593062ef600c059a3eb269c1aa927f7a4a243e12d6a49c0ebd20ef449125e718bd8cf6546292c9064@199.223.254.101:60888",
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
