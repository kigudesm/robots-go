package structures

// Трансляция
type EventInfo struct {
	ID      int64  // id события
	RegTime string // время
	Type    int    // тип
	I1      *int
	I2      *int
	I3      *int
	I4      *int
	I5      *int
}

// Описывает конфигурацию блокировки
type BlockConfig struct {
	Description string
	Block       map[string]bool // Множество блокируемых событий от сканера
	Cancel      map[int]bool    // Множество отменяющих событий трансляции
}
