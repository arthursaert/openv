package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/arthur/openv/core"
)

type PushPayload struct {
	RepositoryID string           `json:"repository_id"`
	Commits      []core.Commit    `json:"commits"`
	Config       core.OpenVConfig `json:"config"`
}

func Push(args []string) {
	if len(args) < 1 {
		fmt.Println("❌ Uso: openv push <server_url>")
		os.Exit(1)
	}

	serverURL := args[0]

	// Carrega config
	config, err := core.LoadConfig(".openv")
	if err != nil {
		fmt.Printf("❌ Erro: %v\n", err)
		os.Exit(1)
	}

	// Prepara payload
	payload := PushPayload{
		RepositoryID: config.RepositoryID,
		Commits:      config.Commits,
		Config:       config,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("❌ Erro ao serializar: %v\n", err)
		os.Exit(1)
	}

	// Envia para servidor
	resp, err := http.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Erro de conexão: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		fmt.Println("✅ Push realizado com sucesso!")
		fmt.Printf("📡 Resposta: %s\n", string(body))
	} else {
		fmt.Printf("❌ Erro no servidor: %d\n", resp.StatusCode)
		fmt.Printf("📡 Resposta: %s\n", string(body))
	}
}
