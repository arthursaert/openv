.PHONY: all build clean linux windows mac all-platforms

# Pasta de saída
BIN_DIR := bin

# Cria pasta bin/ se não existir
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

all: build

build: linux

# Linux 64-bit + 32-bit
linux: $(BIN_DIR)
	@echo "🐧 Compilando para Linux..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_DIR)/openv-linux-amd64 .
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_DIR)/openv-server-linux-amd64 ./server
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o $(BIN_DIR)/openv-linux-386 .
	GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -o $(BIN_DIR)/openv-server-linux-386 ./server
	@echo "✅ Linux (amd64 + 386) compilado!"

# Windows 64-bit + 32-bit
windows: $(BIN_DIR)
	@echo "🪟 Compilando para Windows..."
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_DIR)/openv-windows-amd64.exe .
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_DIR)/openv-server-windows-amd64.exe ./server
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -o $(BIN_DIR)/openv-windows-386.exe .
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -o $(BIN_DIR)/openv-server-windows-386.exe ./server
	@echo "✅ Windows (amd64 + 386) compilado!"

# Mac Intel + Apple Silicon
mac: $(BIN_DIR)
	@echo "🍎 Compilando para macOS..."
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_DIR)/openv-darwin-amd64 .
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o $(BIN_DIR)/openv-server-darwin-amd64 ./server
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o $(BIN_DIR)/openv-darwin-arm64 .
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o $(BIN_DIR)/openv-server-darwin-arm64 ./server
	@echo "✅ macOS (amd64 + arm64) compilado!"

# Compila TUDO de uma vez
all-platforms: $(BIN_DIR)
	@echo "🌍 Compilando para TODAS as plataformas..."
	@$(MAKE) linux
	@$(MAKE) windows
	@$(MAKE) mac
	@echo ""
	@echo "✅ TODAS as plataformas compiladas!"
	@echo "📁 Binários em: $(BIN_DIR)/"
	@ls -lh $(BIN_DIR)/

# Limpa tudo
clean:
	@echo "🧹 Limpando..."
	rm -rf $(BIN_DIR)
	@echo "✅ Limpo!"

# Ajuda
help:
	@echo "🚀 OpenV Build System"
	@echo ""
	@echo "Uso:"
	@echo "  make all           - Compila Linux (padrão)"
	@echo "  make linux         - Linux 64-bit + 32-bit"
	@echo "  make windows       - Windows 64-bit + 32-bit"
	@echo "  make mac           - macOS Intel + Apple Silicon"
	@echo "  make all-platforms - TODAS as plataformas"
	@echo "  make clean         - Limpa binários"
	@echo "  make help          - Mostra esta ajuda"