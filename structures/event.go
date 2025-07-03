package structures

// Трансляция
type EventStruct struct {
	ID      int64  // id события
	RegTime string // время
	Type    int    // тип
	// Параметры события (могут отсутствовать)
	I1 *int
	I2 *int
	I3 *int
	I4 *int
	I5 *int
}
