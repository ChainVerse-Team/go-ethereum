package types

import "github.com/ethereum/go-ethereum/common/hexutil"

type PurgeHistoryResult struct {
	DemotedAscendance   []hexutil.Uint64 `json:"demoted_ascendance"`
	PromotedPaladin     []hexutil.Uint64 `json:"promoted_paladin"`
	DemotedPaladin      []hexutil.Uint64 `json:"demoted_paladin"`
	PromotedTemplar     []hexutil.Uint64 `json:"promoted_templar"`
	DemotedTemplar      []hexutil.Uint64 `json:"demoted_templar"`
	PromotedCavalier    []hexutil.Uint64 `json:"promoted_cavalier"`
	DemotedCavalier     []hexutil.Uint64 `json:"demoted_cavalier"`
	PromotedLegionnaire []hexutil.Uint64 `json:"promoted_legionnaire"`
}
