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

type RewardRecord struct {
	RunningTotal               *big.Int   // running total since genesis
	Daily                      []*big.Int // daily history for last 28 days
	Monthly                    []*big.Int // monthly earnings for 12 months
	AmountEarnedToday          *big.Int   // total amount of tokens earned today
	AmountEarnedSinceLastPurge *big.Int   // total amount of tokens earned since last purge
	AmountEarnedThisMoon       *big.Int   // total amount of tokens earned in this moon
	moonCalendarCounter        int        // is used to track the current date on a 28-day calendar (starts at 0)
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
type rpcRewardRecord struct {
	RunningTotal               hexutil.Big    `json:"runningTotal"`
	Daily                      []hexutil.Big  `json:"daily"`
	Monthly                    []hexutil.Big  `json:"monthly"`
	AmountEarnedToday          hexutil.Big    `json:"amountEarnedToday"`
	AmountEarnedSinceLastPurge hexutil.Big    `json:"amountEarnedSinceLastPurge"`
	AmountEarnedThisMoon       hexutil.Big    `json:"amountEarnedThisMoon"`
	MoonCalendarCounter        hexutil.Uint64 `json:"moonCalendarCounter"`
}

type RPCValidatorRewardRecord struct {
	rpcRewardRecord `json:"rewardRecord"`
	Address         common.Address `json:"address"`
}

type RPCCovenantNFTRewardRecord struct {
	rpcRewardRecord `json:"rewardRecord"`
	TokenID         hexutil.Uint64 `json:"tokenID"`
}

func newRewardRecord() *RewardRecord {
	return &RewardRecord{
		RunningTotal:               big.NewInt(0),
		Daily:                      make([]*big.Int, 0),
		Monthly:                    make([]*big.Int, 0),
		AmountEarnedToday:          big.NewInt(0),
		AmountEarnedSinceLastPurge: big.NewInt(0),
		AmountEarnedThisMoon:       big.NewInt(0),
		moonCalendarCounter:        0,
	}
}

func toRewardRecord(rpcRc rpcRewardRecord) *RewardRecord {
	rs := newRewardRecord()
	rs.RunningTotal = rpcRc.RunningTotal.ToInt()
	for _, d := range rpcRc.Daily {
		rs.Daily = append(rs.Daily, d.ToInt())
	}
	for _, m := range rpcRc.Monthly {
		rs.Monthly = append(rs.Monthly, m.ToInt())
	}
	rs.AmountEarnedToday = rpcRc.AmountEarnedToday.ToInt()
	rs.AmountEarnedSinceLastPurge = rpcRc.AmountEarnedSinceLastPurge.ToInt()
	rs.AmountEarnedThisMoon = rpcRc.AmountEarnedThisMoon.ToInt()
	rs.moonCalendarCounter = int(rpcRc.MoonCalendarCounter)

	return rs
}

func (v *RPCValidatorRewardRecord) ToValidatorRewardRecord() *ValidatorRewardRecord {
	rs := &ValidatorRewardRecord{
		RewardRecord: toRewardRecord(v.rpcRewardRecord),
		Address:      v.Address,
	}

	return rs
}

func (c *RPCCovenantNFTRewardRecord) ToCovenantNFTRewardRecord() *CovenantNFTRewardRecord {
	rs := &CovenantNFTRewardRecord{
		RewardRecord: toRewardRecord(c.rpcRewardRecord),
		TokenID:      uint64(c.TokenID),
	}

	return rs
}

// GetTodayIndex returns how many days have passed on 28-day calendar
func (r *RewardRecord) GetTodayIndex() int {
	if r.moonCalendarCounter < 0 || r.moonCalendarCounter > MaxDateInMonth {
		return 0
	}
	return r.moonCalendarCounter
}

func (r *RewardRecord) GetTotalTokensEarnedToday() *big.Int {
	return r.AmountEarnedToday
}

func (r *RewardRecord) GetTotalTokensEarnedLastMonth() *big.Int {
	if r == nil || r.Monthly == nil {
		return nil
	}
	if len(r.Monthly) == 0 {
		return nil
	}
	return r.Monthly[len(r.Monthly)-1]
}

// GetTotalTokensEarnedSince returns total amount of tokens earned from index
func (r *RewardRecord) GetTotalTokensEarnedSince(i int) *big.Int {
	today := r.GetTotalTokensEarnedToday()
	if r == nil || r.Daily == nil || today == nil || i < 0 {
		return nil
	}

	totalSinceIndex := big.NewInt(0)
	for j := i; j < len(r.Daily); j++ {
		if r.Daily[j] == nil {
			return nil
		}
		totalSinceIndex = totalSinceIndex.Add(totalSinceIndex, r.Daily[j])
	}
	totalSinceIndex = totalSinceIndex.Add(totalSinceIndex, today)

	return totalSinceIndex
}

func (r *RewardRecord) GetTotalTokensEarnedSinceLastMonth() *big.Int {
	return r.AmountEarnedThisMoon
}

func (r *RewardRecord) GetTotalTokensEarnedSinceLastPurge() *big.Int {
	return r.AmountEarnedSinceLastPurge
}
