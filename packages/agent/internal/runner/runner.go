package runner

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/JoaoPedroLanca/Exray/agent/internal/classifier"
	"github.com/JoaoPedroLanca/Exray/agent/internal/emitter"
)

type Runner struct {
	emitter *emitter.Emitter
	cmds    map[int]*exec.Cmd
	mu      sync.Mutex
}

func NewRunner(e *emitter.Emitter) *Runner {
	return &Runner{
		emitter: e,
		cmds:    make(map[int]*exec.Cmd),
	}
}

func (r *Runner) Start(ctx context.Context, bin string, args []string, stage emitter.StageType) error {
	cmd := exec.CommandContext(ctx, bin, args...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("runner: falha ao criar stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("runner: falha ao criar stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("runner: falha ao iniciar processo %q: %w", bin, err)
	}

	pid := cmd.Process.Pid

	r.mu.Lock()
	r.cmds[pid] = cmd
	r.mu.Unlock()

	startTime := time.Now()

	_ = r.emitter.Emit(emitter.EventTypeProcessStart, stage, map[string]interface{}{
		"bin":  bin,
		"args": args,
		"pid":  pid,
	})

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		r.streamLines(stdoutPipe, stage, "stdout")
	}()

	go func() {
		defer wg.Done()
		r.streamLines(stderrPipe, stage, "stderr")
	}()

	wg.Wait()

	exitCode := 0
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	durationMs := time.Since(startTime).Milliseconds()

	// Remove o processo do mapa após encerrar
	r.mu.Lock()
	delete(r.cmds, pid)
	r.mu.Unlock()

	_ = r.emitter.Emit(emitter.EventTypeProcessExit, stage, map[string]interface{}{
		"exitCode":   exitCode,
		"durationMs": durationMs,
	})

	return nil
}

func (r *Runner) streamLines(reader io.Reader, stage emitter.StageType, source string) {
	scanner := bufio.NewScanner(reader)

	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		eventType := classifier.Classify(line)
		_ = r.emitter.Emit(eventType, stage, map[string]interface{}{
			"message": line,
			"source":  source,
		})
	}
}

func (r *Runner) Stop() error {
	r.mu.Lock()
	// Copia o mapa para não segurar o lock durante os signals
	toStop := make(map[int]*exec.Cmd, len(r.cmds))
	for pid, cmd := range r.cmds {
		toStop[pid] = cmd
	}
	r.mu.Unlock()

	for _, cmd := range toStop {
		if cmd.Process == nil {
			continue
		}
		if runtime.GOOS == "windows" {
			_ = cmd.Process.Signal(os.Interrupt)
		} else {
			_ = cmd.Process.Signal(syscall.SIGTERM)
		}
	}

	if len(toStop) == 0 {
		return nil
	}

	go func() {
		timer := time.NewTimer(3 * time.Second)
		defer timer.Stop()
		<-timer.C

		r.mu.Lock()
		survivors := make([]*exec.Cmd, 0, len(r.cmds))
		for _, cmd := range r.cmds {
			survivors = append(survivors, cmd)
		}
		r.mu.Unlock()

		for _, cmd := range survivors {
			if cmd.Process != nil {
				_ = cmd.Process.Kill()
			}
		}
	}()

	return nil
}
