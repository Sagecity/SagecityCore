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

package types

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/SagecityCore/sagecity/common"
	"github.com/SagecityCore/sagecity/crypto"
	"github.com/SagecityCore/sagecity/params"
)

var (
	ErrInvalidChainId = errors.New("invalid chain id for signer")
)

// sigCache is used to cache the derived sender and contains
// the signer used to derive it.
type sigCache struct {
	signer Signer
	from   common.Address
}

// MakeSigner returns a Signer based on the given chain config and block number.
func MakeSigner(config *params.ChainConfig, blockNumber *big.Int) Signer {
	println("MakeSigner var signer Signer")
	var signer Signer
	println("MakeSigner switch")
	switch {
	case config.IsEIP155(blockNumber):
		println("MakeSigner switch case 1")
		signer = NewEIP155Signer(config.ChainId)
	case config.IsHomestead(blockNumber):
		println("MakeSigner switch case 2")
		signer = HomesteadSigner{}
	default:
		println("MakeSigner switch case 3")
		signer = FrontierSigner{}
	}
	println("MakeSigner switch case 4")
	return signer
}

// SignTx signs the transaction using the given signer and private key
func SignTx(tx *Transaction, s Signer, prv *ecdsa.PrivateKey) (*Transaction, error) {
	println("SignTx h := s.Hash(tx)")
	h := s.Hash(tx)
	println("SignTx sig, err := crypto.Sign(h[:], prv)")
	sig, err := crypto.Sign(h[:], prv)
	println("SignTx if")
	if err != nil {
		return nil, err
	}
	println("SignTx return tx.WithSignature(s, sig)")
	return tx.WithSignature(s, sig)
}

// Sender returns the address derived from the signature (V, R, S) using secp256k1
// elliptic curve and an error if it failed deriving or upon an incorrect
// signature.
//
// Sender may cache the address, allowing it to be used regardless of
// signing method. The cache is invalidated if the cached signer does
// not match the signer used in the current call.
func Sender(signer Signer, tx *Transaction) (common.Address, error) {
	println("Sender if sc := tx.from.Load(); sc != nil")
	if sc := tx.from.Load(); sc != nil {
		println("Sender sigCache := sc.(sigCache)")
		sigCache := sc.(sigCache)
		// If the signer used to derive from in a previous
		// call is not the same as used current, invalidate
		// the cache.
		println("Sender if sigCache.signer.Equal(signer)")
		if sigCache.signer.Equal(signer) {
			return sigCache.from, nil
		}
	}
	println("Sender addr, err := signer.Sender(tx)")
	addr, err := signer.Sender(tx)

	println("Sender if err != nil")
	if err != nil {
		println("addr, err := signer.Sender(tx) trajja3 erreur")
		return common.Address{}, err
	}
	println("Sender tx.from.Store(sigCache{signer: signer, from: addr})")
	tx.from.Store(sigCache{signer: signer, from: addr})
	println("Sender return addr, nil")
	return addr, nil
}

// Signer encapsulates transaction signature handling. Note that this interface is not a
// stable API and may change at any time to accommodate new protocol rules.
type Signer interface {
	// Sender returns the sender address of the transaction.
	Sender(tx *Transaction) (common.Address, error)
	// SignatureValues returns the raw R, S, V values corresponding to the
	// given signature.
	SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error)
	// Hash returns the hash to be signed.
	Hash(tx *Transaction) common.Hash
	// Equal returns true if the given signer is the same as the receiver.
	Equal(Signer) bool
}

// EIP155Transaction implements Signer using the EIP155 rules.
type EIP155Signer struct {
	chainId, chainIdMul *big.Int
}

func NewEIP155Signer(chainId *big.Int) EIP155Signer {
	println("NewEIP155Signer if chainId == nil")
	if chainId == nil {
		chainId = new(big.Int)
	}
	println("NewEIP155Signer return EIP155Signer")
	return EIP155Signer{
		chainId:    chainId,
		chainIdMul: new(big.Int).Mul(chainId, big.NewInt(2)),
	}
}

func (s EIP155Signer) Equal(s2 Signer) bool {
	println("Equal eip155, ok := s2.(EIP155Signer)")
	eip155, ok := s2.(EIP155Signer)
	println("Equal return ok && eip155.chainId.Cmp(s.chainId) == 0")
	return ok && eip155.chainId.Cmp(s.chainId) == 0
}

var big8 = big.NewInt(8)

func (s EIP155Signer) Sender(tx *Transaction) (common.Address, error) {

	println("Sender if !tx.Protected()")
	if !tx.Protected() {
		return HomesteadSigner{}.Sender(tx)
	}
	println("Sender if tx.ChainId().Cmp(s.chainId) != 0")
	if tx.ChainId().Cmp(s.chainId) != 0 {
		println("sender return common.Address{}, ErrInvalidChainId")
		return common.Address{}, ErrInvalidChainId
	}
	println("Sender V := new(big.Int).Sub(tx.data.V, s.chainIdMul)")
	V := new(big.Int).Sub(tx.data.V, s.chainIdMul)
	println("Sender V.Sub(V, big8)")
	V.Sub(V, big8)
	println("Sender return recoverPlain(s.Hash(tx), tx.data.R, tx.data.S, V, true, tx)")
	return recoverPlain(s.Hash(tx), tx.data.R, tx.data.S, V, true, tx)
}

