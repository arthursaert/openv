package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arthur/openv/core"
)

// Status mostra o estado atual do repositório
func Status() {
	config, err := core.LoadConfig(".openv")
	if err != nil {
		fmt.Printf("❌ Erro ao carregar .openv: %v\n", err)
		fmt.Println("💡 Dica: Execute 'openv init' para inicializar um repositório")
		return
	}

	// Detecta arquivos modificados
	changedFiles, err := core.DetectChangedFiles(config)
	if err != nil {
		fmt.Printf("❌ Erro ao detectar mudanças: %v\n", err)
		return
	}

	// Busca arquivos não versionados
	untracked := findUntrackedFiles(config)

	// Busca arquivos deletados (estavam versionados, sumiram)
	deleted := findDeletedFiles(config)

	// === CABEÇALHO ===
	fmt.Println("📊 Status do Repositório")
	fmt.Println(strings.Repeat("=", 40))

	// === ARQUIVOS MODIFICADOS ===
	fmt.Println()
	if len(changedFiles) == 0 {
		fmt.Println("✅ Nenhum arquivo modificado")
	} else {
		fmt.Printf("📝 Arquivos modificados: %d\n", len(changedFiles))
		for _, f := range changedFiles {
			printFileStatus(f)
		}
	}

	// === ARQUIVOS DELETADOS ===
	fmt.Println()
	if len(deleted) == 0 {
		fmt.Println("✅ Nenhum arquivo deletado")
	} else {
		fmt.Printf("❌ Arquivos deletados: %d\n", len(deleted))
		for _, path := range deleted {
			fmt.Printf("   ❌ %s (deletado)\n", path)
		}
	}

	// === ARQUIVOS NÃO VERSIONADOS ===
	fmt.Println()
	if len(untracked) == 0 {
		fmt.Println("✅ Nenhum arquivo não versionado")
	} else {
		fmt.Printf("🆕 Arquivos não versionados: %d\n", len(untracked))
		for _, f := range untracked {
			fmt.Printf("   ?? %s\n", f)
		}
	}

	// === RESUMO ===
	fmt.Println()
	fmt.Println("📈 Resumo:")
	fmt.Printf("   • Commits totais: %d\n", len(config.Commits))
	fmt.Printf("   • Arquivos versionados: %d\n", countUniqueVersionedFiles(config))
	fmt.Printf("   • Modificados: %d\n", len(changedFiles))
	fmt.Printf("   • Deletados: %d\n", len(deleted))
	fmt.Printf("   • Não versionados: %d\n", len(untracked))

	// === ÚLTIMO COMMIT ===
	if len(config.Commits) > 0 {
		lastCommit := config.Commits[len(config.Commits)-1]
		fmt.Println()
		fmt.Println("🕐 Último commit:")
		fmt.Printf("   • ID: %s\n", lastCommit.ID[:8])
		fmt.Printf("   • Mensagem: %s\n", lastCommit.Message)
		fmt.Printf("   • Arquivos: %d\n", len(lastCommit.Files))
	}
}

// printFileStatus exibe o status de um arquivo com ícones apropriados
func printFileStatus(f core.FileChange) {
	if f.Binary {
		fmt.Printf("   📦 %s (binário, %s)\n", f.Path, formatSize(f.Size))
	} else if f.Compressed {
		fmt.Printf("   🗜️  %s (comprimido, %d mudança(s))\n", f.Path, len(f.LineChanges))
	} else if len(f.LineChanges) > 0 {
		fmt.Printf("   ✏️  %s (%d mudança(s))\n", f.Path, len(f.LineChanges))
	} else {
		fmt.Printf("   • %s\n", f.Path)
	}
}

// findUntrackedFiles encontra arquivos que existem mas não estão versionados
func findUntrackedFiles(config core.OpenVConfig) []string {
	var untracked []string
	ignoreList := loadIgnoreFile(".openvignore")

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Ignora diretórios e arquivos do sistema
		if info.IsDir() {
			return nil
		}
		if path == ".openv" || path == ".openvignore" {
			return nil
		}
		if strings.HasPrefix(path, ".git/") || path == ".git" {
			return nil
		}
		if strings.HasPrefix(path, "build/") || strings.HasPrefix(path, "bin/") {
			return nil
		}
		if shouldIgnore(path, ignoreList) {
			return nil
		}

		// Verifica se já está versionado em algum commit
		found := false
		for _, c := range config.Commits {
			for _, f := range c.Files {
				if f.Path == path {
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		if !found {
			untracked = append(untracked, path)
		}

		return nil
	})

	return untracked
}

// findDeletedFiles encontra arquivos que estavam versionados mas foram deletados
func findDeletedFiles(config core.OpenVConfig) []string {
	var deleted []string

	// Pega o último commit para comparar
	if len(config.Commits) == 0 {
		return deleted
	}

	lastCommit := config.Commits[len(config.Commits)-1]

	for _, f := range lastCommit.Files {
		// Verifica se o arquivo ainda existe
		if _, err := os.Stat(f.Path); os.IsNotExist(err) {
			deleted = append(deleted, f.Path)
		}
	}

	return deleted
}

// countUniqueVersionedFiles conta quantos arquivos únicos já foram versionados
func countUniqueVersionedFiles(config core.OpenVConfig) int {
	seen := make(map[string]bool)
	for _, c := range config.Commits {
		for _, f := range c.Files {
			seen[f.Path] = true
		}
	}
	return len(seen)
}

// formatSize formata tamanho em bytes para formato legível
func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

// loadIgnoreFile carrega a lista de arquivos ignorados do .openvignore
func loadIgnoreFile(path string) []string {
	var ignores []string

	data, err := os.ReadFile(path)
	if err != nil {
		return ignores
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			ignores = append(ignores, line)
		}
	}

	return ignores
}

// shouldIgnore verifica se um arquivo deve ser ignorado
func shouldIgnore(path string, ignoreList []string) bool {
	for _, pattern := range ignoreList {
		if pattern == path {
			return true
		}
		if strings.HasPrefix(pattern, "*") {
			ext := strings.TrimPrefix(pattern, "*")
			if strings.HasSuffix(path, ext) {
				return true
			}
		}
		if strings.HasSuffix(pattern, "/") {
			dir := strings.TrimSuffix(pattern, "/")
			if strings.HasPrefix(path, dir+"/") || path == dir {
				return true
			}
		}
	}
	return false
}
