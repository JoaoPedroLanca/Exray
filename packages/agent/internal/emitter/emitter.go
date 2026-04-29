package emitter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/oklog/ulid/v2"
)

type Emitter struct {
	mu  sync.Mutex
	w   *bufio.Writer
	seq atomic.Int64
}

func NewEmitter() *Emitter {
	return &Emitter{
		w: bufio.NewWriter(os.Stdout),
	}
}

func (e *Emitter) Emit(eventType EventType, stage StageType, payload map[string]interface{}) error {
	id := ulid.Make().String()
	seq := e.seq.Add(1)

	event := Event{
		Id:        id,
		Type:      eventType,
		Stage:     stage,
		Timestamp: time.Now().UnixMilli(),
		Seq:       seq,
		Payload:   payload,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("emitter: Falha ao serializar evento: %w", err)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if _, err := e.w.Write(data); err != nil {
		return fmt.Errorf("emitter: Falha ao escrever evento: %w", err)
	}
	if err := e.w.WriteByte('\n'); err != nil {
		return fmt.Errorf("emitter: Falha ao escrever nova linha: %w", err)
	}
	if err := e.w.Flush(); err != nil {
		return fmt.Errorf("emitter: Falha ao flushar buffer: %w", err)
	}
	return nil
}