// WithSignature returns a new transaction with the given signature. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (s EIP155Signer) SignatureValues(tx *Transaction, sig []byte) (R, S, V *big.Int, err error) {
	println("SignatureValues R, S, V, err = HomesteadSigner{}.SignatureValues(tx, sig)")
	R, S, V, err = HomesteadSigner{}.SignatureValues(tx, sig)
	println("SignatureValues if err != nil")
	if err != nil {
		return nil, nil, nil, err
	}
	println("SignatureValues if s.chainId.Sign() != 0 ")
	if s.chainId.Sign() != 0 {
		println("SignatureValues V = big.NewInt(int64(sig[64] + 35))")
		V = big.NewInt(int64(sig[64] + 35))
		println("SignatureValues V.Add(V, s.chainIdMul)")
		V.Add(V, s.chainIdMul)
	}
	println("SignatureValues return R, S, V, nil")
	return R, S, V, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (s EIP155Signer) Hash(tx *Transaction) common.Hash {
	println("Hash return rlpHash([]interface{}")
	return rlpHash([]interface{}{
		tx.data.AccountNonce,
		tx.data.Price,
		tx.data.GasLimit,
		tx.data.Recipient,
		tx.data.Amount,
		tx.data.Payload,
		s.chainId, uint(0), uint(0),
	})
}

// HomesteadTransaction implements TransactionInterface using the
// homestead rules.
type HomesteadSigner struct{ FrontierSigner }

func (s HomesteadSigner) Equal(s2 Signer) bool {
	println("Equal _, ok := s2.(HomesteadSigner)")
	_, ok := s2.(HomesteadSigner)
	println("equal return ok")
	return ok
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (hs HomesteadSigner) SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error) {
	println("SignatureValues return hs.FrontierSigner.SignatureValues(tx, sig)")
	return hs.FrontierSigner.SignatureValues(tx, sig)
}

func (hs HomesteadSigner) Sender(tx *Transaction) (common.Address, error) {
	println("Sender return recoverPlain(hs.Hash(tx), tx.data.R, tx.data.S, tx.data.V, true,tx)")
	return recoverPlain(hs.Hash(tx), tx.data.R, tx.data.S, tx.data.V, true,tx)
}

type FrontierSigner struct{}

func (s FrontierSigner) Equal(s2 Signer) bool {
	println("Equal _, ok := s2.(FrontierSigner)")
	_, ok := s2.(FrontierSigner)
	println("Equal return ok")
	return ok
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (fs FrontierSigner) SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error) {
	println("Equal return ok")
	if len(sig) != 65 {
		panic(fmt.Sprintf("wrong size for signature: got %d, want 65", len(sig)))
	}
	println("SignatureValues r = new(big.Int).SetBytes(sig[:32])")
	r = new(big.Int).SetBytes(sig[:32])
	println("SignatureValues s = new(big.Int).SetBytes(sig[32:64])")
	s = new(big.Int).SetBytes(sig[32:64])
	println("SignatureValues v = new(big.Int).SetBytes([]byte{sig[64] + 27})")
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	println("SignatureValues return r, s, v, nil")
	return r, s, v, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (fs FrontierSigner) Hash(tx *Transaction) common.Hash {
	println("Hash return rlpHash([]interface{}")
	return rlpHash([]interface{}{
		tx.data.AccountNonce,
		tx.data.Price,
		tx.data.GasLimit,
		tx.data.Recipient,
		tx.data.Amount,
		tx.data.Payload,
	})
}

func (fs FrontierSigner) Sender(tx *Transaction) (common.Address, error) {
		println("baypass return recoverPlain(fs.Hash(tx), tx.data.R, tx.data.S, tx.data.V, false, tx)")
	return recoverPlain(fs.Hash(tx), tx.data.R, tx.data.S, tx.data.V, false, tx)
}

func recoverPlain(sighash common.Hash, R, S, Vb *big.Int, homestead bool, tx *Transaction) (common.Address, error) {


	if Vb.BitLen() > 8 {
		println("Vb.BitLen() > 8")
		return common.Address{}, ErrInvalidSig
	}
	V := byte(Vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(V, R, S, homestead) {
		println("invalid transaction")
		return common.Address{}, ErrInvalidSig
	}
	// encode the snature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, 65)
	println("step copy(sig[32-len(r):32], r)")
	copy(sig[32-len(r):32], r)
	println("step copy(sig[64-len(s):64], s)")
	copy(sig[64-len(s):64], s)
	println("step sig[64] = V")
	sig[64] = V
	// recover the public key from the snature
	println("step 	pub, err := crypto.Ecrecover(sighash[:], sig)")
	var sigHash common.Hash = tx.Hash()
	pub, err := crypto.Ecrecover(sigHash[:], sig)
	if err != nil {
		println("error in condition : crypto.Ecrecover(sighash[:], sig)")
		return common.Address{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return common.Address{}, errors.New("invalid public key")
	}
	println("step 	var addr common.Address")
	var addr common.Address
	println("step 	copy(addr[:], crypto.Keccak256(pub[1:])[12:])")
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	println("step 	return addr, nil")
	return addr, nil
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	println("deriveChainId 	if v.BitLen() <= 64")
	if v.BitLen() <= 64 {
		println("deriveChainId 	v := v.Uint64()")
		v := v.Uint64()
		println("deriveChainId 	if v == 27 || v == 28 ")
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		println("deriveChainId return new(big.Int).SetUint64((v - 35) / 2)")
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	println("deriveChainId 	v = new(big.Int).Sub(v, big.NewInt(35))")
	v = new(big.Int).Sub(v, big.NewInt(35))
	println("deriveChainId 	return v.Div(v, big.NewInt(2))")
	return v.Div(v, big.NewInt(2))
}
