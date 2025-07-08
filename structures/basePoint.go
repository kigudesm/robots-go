package structures

type Factor struct {
	FactorID string
	Value    string
	Disabled bool
}

type BasePointStruct struct {
	ID            string
	EventKind     string
	ProviderLayer string
	TimeStamp     int64
	ScoreC1       int
	ScoreC2       int
	Factors       []Factor
}
