package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	Role             hexutil.Uint64 `json:"role"`
	Epoch            hexutil.Uint64 `json:"epoch"`
	Number           hexutil.Uint64 `json:"number"`
	Amount           hexutil.Big    `json:"amount"`
	TotalFromGenesis hexutil.Big    `json:"totalFromGenesis"`
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

func (br *BlockReward) GetRole() string {
	switch br.Role {
	case Validator:
		return "validator"
	case Covenant:
		return "covenant"
	default:
		return ""
	}
}

func (r *RPCBlockReward) ToBlockReward() *BlockReward {
	res := &BlockReward{
		Address:          r.Address,
		Epoch:            uint64(r.Epoch),
		Number:           uint64(r.Number),
		Amount:           (*big.Int)(&r.Amount),
		TotalFromGenesis: (*big.Int)(&r.TotalFromGenesis),
	}
	res.SetRole(uint64(r.Role))

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

func (b *BlockRewards) Len() int {
	return len(b.blockRewards)
}

func (b *BlockRewards) Array() []BlockReward {
	return b.blockRewards
}
