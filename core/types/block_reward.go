package types

import (
	"fmt"
	"math/big"
	"strings"

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
	ValidatorRecords    []ValidatorRewardRecord
	CovenantRecords     []CovenantNFTRewardRecord
	moonCalendarCounter int
}

type RewardRecord struct {
	RunningTotal               *big.Int   // running total since genesis
	Daily                      []*big.Int // daily history for last 28 days
	Monthly                    []*big.Int // monthly earnings for 12 months
	AmountEarnedToday          *big.Int   // total amount of tokens earned today
	AmountEarnedSinceLastPurge *big.Int   // total amount of tokens earned since last purge
	AmountEarnedThisMoon       *big.Int   // total amount of tokens earned in this moon
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
	RunningTotal               *hexutil.Big   `json:"runningTotal"`
	Daily                      []*hexutil.Big `json:"daily"`
	Monthly                    []*hexutil.Big `json:"monthly"`
	AmountEarnedToday          *hexutil.Big   `json:"amountEarnedToday"`
	AmountEarnedSinceLastPurge *hexutil.Big   `json:"amountEarnedSinceLastPurge"`
	AmountEarnedThisMoon       *hexutil.Big   `json:"amountEarnedThisMoon"`
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
	}
}

func toRewardRecord(rpcRc rpcRewardRecord) *RewardRecord {
	rs := newRewardRecord()
	rs.RunningTotal = (*big.Int)(rpcRc.RunningTotal)
	for _, d := range rpcRc.Daily {
		rs.Daily = append(rs.Daily, (*big.Int)(d))
	}
	for _, m := range rpcRc.Monthly {
		rs.Monthly = append(rs.Monthly, (*big.Int)(m))
	}
	rs.AmountEarnedToday = (*big.Int)(rpcRc.AmountEarnedToday)
	rs.AmountEarnedSinceLastPurge = (*big.Int)(rpcRc.AmountEarnedSinceLastPurge)
	rs.AmountEarnedThisMoon = (*big.Int)(rpcRc.AmountEarnedThisMoon)

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

func (s *NodeRewardStorage) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("MoonCounter: %d\n", s.moonCalendarCounter))
	sb.WriteString(fmt.Sprintf("ValidatorRecords: %v\n", s.ValidatorRecords))
	sb.WriteString(fmt.Sprintf("CovenantRecords: %v\n", s.CovenantRecords))

	return sb.String()
}

func (r *RewardRecord) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[Daily: %v, ", r.Daily))
	sb.WriteString(fmt.Sprintf("Monthly: %v, ", r.Monthly))
	sb.WriteString(fmt.Sprintf("RunningTotal: %v, ", r.RunningTotal))
	sb.WriteString(fmt.Sprintf("AmountEarnedThisMoon: %v, ", r.AmountEarnedThisMoon))
	sb.WriteString(fmt.Sprintf("AmountEarnedSinceLastPurge: %v, ", r.AmountEarnedSinceLastPurge))
	sb.WriteString(fmt.Sprintf("AmountEarnedToday: %v]", r.AmountEarnedToday))
	return sb.String()
}

func (v *ValidatorRewardRecord) String() string {
	return fmt.Sprintf("[Address: %v, record: %s]", v.Address, v.RewardRecord.String())
}

func (c *CovenantNFTRewardRecord) String() string {
	return fmt.Sprintf("[TokenID: %v, record: %s]", c.TokenID, c.RewardRecord.String())
}
