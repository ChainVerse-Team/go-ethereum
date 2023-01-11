package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	MaxDateInMonth = 28
	MaxMonthInYear = 12
	PurgePeriod    = 14
)

type BlockRewardInsertType uint8

type NodeRewardStorage struct {
	ValidatorRecords []ValidatorRewardRecord
	CovenantRecords  []CovenantNFTRewardRecord
}

// RewardRecord new day or month will be pushed at the end
type RewardRecord struct {
	Daily   []*big.Int // starts from 0..27
	Monthly []*big.Int // starts from 0..11. Reset Daily when new month begins
}

type ValidatorRewardRecord struct {
	*RewardRecord
	Address common.Address
}

type CovenantNFTRewardRecord struct {
	*RewardRecord
	TokenID uint64
}

// rpc type for RewardRecord
type rewardRecord struct {
	Daily   []hexutil.Big `json:"daily"`
	Monthly []hexutil.Big `json:"monthly"`
}

type RPCValidatorRewardRecord struct {
	rewardRecord `json:"rewardRecord"`
	Address      common.Address `json:"address"`
}

type RPCCovenantNFTRewardRecord struct {
	rewardRecord `json:"rewardRecord"`
	TokenID      hexutil.Uint64 `json:"tokenID"`
}

func (v *RPCValidatorRewardRecord) ToValidatorRewardRecord() *ValidatorRewardRecord {
	rs := &ValidatorRewardRecord{
		RewardRecord: &RewardRecord{
			Daily:   make([]*big.Int, len(v.Daily)),
			Monthly: make([]*big.Int, len(v.Monthly)),
		},
		Address: v.Address,
	}
	for i, d := range v.Daily {
		rs.Daily[i] = d.ToInt()
	}
	for i, m := range v.Monthly {
		rs.Monthly[i] = m.ToInt()
	}

	return rs
}

func (c *RPCCovenantNFTRewardRecord) ToCovenantNFTRewardRecord() *CovenantNFTRewardRecord {
	rs := &CovenantNFTRewardRecord{
		RewardRecord: &RewardRecord{
			Daily:   make([]*big.Int, len(c.Daily)),
			Monthly: make([]*big.Int, len(c.Monthly)),
		},
		TokenID: uint64(c.TokenID),
	}
	for i, d := range c.Daily {
		rs.Daily[i] = d.ToInt()
	}
	for i, m := range c.Monthly {
		rs.Monthly[i] = m.ToInt()
	}

	return rs
}

func (r *RewardRecord) GetRunningTotalToday() *big.Int {
	if len(r.Daily) == 0 {
		return nil
	}
	return r.Daily[len(r.Daily)-1]
}

func (r *RewardRecord) GetRunningLastMonth() *big.Int {
	if len(r.Monthly) == 0 {
		return nil
	}
	return r.Monthly[len(r.Monthly)-1]
}

func (r *RewardRecord) GetRunningTotalAt(i int) *big.Int {
	_size := len(r.Daily)
	if i < 0 || i >= _size {
		return nil
	}
	return r.Daily[i]
}

func (r *RewardRecord) GetRunningTotalSinceLastMonth() *big.Int {
	today := r.GetRunningTotalToday()
	if today == nil {
		return nil
	}

	lastMth := r.GetRunningLastMonth()
	if lastMth == nil {
		return today
	}

	tmp := big.NewInt(0).Set(today)
	return tmp.Sub(tmp, lastMth)
}

func (r *RewardRecord) GetRunningTotalSinceLastPurge() *big.Int {
	today := r.GetRunningTotalToday()
	if today == nil {
		return nil
	}

	tmp := big.NewInt(0).Set(today)
	todayInd := len(r.Daily) - 1
	if todayInd < PurgePeriod {
		// last purge index could be 0 if there's no month record. Hence, the current running total today
		// in case there's monthly records. LastPurge = Today - LastMonth
		return r.GetRunningTotalSinceLastMonth()
	}
	purge := r.GetRunningTotalAt(PurgePeriod - 1)
	if purge == nil {
		return nil
	}
	return tmp.Sub(tmp, purge)
}

func (r *RewardRecord) GetRunningTotalOverLastWeek() *big.Int {
	if r == nil || r.Daily == nil {
		return nil
	}
	_size := len(r.Daily)
	weekIndex := _size/7*7 - 1
	today := r.GetRunningTotalToday()
	if today == nil {
		return nil
	}
	if weekIndex == -1 {
		return r.GetRunningTotalSinceLastMonth()
	}
	tmp := big.NewInt(0).Set(today)
	lastWeekTotal := r.GetRunningTotalAt(weekIndex)
	if lastWeekTotal == nil {
		return nil
	}
	return tmp.Sub(tmp, lastWeekTotal)
}
