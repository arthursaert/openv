package core

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
)

// LineChange representa uma mudança em uma linha específica
type LineChange struct {
	LineNumber int    `json:"line_number"`
	OldContent string `json:"old_content"`
	NewContent string `json:"new_content"`
	ChangeType string `json:"change_type"`
}

// FileChange representa um arquivo modificado com conteúdo e diff
type FileChange struct {
	Path        string       `json:"path"`
	Hash        string       `json:"hash"`
	Size        int64        `json:"size"`
	Modified    string       `json:"modified"`
	Content     string       `json:"content"`
	Compressed  bool         `json:"compressed,omitempty"` // ✅ NOVO!
	Binary      bool         `json:"binary,omitempty"`     // ✅ NOVO!
	LineChanges []LineChange `json:"line_changes,omitempty"`
}

// Commit representa um commit no histórico
type Commit struct {
	ID      string       `json:"id"`
	Message string       `json:"message"`
	Files   []FileChange `json:"files"`
	Dir     string       `json:"dir"`
}

// OpenVConfig representa a configuração do repositório
type OpenVConfig struct {
	RepositoryID string   `json:"repository_id"`
	Version      string   `json:"version"`
	Commits      []Commit `json:"commits"`
}

// GenerateRepositoryID gera um ID único para o repositório
func GenerateRepositoryID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// LoadConfig carrega a configuração do arquivo .openv
func LoadConfig(path string) (OpenVConfig, error) {
	var config OpenVConfig

	data, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// SaveConfig salva a configuração no arquivo .openv
func SaveConfig(path string, config OpenVConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
