package cmd

import (
	"fmt"
	"os"

	"github.com/arthur/openv/core"
)

func Init() {
	// Verifica se já existe .openv
	if _, err := os.Stat(".openv"); err == nil {
		fmt.Println("⚠️  Repositório OpenV já inicializado!")
		return
	}

	// Cria arquivo .openv
	config := core.OpenVConfig{
		RepositoryID: core.GenerateRepoID(),
		Version:      "1.0",
		Commits:      []core.Commit{},
	}

	err := core.SaveConfig(".openv", config)
	if err != nil {
		fmt.Printf("❌ Erro ao inicializar: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Repositório OpenV inicializado!")
	fmt.Println("📁 Arquivo .openv criado")
}
