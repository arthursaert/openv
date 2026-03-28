package cmd

import (
	"fmt"
	"strings"

	"github.com/arthur/openv/core"
)

func Log() {
	config, err := core.LoadConfig(".openv")
	if err != nil {
		fmt.Printf("❌ Erro: %v\n", err)
		return
	}

	if len(config.Commits) == 0 {
		fmt.Println("📭 Nenhum commit ainda")
		return
	}

	fmt.Println("📜 Histórico de Commits")
	fmt.Println(strings.Repeat("=", 40))

	for i := len(config.Commits) - 1; i >= 0; i-- {
		commit := config.Commits[i]
		fmt.Printf("\n🔖 [%d] %s\n", i, commit.ID[:8])
		fmt.Printf("📝 %s\n", commit.Message)
		fmt.Printf("📁 %d arquivo(s)\n", len(commit.Files))
		fmt.Printf("📂 Diretório: %s\n", commit.Dir)

		// ✅ Mostra mudanças por linha
		if len(commit.Files) > 0 {
			fmt.Println("   ─────────────────────────────")
			for _, file := range commit.Files {
				fmt.Printf("   📄 %s\n", file.Path)
				if len(file.LineChanges) > 0 {
					for _, change := range file.LineChanges {
						switch change.ChangeType {
						case "added":
							fmt.Printf("      ➕ +%d: %s\n", change.LineNumber, change.NewContent)
						case "deleted":
							fmt.Printf("      ➖ -%d: %s\n", change.LineNumber, change.OldContent)
						case "modified":
							fmt.Printf("      ✏️  ~%d: %s → %s\n", change.LineNumber, change.OldContent, change.NewContent)
						}
					}
				} else {
					fmt.Println("      (sem mudanças de linha detectadas)")
				}
			}
		}
	}
}
