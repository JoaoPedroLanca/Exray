package sanitizer

import (
	"fmt"
	"os/exec"
	"regexp"
)

var allowedBinaries = map[string]bool{
	"go":      true,
	"dotnet":  true,
	"npm":     true,
	"cargo":   true,
	"python":  true,
	"python3": true,
	"node":    true,
	"npx":     true,
	"yarn":    true,
	"bun":     true,
}

// Regex detecta metacharacters comuns usados para injeção de comandos em shell
var shellMetaRegex = regexp.MustCompile(`&&|\|\||;|\||>|<|` + "`" + `|\$\(|\$\{`)

func ValidateCommand(bin string, args []string) error {
	if !allowedBinaries[bin] {
		return fmt.Errorf("binário não permitido: %q - lista permitida: go, dotnet, npm, cargo, python, python3, node, npx, yarn, bun", bin)
	}

	if _, err := exec.LookPath(bin); err != nil {
		return fmt.Errorf("binário %q não encontrado no PATH: %w", bin, err)
	}

	for _, arg := range args {
		if shellMetaRegex.MatchString(arg) {
			return fmt.Errorf("argumento %q contém metacaracteres de shell proibido:", arg)
		}
	}

	return nil
}
