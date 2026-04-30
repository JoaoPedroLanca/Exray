package classifier

import (
	"regexp"
	"strings"

	"github.com/JoaoPedroLanca/Exray/agent/internal/emitter"
)

// matches build-tool error codes like MSB3021:, CS0001:, TS2304:
var buildErrorCodeRe = regexp.MustCompile(`\b[a-z]+\d+:`)

func Classify(line string) emitter.EventType {
	lower := strings.ToLower(line)

	errorPatterns := []string{
		"panic: ",
		"fatal: ",
		"exception",
		"critical: ",
		"err: ",
		"failed: ",
		" fail",
		"error:",
		"error ",
	}

	for _, pattern := range errorPatterns {
		if strings.Contains(lower, pattern) {
			return emitter.EventTypeError
		}
	}

	if buildErrorCodeRe.MatchString(lower) {
		return emitter.EventTypeError
	}

	warningPatterns := []string{
		"warning:",
		"warn:",
		"deprecated",
		"caution",
		"warning",
		"warn",
	}

	for _, pattern := range warningPatterns {
		if strings.Contains(lower, pattern) {
			return emitter.EventTypeWarning
		}
	}

	return emitter.EventTypeLog
}
