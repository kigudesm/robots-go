package structures

// Период в матче
type PartStruct struct {
	IsGoing bool // Идет ли игра
	Nmb     int  // Номер периода (Если игра идет: 3 или 4 дополнительное время, 5 серия пенальти;
	// Если игра не идет: 0 матч не начался,  i перерыв после i-ого тайма, -1 матч закончился)
}

// Текущее состояние матча
type MatchStateCurrent struct {
	TimeStampCurrent int64             // Текущее время в формате timestamp
	Events           []EventStruct     // Трансляция на текущий момент
	Part             PartStruct        // Тайм в матче
	Timer            int64             // Значение таймера
	ActiveEventKinds []string          // eventKinds которые требуют вычислений
	Penalty          int               // 0 - если пенальти не пробивается, i - команда пробивающая пенальти
	Removal          [2]int            // количество удалений в обеих командах
	Scores           map[string][2]int // счет по каждому activeEventKind
	History
}
