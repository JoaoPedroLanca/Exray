package classifier

import (
	"testing"

	"github.com/JoaoPedroLanca/Exray/agent/internal/emitter"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected emitter.EventType
	}{
		// ── Erros ───────────────────────────────────────────────────────────
		{"error lowercase", "error: undefined variable x", emitter.EventTypeError},
		{"error uppercase", "ERROR: undefined variable x", emitter.EventTypeError},
		{"panic runtime", "panic: runtime error: index out of range [0]", emitter.EventTypeError},
		{"fatal", "fatal: repository not found", emitter.EventTypeError},
		{"exception java", "java.lang.NullPointerException at Main.java:42", emitter.EventTypeError},
		{"err colon", "err: connection refused", emitter.EventTypeError},
		{"failed build", "Build FAILED", emitter.EventTypeError},
		{"go test fail", "--- FAIL: TestSomething (0.12s)", emitter.EventTypeError},
		{"critical alert", "CRITICAL: disk space low", emitter.EventTypeError},
		{"dotnet error", "MSB3021: Unable to copy file", emitter.EventTypeError},

		// ── Warnings ────────────────────────────────────────────────────────
		{"warning lowercase", "warning: deprecated function foo", emitter.EventTypeWarning},
		{"warn prefix", "warn: config file not found, using defaults", emitter.EventTypeWarning},
		{"deprecated", "deprecated: use NewFunction instead of OldFunction", emitter.EventTypeWarning},

		// ── Logs normais ─────────────────────────────────────────────────────
		{"server started", "server started on port 8080", emitter.EventTypeLog},
		{"build success", "Build succeeded", emitter.EventTypeLog},
		{"go test ok", "ok  \tgithub.com/user/repo\t0.123s", emitter.EventTypeLog},
		{"info message", "info: application started", emitter.EventTypeLog},

		// ── Falsos positivos que devem ser Log ───────────────────────────────
		// "scaffold" não contém " fail" (sem espaço antes de "fail")
		{"scaffold no fp", "scaffold created successfully", emitter.EventTypeLog},
		// "workflow" não contém nenhum padrão
		{"workflow no fp", "workflow completed in 2.3s", emitter.EventTypeLog},
		// "default" não contém " fail"
		{"default no fp", "using default configuration", emitter.EventTypeLog},
		// "profile" não contém " fail"
		{"profile no fp", "profile loaded from ~/.bashrc", emitter.EventTypeLog},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Classify(tt.line)
			if got != tt.expected {
				t.Errorf("Classify(%q)\n  got:  %v\n  want: %v", tt.line, got, tt.expected)
			}
		})
	}
}
