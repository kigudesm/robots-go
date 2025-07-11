package structures

// Период в матче
type Part struct {
	IsGoing bool // Идет ли игра
	Nmb     int  // Номер периода (Если игра идет: 3 или 4 дополнительное время, 5 серия пенальти;
	// Если игра не идет: 0 матч не начался,  i перерыв после i-ого тайма, -1 матч закончился)
}

type EventKindsActive struct {
	// BPID      string // Id опорной точки по которой будем считать
	IsLive    bool // Нужно показывать на сайте
	IsActive  bool // Нужно вычислять парамтеры оптимизации
	IsBlocked bool // В блоке
}

// Текущее состояние матча
type MatchStateCurrent struct {
	Timestamp int64 // Текущее время в формате timestamp
	// Events     []EventInfo            // Трансляция на текущий момент
	Half       Part                        // Тайм в матче
	Timer      int64                       // Значение таймера
	InjuryTime [2]int                      // Компенсированное время фактическое
	EventKinds map[string]EventKindsActive // eventKinds которые требуют действий
	Suspended  bool                        // Матч приостановлен
	Var        bool                        // Просмотр видеорефери
	Penalty    int                         // 0 - если пенальти не пробивается, i - команда пробивающая пенальти
	Removal    [2]int                      // Количество удалений в каждой из команд
	History    map[string][]EventInfo      // История событий за период игры (в обратном порядке) для каждого eventKind
	ScoreHost  map[string]int              // Счет (количество событий за период игры) команды хозяев для каждого eventKind
	ScoreAway  map[string]int              // Счет (количество событий за период игры) команды гостей для каждого eventKind
	Alerts     *[]string                   // Алерты
}
