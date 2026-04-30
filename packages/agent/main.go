package agent

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JoaoPedroLanca/Exray/agent/internal/emitter"
	"github.com/JoaoPedroLanca/Exray/agent/internal/runner"
	"github.com/JoaoPedroLanca/Exray/agent/internal/sanitizer"
)

func main() {
	emit := emitter.NewEmitter()
	run := runner.NewRunner(emit)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		_ = run.Stop()
		os.Exit(0)
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var cmd emitter.AgentCommand
		if err := json.Unmarshal([]byte(line), &cmd); err != nil {
			_ = emit.Emit(emitter.EventTypeError, emitter.StageTypeUnknown, map[string]interface{}{
				"message": fmt.Sprintf("falha ao parsear comando stdin: %v", err),
				"source":  "agent",
			})
			continue
		}

		switch cmd.Command {
		case "stop":
			if err := run.Stop(); err != nil {
				_ = emit.Emit(emitter.EventTypeError, emitter.StageTypeUnknown, map[string]interface{}{
					"message": fmt.Sprintf("erro ao parar processo: %v", err),
					"source":  "agent",
				})
			}

		case "run":
			if err := sanitizer.ValidateCommand(cmd.Bin, cmd.Args); err != nil {
				_ = emit.Emit(emitter.EventTypeError, emitter.StageTypeUnknown, map[string]interface{}{
					"message": fmt.Sprintf("comando inválido: %v", err),
					"source":  "agent",
				})
				continue
			}

			stage := cmd.Stage
			if stage == "" {
				stage = emitter.StageTypeUnknown
			}

			var ctx context.Context
			var cancel context.CancelFunc
			if cmd.Timeout > 0 {
				ctx, cancel = context.WithTimeout(
					context.Background(),
					time.Duration(cmd.Timeout)*time.Second,
				)
			} else {
				ctx, cancel = context.WithCancel(context.Background())
			}

			go func(c emitter.AgentCommand, cancelFn context.CancelFunc) {
				defer cancelFn()
				if err := run.Start(ctx, c.Bin, c.Args, stage); err != nil {
					_ = emit.Emit(emitter.EventTypeError, stage, map[string]interface{}{
						"message": fmt.Sprintf("erro ao executar processo: %v", err),
						"source":  "agent",
					})
				}
			}(cmd, cancel)

		default:
			_ = emit.Emit(emitter.EventTypeError, emitter.StageTypeUnknown, map[string]interface{}{
				"message": fmt.Sprintf("comando desconhecido: %q", cmd.Command),
				"source":  "agent",
			})
		}
	}

	_ = run.Stop()
}
