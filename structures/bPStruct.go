package structures

type Factor struct {
	ID        string
	Value     float64
	Parameter *float64
	Disabled  bool
}

type BasePoint struct {
	ID            string
	EventKind     string
	ProviderLayer string
	TimeStamp     int64
	ScoreC1       int
	ScoreC2       int
	Factors       []Factor
}
