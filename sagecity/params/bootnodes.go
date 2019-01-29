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
	"enode://d291f49f4dafd0def9ae675ecc65c0cd72c8c7de7c28a22011b91df5ea39eeed2c3979cc59b256dc54511fa17217851aa3af4ffe93dababd50483e8d6741d3d4@66.206.38.18:60888",// Amsterdam
	"enode://9d331972c01d2ff45ca3f2b2945dcefd65a48611813fa45d544600edf064366165daf86553db013df5dfcb97afeed1400fae4b70d2318026bdad90428093ef2d@94.26.27.229:60888",
	"enode://8b04f78863c015b95f99bf6f3c96864000acbb117c75d85f99ff41c7ce5d0942a673027d81862b77faf0b93f8a41c26e1df28a62c090e7db43d6719ad1a44c59@199.223.254.101:60888"}

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
