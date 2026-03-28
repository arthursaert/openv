package core

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ✅ FileChange e LineChange NÃO estão aqui - já estão em config.go!

// DetectChangedFiles detecta arquivos modificados no diretório atual
func DetectChangedFiles(config OpenVConfig) ([]FileChange, error) {
	var changed []FileChange

	ignoreList := loadIgnoreFile(".openvignore")

	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == ".openv" || path == ".openvignore" {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasPrefix(path, ".git/") || path == ".git" {
			return nil
		}

		if strings.HasPrefix(path, "build/") {
			return nil
		}

		if shouldIgnore(path, ignoreList) {
			return nil
		}

		hash, err := calculateHash(path)
		if err != nil {
			return err
		}

		exists := false
		if len(config.Commits) > 0 {
			lastCommit := config.Commits[len(config.Commits)-1]
			for _, file := range lastCommit.Files {
				if file.Path == path && file.Hash == hash {
					exists = true
					break
				}
			}
		}

		if !exists {
			content, err := readFileContent(path)
			if err != nil {
				return err
			}

			// ✅ Calcula diff se tiver versão anterior
			var lineChanges []LineChange
			if len(config.Commits) > 0 {
				lastCommit := config.Commits[len(config.Commits)-1]
				for _, f := range lastCommit.Files {
					if f.Path == path {
						oldContent, _ := base64.StdEncoding.DecodeString(f.Content)
						lineChanges = CalculateDiff(string(oldContent), string(content))
						break
					}
				}
			}

			changed = append(changed, FileChange{
				Path:        path,
				Hash:        hash,
				Size:        info.Size(),
				Modified:    info.ModTime().Format(time.RFC3339),
				Content:     content,
				LineChanges: lineChanges,
			})
		}

		return nil
	})

	return changed, err
}

// DetectChangedFilesInDir detecta arquivos modificados em um diretório específico
func DetectChangedFilesInDir(config OpenVConfig, dirPath string) ([]FileChange, error) {
	var changed []FileChange

	ignoreList := loadIgnoreFile(".openvignore")

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".openv") || strings.HasSuffix(path, ".openvignore") {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if strings.Contains(path, ".git/") {
			return nil
		}

		if strings.Contains(path, "/build/") {
			return nil
		}

		if shouldIgnore(path, ignoreList) {
			return nil
		}

		hash, err := calculateHash(path)
		if err != nil {
			return err
		}

		exists := false
		if len(config.Commits) > 0 {
			lastCommit := config.Commits[len(config.Commits)-1]
			for _, file := range lastCommit.Files {
				if file.Path == path && file.Hash == hash {
					exists = true
					break
				}
			}
		}

		if !exists {
			content, err := readFileContent(path)
			if err != nil {
				return err
			}

			// ✅ Calcula diff se tiver versão anterior
			var lineChanges []LineChange
			if len(config.Commits) > 0 {
				lastCommit := config.Commits[len(config.Commits)-1]
				for _, f := range lastCommit.Files {
					if f.Path == path {
						oldContent, _ := base64.StdEncoding.DecodeString(f.Content)
						lineChanges = CalculateDiff(string(oldContent), string(content))
						break
					}
				}
			}

			// ✅ Converte para caminho relativo
			relPath, err := filepath.Rel(".", path)
			if err != nil {
				relPath = path
			}

			changed = append(changed, FileChange{
				Path:        relPath,
				Hash:        hash,
				Size:        info.Size(),
				Modified:    info.ModTime().Format(time.RFC3339),
				Content:     content,
				LineChanges: lineChanges,
			})
		}

		return nil
	})

	return changed, err
}

// DetectChangedFile detecta mudanças em um arquivo único
func DetectChangedFile(config OpenVConfig, filePath string) ([]FileChange, error) {
	var changed []FileChange

	info, err := os.Stat(filePath)
	if err != nil {
		return changed, err
	}

	if info.IsDir() {
		return changed, fmt.Errorf("caminho é um diretório: use DetectChangedFilesInDir")
	}

	ignoreList := loadIgnoreFile(".openvignore")
	if shouldIgnore(filePath, ignoreList) {
		return changed, nil
	}

	hash, err := calculateHash(filePath)
	if err != nil {
		return changed, err
	}

	exists := false
	if len(config.Commits) > 0 {
		lastCommit := config.Commits[len(config.Commits)-1]
		for _, file := range lastCommit.Files {
			// ✅ Compara caminhos relativos
			if file.Path == filePath && file.Hash == hash {
				exists = true
				break
			}
		}
	}

	if !exists {
		content, err := readFileContent(filePath)
		if err != nil {
			return changed, err
		}

		// ✅ Calcula diff se tiver versão anterior
		var lineChanges []LineChange
		if len(config.Commits) > 0 {
			lastCommit := config.Commits[len(config.Commits)-1]
			for _, f := range lastCommit.Files {
				if f.Path == filePath {
					oldContent, _ := base64.StdEncoding.DecodeString(f.Content)
					lineChanges = CalculateDiff(string(oldContent), string(content))
					break
				}
			}
		}

		// ✅ Converte para caminho relativo
		relPath, err := filepath.Rel(".", filePath)
		if err != nil {
			relPath = filePath
		}

		changed = append(changed, FileChange{
			Path:        relPath,
			Hash:        hash,
			Size:        info.Size(),
			Modified:    info.ModTime().Format(time.RFC3339),
			Content:     content,
			LineChanges: lineChanges,
		})
	}

	return changed, nil
}

// calculateHash calcula o hash SHA256 de um arquivo
func calculateHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// readFileContent lê o conteúdo do arquivo e retorna em base64
func readFileContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// DecodeFileContent decodifica conteúdo de base64 para bytes
func DecodeFileContent(content string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(content)
}

// CalculateDiff compara dois conteúdos e retorna as mudanças por linha
func CalculateDiff(oldContent, newContent string) []LineChange {
	var changes []LineChange

	oldLines := strings.Split(oldContent, "\n")
	newLines := strings.Split(newContent, "\n")

	maxLines := len(oldLines)
	if len(newLines) > maxLines {
		maxLines = len(newLines)
	}

	for i := 0; i < maxLines; i++ {
		var oldLine, newLine string
		var changeType string

		if i < len(oldLines) {
			oldLine = oldLines[i]
		}
		if i < len(newLines) {
			newLine = newLines[i]
		}

		if i >= len(oldLines) {
			changeType = "added"
		} else if i >= len(newLines) {
			changeType = "deleted"
		} else if oldLine != newLine {
			changeType = "modified"
		} else {
			continue
		}

		changes = append(changes, LineChange{
			LineNumber: i + 1,
			OldContent: oldLine,
			NewContent: newLine,
			ChangeType: changeType,
		})
	}

	return changes
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
