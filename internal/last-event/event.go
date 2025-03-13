package lastevent

import (
	"encoding/binary"
	"os"
	"sync/atomic"
	"time"
)

const lastEventFile = "last_event.dat"

var LastEvent atomic.Int64

func InitLastEvent() {
	data, err := os.ReadFile(lastEventFile)
	if err == nil && len(data) == 8 {
		value := int64(binary.LittleEndian.Uint64(data))
		LastEvent.Store(value)
	}

	go periodicSave()
}

func SaveLastEvent() error {
	value := LastEvent.Load()
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, uint64(value))
	return os.WriteFile(lastEventFile, data, 0644)
}

func periodicSave() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		SaveLastEvent()
	}
}
