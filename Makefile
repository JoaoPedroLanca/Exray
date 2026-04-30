AGENT_DIR     := packages/agent
AGENT_BIN_DIR := packages/agent/bin
EXT_BIN_DIR   := packages/extension/bin
CMD_PATH      := ./cmd/agent

.PHONY: build-agent build-agent-all test-agent lint-agent clean

## Compila o agent para a plataforma atual e copia para packages/extension/bin/
build-agent:
	@mkdir -p $(AGENT_BIN_DIR) $(EXT_BIN_DIR)
	cd $(AGENT_DIR) && go build -ldflags="-s -w" \
		-o ../../$(AGENT_BIN_DIR)/agent$(shell go env GOEXE) $(CMD_PATH)
	cp $(AGENT_BIN_DIR)/agent$(shell go env GOEXE) $(EXT_BIN_DIR)/
	@echo "✅ agent compilado para $$(go env GOOS)/$$(go env GOARCH)"

## Compila para as 5 plataformas suportadas
build-agent-all:
	@mkdir -p $(AGENT_BIN_DIR) $(EXT_BIN_DIR)
	cd $(AGENT_DIR) && GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ../../$(AGENT_BIN_DIR)/agent-windows-amd64.exe $(CMD_PATH)
	cd $(AGENT_DIR) && GOOS=linux   GOARCH=amd64 go build -ldflags="-s -w" -o ../../$(AGENT_BIN_DIR)/agent-linux-amd64      $(CMD_PATH)
	cd $(AGENT_DIR) && GOOS=linux   GOARCH=arm64 go build -ldflags="-s -w" -o ../../$(AGENT_BIN_DIR)/agent-linux-arm64      $(CMD_PATH)
	cd $(AGENT_DIR) && GOOS=darwin  GOARCH=amd64 go build -ldflags="-s -w" -o ../../$(AGENT_BIN_DIR)/agent-darwin-amd64    $(CMD_PATH)
	cd $(AGENT_DIR) && GOOS=darwin  GOARCH=arm64 go build -ldflags="-s -w" -o ../../$(AGENT_BIN_DIR)/agent-darwin-arm64    $(CMD_PATH)
	cp $(AGENT_BIN_DIR)/agent-windows-amd64.exe $(EXT_BIN_DIR)/
	cp $(AGENT_BIN_DIR)/agent-linux-amd64       $(EXT_BIN_DIR)/
	cp $(AGENT_BIN_DIR)/agent-linux-arm64       $(EXT_BIN_DIR)/
	cp $(AGENT_BIN_DIR)/agent-darwin-amd64      $(EXT_BIN_DIR)/
	cp $(AGENT_BIN_DIR)/agent-darwin-arm64      $(EXT_BIN_DIR)/
	@echo "✅ 5 binários compilados em $(AGENT_BIN_DIR)/ e copiados para $(EXT_BIN_DIR)/"

## Roda todos os testes Go com race detector
test-agent:
	cd $(AGENT_DIR) && go test ./... -race -count=1 -v

## Lint com golangci-lint
## Instalar: go install ...
lint-agent:
	cd $(AGENT_DIR) && golangci-lint run ./...

## Remove todos os artefatos compilados
clean:
	rm -rf $(AGENT_BIN_DIR)
	rm -rf packages/extension/dist
	rm -rf packages/webview/dist
	find . -name "*.vsix" -delete
	@echo "✅ Artefatos removidos"