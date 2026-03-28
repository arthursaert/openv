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

// ✅ FileChange NÃO está aqui - já está em config.go!

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

			changed = append(changed, FileChange{
				Path:     path,
				Hash:     hash,
				Size:     info.Size(),
				Modified: info.ModTime().Format(time.RFC3339),
				Content:  content,
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

			changed = append(changed, FileChange{
				Path:     path,
				Hash:     hash,
				Size:     info.Size(),
				Modified: info.ModTime().Format(time.RFC3339),
				Content:  content,
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

		changed = append(changed, FileChange{
			Path:     filePath,
			Hash:     hash,
			Size:     info.Size(),
			Modified: info.ModTime().Format(time.RFC3339),
			Content:  content,
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

// ✅ NOVA FUNÇÃO: Compara dois conteúdos e retorna as mudanças por linha
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
			// Linha adicionada
			changeType = "added"
		} else if i >= len(newLines) {
			// Linha deletada
			changeType = "deleted"
		} else if oldLine != newLine {
			// Linha modificada
			changeType = "modified"
		} else {
			// Linha igual, não adiciona
			continue
		}

		changes = append(changes, LineChange{
			LineNumber: i + 1, // Linhas começam em 1
			OldContent: oldLine,
			NewContent: newLine,
			ChangeType: changeType,
		})
	}

	return changes
}

// ✅ NOVA FUNÇÃO: Aplica diff para restaurar conteúdo
func ApplyDiff(baseContent string, changes []LineChange, revert bool) string {
	lines := strings.Split(baseContent, "\n")

	for _, change := range changes {
		lineIdx := change.LineNumber - 1 // Converte para 0-based

		if lineIdx >= len(lines) && change.ChangeType == "added" && !revert {
			// Adiciona linha no final
			lines = append(lines, change.NewContent)
		} else if lineIdx < len(lines) {
			if revert {
				// Reverte: usa OldContent
				if change.ChangeType == "added" {
					// Remove linha adicionada
					lines = append(lines[:lineIdx], lines[lineIdx+1:]...)
				} else if change.ChangeType == "deleted" {
					// Restaura linha deletada
					lines = append(lines[:lineIdx], append([]string{change.OldContent}, lines[lineIdx:]...)...)
				} else {
					// Restaura modificação
					lines[lineIdx] = change.OldContent
				}
			} else {
				// Aplica: usa NewContent
				if change.ChangeType == "added" {
					lines = append(lines[:lineIdx], append([]string{change.NewContent}, lines[lineIdx:]...)...)
				} else if change.ChangeType == "deleted" {
					lines = append(lines[:lineIdx], lines[lineIdx+1:]...)
				} else {
					lines[lineIdx] = change.NewContent
				}
			}
		}
	}

	return strings.Join(lines, "\n")
}

// ✅ NOVA FUNÇÃO: Busca arquivo em um commit específico
func FindFileInCommit(config OpenVConfig, commitIndex int, filePath string) (*FileChange, error) {
	if commitIndex < 0 || commitIndex >= len(config.Commits) {
		return nil, fmt.Errorf("commit %d não existe", commitIndex)
	}

	commit := config.Commits[commitIndex]
	for _, file := range commit.Files {
		if file.Path == filePath {
			return &file, nil
		}
	}

	return nil, fmt.Errorf("arquivo %s não encontrado no commit %d", filePath, commitIndex)
}
