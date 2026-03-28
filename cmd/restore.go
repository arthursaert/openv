package cmd

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/arthur/openv/core"
)

func Restore(args []string) {
	if len(args) < 2 {
		fmt.Println("❌ Uso: openv restore <commit_id> <arquivo>")
		fmt.Println("   Ex: openv restore abc1234 teste.txt")
		fmt.Println("       openv restore 2 teste.txt  (usa índice do commit)")
		os.Exit(1)
	}

	commitRef := args[0]
	filePath := args[1]

	// Carrega config
	config, err := core.LoadConfig(".openv")
	if err != nil {
		fmt.Printf("❌ Erro ao carregar .openv: %v\n", err)
		os.Exit(1)
	}

	// Encontra o commit
	var commitIndex int
	var targetCommit *core.Commit

	// Tenta buscar por índice primeiro
	if idx, err := strconv.Atoi(commitRef); err == nil {
		if idx < 0 || idx >= len(config.Commits) {
			fmt.Printf("❌ Commit %d não existe\n", idx)
			os.Exit(1)
		}
		commitIndex = idx
		targetCommit = &config.Commits[idx]
	} else {
		// Busca por ID parcial
		found := false
		for i, c := range config.Commits {
			if strings.HasPrefix(c.ID, commitRef) || commitRef == c.ID {
				commitIndex = i
				targetCommit = &c
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("❌ Commit %s não encontrado\n", commitRef)
			os.Exit(1)
		}
	}

	// ✅ Encontra o arquivo no commit (múltiplas tentativas de match)
	var targetFile *core.FileChange
	for _, f := range targetCommit.Files {
		// Match 1: Caminho direto
		if f.Path == filePath {
			targetFile = &f
			break
		}
		// Match 2: Caminho absoluto
		absInput, _ := filepath.Abs(filePath)
		absInCommit, _ := filepath.Abs(f.Path)
		if absInput == absInCommit {
			targetFile = &f
			break
		}
		// Match 3: Nome do arquivo (último componente)
		if filepath.Base(f.Path) == filepath.Base(filePath) {
			targetFile = &f
			break
		}
	}

	if targetFile == nil {
		fmt.Printf("❌ Arquivo %s não encontrado no commit %d\n", filePath, commitIndex)
		fmt.Printf("📁 Arquivos nesse commit:\n")
		for _, f := range targetCommit.Files {
			fmt.Printf("   • %s\n", f.Path)
		}
		os.Exit(1)
	}

	// Decodifica conteúdo
	content, err := base64.StdEncoding.DecodeString(targetFile.Content)
	if err != nil {
		fmt.Printf("❌ Erro ao decodificar conteúdo: %v\n", err)
		os.Exit(1)
	}

	// Restaura arquivo
	err = os.WriteFile(filePath, content, 0644)
	if err != nil {
		fmt.Printf("❌ Erro ao restaurar arquivo: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Arquivo %s restaurado do commit %s\n", filePath, targetCommit.ID[:8])
	fmt.Printf("📝 Mensagem do commit: %s\n", targetCommit.Message)
	fmt.Printf("📊 %d mudança(s) de linha nesse arquivo:\n", len(targetFile.LineChanges))

	for _, change := range targetFile.LineChanges {
		switch change.ChangeType {
		case "added":
			fmt.Printf("   ➕ Linha %d: %s\n", change.LineNumber, change.NewContent)
		case "deleted":
			fmt.Printf("   ➖ Linha %d: %s\n", change.LineNumber, change.OldContent)
		case "modified":
			fmt.Printf("   ✏️  Linha %d: %s → %s\n", change.LineNumber, change.OldContent, change.NewContent)
		}
	}
}
