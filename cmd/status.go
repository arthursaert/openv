package cmd

import (
	"fmt"

	"github.com/arthur/openv/core"
)

func Status() {
	config, err := core.LoadConfig(".openv")
	if err != nil {
		fmt.Printf("❌ Erro: %v\n", err)
		return
	}

	changedFiles, err := core.DetectChangedFiles(config)
	if err != nil {
		fmt.Printf("❌ Erro: %v\n", err)
		return
	}

	fmt.Println("📊 Status do Repositório")
	fmt.Println("=" + repeat("=", 39))

	if len(changedFiles) == 0 {
		fmt.Println("✅ Nenhum arquivo modificado")
	} else {
		fmt.Printf("📝 Arquivos modificados: %d\n", len(changedFiles))
		for _, f := range changedFiles {
			fmt.Printf("   M %s\n", f.Path)
		}
	}

	fmt.Println()
	fmt.Printf("📦 Total de commits: %d\n", len(config.Commits))
	if len(config.Commits) > 0 {
		lastCommit := config.Commits[len(config.Commits)-1]
		fmt.Printf("🕐 Último commit: %s\n", lastCommit.ID[:8])
	}
}

func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
