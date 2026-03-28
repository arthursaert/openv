package core

import (
	"encoding/json"
	"os"
)

// ✅ LineChange declarado AQUI (único lugar!)
type LineChange struct {
	LineNumber int    `json:"line_number"`
	OldContent string `json:"old_content"`
	NewContent string `json:"new_content"`
	ChangeType string `json:"change_type"`
}

// ✅ FileChange com LineChanges
type FileChange struct {
	Path        string       `json:"path"`
	Hash        string       `json:"hash"`
	Size        int64        `json:"size"`
	Modified    string       `json:"modified"`
	Content     string       `json:"content"`
	LineChanges []LineChange `json:"line_changes,omitempty"`
}

type Commit struct {
	ID      string       `json:"id"`
	Message string       `json:"message"`
	Files   []FileChange `json:"files"`
	Dir     string       `json:"dir"`
}

type OpenVConfig struct {
	RepositoryID string   `json:"repository_id"`
	Version      string   `json:"version"`
	Commits      []Commit `json:"commits"`
}

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

func SaveConfig(path string, config OpenVConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
