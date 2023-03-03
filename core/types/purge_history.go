package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"strings"
)

type PurgeHistoryResultJSON struct {
	DemotedAscendance   []hexutil.Uint64 `json:"demoted_ascendance"`
	PromotedPaladin     []hexutil.Uint64 `json:"promoted_paladin"`
	DemotedPaladin      []hexutil.Uint64 `json:"demoted_paladin"`
	PromotedTemplar     []hexutil.Uint64 `json:"promoted_templar"`
	DemotedTemplar      []hexutil.Uint64 `json:"demoted_templar"`
	PromotedCavalier    []hexutil.Uint64 `json:"promoted_cavalier"`
	DemotedCavalier     []hexutil.Uint64 `json:"demoted_cavalier"`
	PromotedLegionnaire []hexutil.Uint64 `json:"promoted_legionnaire"`
}

type PurgeHistoryResult struct {
	DemotedAscendance   []uint64
	PromotedPaladin     []uint64
	DemotedPaladin      []uint64
	PromotedTemplar     []uint64
	DemotedTemplar      []uint64
	PromotedCavalier    []uint64
	DemotedCavalier     []uint64
	PromotedLegionnaire []uint64
}

func (p *PurgeHistoryResultJSON) ToPurgeHistoryResult() *PurgeHistoryResult {
	return &PurgeHistoryResult{
		DemotedAscendance:   convertUtilUint64ArrayToUint64Array(p.DemotedAscendance),
		PromotedPaladin:     convertUtilUint64ArrayToUint64Array(p.PromotedPaladin),
		DemotedPaladin:      convertUtilUint64ArrayToUint64Array(p.DemotedPaladin),
		PromotedTemplar:     convertUtilUint64ArrayToUint64Array(p.PromotedTemplar),
		DemotedTemplar:      convertUtilUint64ArrayToUint64Array(p.DemotedTemplar),
		PromotedCavalier:    convertUtilUint64ArrayToUint64Array(p.PromotedCavalier),
		DemotedCavalier:     convertUtilUint64ArrayToUint64Array(p.DemotedCavalier),
		PromotedLegionnaire: convertUtilUint64ArrayToUint64Array(p.PromotedLegionnaire),
	}
}

func convertUtilUint64ArrayToUint64Array(arr []hexutil.Uint64) []uint64 {
	rs := make([]uint64, 0)
	for _, ele := range arr {
		rs = append(rs, uint64(ele))
	}
	return rs
}

func (p *PurgeHistoryResult) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Demoted Ascendance:\t %v\n", p.DemotedAscendance))
	sb.WriteString(fmt.Sprintf("Promoted Paladin:\t %v\n", p.PromotedPaladin))
	sb.WriteString(fmt.Sprintf("Demoted Paladin:\t %v\n", p.DemotedPaladin))
	sb.WriteString(fmt.Sprintf("Promoted Templar:\t %v\n", p.PromotedTemplar))
	sb.WriteString(fmt.Sprintf("Demoted Templar:\t %v\n", p.DemotedTemplar))
	sb.WriteString(fmt.Sprintf("Promoted Cavalier:\t %v\n", p.PromotedCavalier))
	sb.WriteString(fmt.Sprintf("Demoted Cavalier:\t %v\n", p.DemotedCavalier))
	sb.WriteString(fmt.Sprintf("Promoted Legionnaire:\t %v\n", p.PromotedLegionnaire))
	sb.WriteString(fmt.Sprintln("Demoted: List of Token IDs sorted by lowest to highest SCREE won"))
	sb.WriteString(fmt.Sprintln("Promoted: List of Token IDs sorted by highest to lowest SCREE won"))
	return sb.String()
}
