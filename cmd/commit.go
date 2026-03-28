package cmd

import (
	"fmt"
	"os"

	"github.com/arthur/openv/core"
)

func Commit(args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Uso: openv commit \"mensagem\"")
		os.Exit(1)
	}

	message := args[0]

	// Carrega config
	config, err := core.LoadConfig(".openv")
	if err != nil {
		fmt.Printf("❌ Erro: %v\n", err)
		os.Exit(1)
	}

	// Detecta arquivos modificados
	changedFiles, err := core.DetectChangedFiles(config)
	if err != nil {
		fmt.Printf("❌ Erro ao detectar mudanças: %v\n", err)
		os.Exit(1)
	}

	if len(changedFiles) == 0 {
		fmt.Println("✅ Nenhum arquivo modificado")
		return
	}

	// Cria commit
	commit := core.Commit{
		ID:      core.GenerateCommitID(),
		Message: message,
		Files:   changedFiles,
	}

	config.Commits = append(config.Commits, commit)

	// Salva config atualizada
	err = core.SaveConfig(".openv", config)
	if err != nil {
		fmt.Printf("❌ Erro ao salvar: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Commit criado: %s\n", commit.ID[:8])
	fmt.Printf("📁 %d arquivo(s) versionado(s)\n", len(changedFiles))
}
