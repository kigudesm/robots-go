package gametypes

type MatchType struct {
	TFirstHalf   int
	TSecondHalf  int
	HalfDuration int
}

var MatchTypes = map[string]MatchType{
	"Football.Common": {
		TFirstHalf:   2805,
		TSecondHalf:  2960,
		HalfDuration: 2700,
	},
	"Football.Friendly": {
		TFirstHalf:   2805,
		TSecondHalf:  2960,
		HalfDuration: 2700,
	},
	"Football.TwoForty": {
		TFirstHalf:   2490,
		TSecondHalf:  2610,
		HalfDuration: 2400,
	},
}
