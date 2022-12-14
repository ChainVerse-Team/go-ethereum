package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type BlockRewardRole uint64

const (
	Validator BlockRewardRole = iota
	Covenant
	Unknown
)

type RPCBlockReward struct {
	Address          common.Address `json:"address"`
	Role             uint64         `json:"role"`
	Epoch            uint64         `json:"epoch"`
	Number           uint64         `json:"number"`
	Amount           *big.Int       `json:"amount"`
	TotalFromGenesis *big.Int       `json:"totalFromGenesis"`
}

type RPCBlockRewards []RPCBlockReward

type BlockReward struct {
	Address          common.Address
	Role             BlockRewardRole
	Epoch            uint64
	Number           uint64
	Amount           *big.Int
	TotalFromGenesis *big.Int
}

type BlockRewards struct {
	blockRewards []BlockReward
}

func (br *BlockReward) SetRole(r uint64) {
	role := BlockRewardRole(r)
	if role > Unknown || role < Validator {
		br.Role = Unknown
		return
	}

	br.Role = role
}

func (r *RPCBlockReward) ToBlockReward() *BlockReward {
	res := &BlockReward{
		Address:          r.Address,
		Epoch:            r.Epoch,
		Number:           r.Number,
		Amount:           r.Amount,
		TotalFromGenesis: r.TotalFromGenesis,
	}
	res.SetRole(r.Role)

	return res
}

func (r *RPCBlockRewards) ToBlockRewards() *BlockRewards {
	res := &BlockRewards{
		blockRewards: make([]BlockReward, len(*r)),
	}
	for i := 0; i < len(*r); i++ {
		tmp := *r
		res.blockRewards[i] = *tmp[i].ToBlockReward()
	}
	return res
}
