package emitter

// Evento do agent
type EventType string

const (
	EventTypeLog          EventType = "log"
	EventTypeError        EventType = "error"
	EventTypeWarning      EventType = "warning"
	EventTypeProcessStart EventType = "process_start"
	EventTypeProcessExit  EventType = "process_exit"
	EventTypeStageChange  EventType = "stage_change"
)

// Estágio do pipeline
type StageType string

const (
	StageTypeBuild   StageType = "build"
	StageTypeRun     StageType = "run"
	StageTypeTest    StageType = "test"
	StageTypeUnknown StageType = "unknown"
)

type Event struct {
	Id        string                 `json:"id"`
	Type      EventType              `json:"type"`
	Stage     StageType              `json:"stage"`
	Timestamp int64                  `json:"timestamp"`
	Seq       int64                  `json:"seq"`
	Payload   map[string]interface{} `json:"payload"`
}

type AgentCommand struct {
	Command string    `json:"command"`
	Bin     string    `json:"bin"`
	Args    []string  `json:"args"`
	Stage   StageType `json:"stage"`
	Timeout int       `json:"timeout"`
}
