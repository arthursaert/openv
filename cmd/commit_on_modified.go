package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/arthur/openv/core"
)

func CommitOnModified(args []string) {
	if len(args) < 2 {
		fmt.Println("❌ Uso: openv --commit-on-m \"mensagem\" <arquivo|diretório>")
		fmt.Println("   Ex: openv --commit-on-m \"Commit\" .")
		fmt.Println("       openv --commit-on-m \"Commit\" arquivo.txt")
		fmt.Println("       openv --commit-on-m \"Commit\" pasta/")
		os.Exit(1)
	}

	message := args[0]
	targetPath := args[1] // ✅ Mantém como o usuário digitou!

	// ✅ Verifica existência sem converter para absoluto
	checkPath := targetPath
	if !filepath.IsAbs(targetPath) {
		checkPath = filepath.Join(".", targetPath)
	}

	info, err := os.Stat(checkPath)
	if os.IsNotExist(err) {
		fmt.Printf("❌ Arquivo ou diretório não existe: %s\n", targetPath)
		os.Exit(1)
	}
	if err != nil {
		fmt.Printf("❌ Erro ao acessar caminho: %v\n", err)
		os.Exit(1)
	}

	// Carrega config
	config, err := core.LoadConfig(".openv")
	if err != nil {
		fmt.Printf("❌ Erro ao carregar .openv: %v\n", err)
		os.Exit(1)
	}

	var changedFiles []core.FileChange

	// ✅ Passa targetPath (não absoluto) pro tracker
	if info.IsDir() {
		changedFiles, err = core.DetectChangedFilesInDir(config, targetPath)
		if err != nil {
			fmt.Printf("❌ Erro ao detectar mudanças: %v\n", err)
			os.Exit(1)
		}
	} else {
		changedFiles, err = core.DetectChangedFile(config, targetPath)
		if err != nil {
			fmt.Printf("❌ Erro ao detectar mudanças: %v\n", err)
			os.Exit(1)
		}
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
		Dir:     targetPath,
	}

	config.Commits = append(config.Commits, commit)

	// Salva config atualizada
	err = core.SaveConfig(".openv", config)
	if err != nil {
		fmt.Printf("❌ Erro ao salvar .openv: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Commit criado: %s\n", commit.ID[:8])
	fmt.Printf("📁 %d arquivo(s) modificado(s) commitado(s)\n", len(changedFiles))
	for _, f := range changedFiles {
		fmt.Printf("   • %s\n", f.Path)
	}
}
