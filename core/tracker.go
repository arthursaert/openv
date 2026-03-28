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

// isBinary detecta se um arquivo é binário pela extensão
func isBinary(path string) bool {
	binaryExts := []string{
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".ico", ".webp",
		".mp4", ".avi", ".mkv", ".mov", ".wmv",
		".mp3", ".wav", ".flac", ".ogg",
		".pdf", ".doc", ".docx", ".xls", ".xlsx",
		".zip", ".rar", ".tar", ".gz", ".7z",
		".exe", ".dll", ".so", ".bin",
		".psd", ".ai", ".eps",
	}

	ext := strings.ToLower(filepath.Ext(path))
	for _, e := range binaryExts {
		if ext == e {
			return true
		}
	}
	return false
}

// isBinaryByContent detecta se um arquivo é binário pelo conteúdo
func isBinaryByContent(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 8192)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false
	}

	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return true
		}
	}

	return false
}

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

		// ✅ Busca conteúdo anterior para calcular diff
		var oldContent string
		var lineChanges []LineChange
		if len(config.Commits) > 0 {
			lastCommit := config.Commits[len(config.Commits)-1]
			for _, f := range lastCommit.Files {
				if f.Path == path {
					oldContentBytes, _ := base64.StdEncoding.DecodeString(f.Content)

					// ✅ Se estava comprimido, descomprime
					if f.Compressed {
						oldContentBytes, _ = DecompressGzip(oldContentBytes)
					}

					oldContent = string(oldContentBytes)
					break
				}
			}
		}

		// ✅ Lê o conteúdo atual
		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// ✅ Detecta se é binário
		binary := isBinary(path) || isBinaryByContent(path)

		// ✅ Calcula diff SÓ se for texto
		if !binary && oldContent != "" {
			lineChanges = CalculateDiff(oldContent, string(contentBytes))
		}

		// ✅ Comprime com GZIP se for binário OU se for maior que 10KB
		var finalContent []byte
		var compressed bool

		if binary || len(contentBytes) > 10240 {
			compressedBytes, err := CompressGzip(contentBytes)
			if err == nil {
				finalContent = compressedBytes
				compressed = true
			} else {
				finalContent = contentBytes
				compressed = false
			}
		} else {
			finalContent = contentBytes
			compressed = false
		}

		// ✅ Codifica em base64
		content := base64.StdEncoding.EncodeToString(finalContent)

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
			Compressed:  compressed,
			Binary:      binary,
			LineChanges: lineChanges,
		})

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

		// ✅ Busca conteúdo anterior para calcular diff
		var oldContent string
		var lineChanges []LineChange
		if len(config.Commits) > 0 {
			lastCommit := config.Commits[len(config.Commits)-1]
			for _, f := range lastCommit.Files {
				if f.Path == path {
					oldContentBytes, _ := base64.StdEncoding.DecodeString(f.Content)

					// ✅ Se estava comprimido, descomprime
					if f.Compressed {
						oldContentBytes, _ = DecompressGzip(oldContentBytes)
					}

					oldContent = string(oldContentBytes)
					break
				}
			}
		}

		// ✅ Lê o conteúdo atual
		contentBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// ✅ Detecta se é binário
		binary := isBinary(path) || isBinaryByContent(path)

		// ✅ Calcula diff SÓ se for texto
		if !binary && oldContent != "" {
			lineChanges = CalculateDiff(oldContent, string(contentBytes))
		}

		// ✅ Comprime com GZIP se for binário OU se for maior que 10KB
		var finalContent []byte
		var compressed bool

		if binary || len(contentBytes) > 10240 {
			compressedBytes, err := CompressGzip(contentBytes)
			if err == nil {
				finalContent = compressedBytes
				compressed = true
			} else {
				finalContent = contentBytes
				compressed = false
			}
		} else {
			finalContent = contentBytes
			compressed = false
		}

		// ✅ Codifica em base64
		content := base64.StdEncoding.EncodeToString(finalContent)

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
			Compressed:  compressed,
			Binary:      binary,
			LineChanges: lineChanges,
		})

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

	// ✅ Busca conteúdo anterior para calcular diff
	var oldContent string
	var lineChanges []LineChange
	if len(config.Commits) > 0 {
		lastCommit := config.Commits[len(config.Commits)-1]
		for _, f := range lastCommit.Files {
			if f.Path == filePath {
				oldContentBytes, _ := base64.StdEncoding.DecodeString(f.Content)

				// ✅ Se estava comprimido, descomprime
				if f.Compressed {
					oldContentBytes, _ = DecompressGzip(oldContentBytes)
				}

				oldContent = string(oldContentBytes)
				break
			}
		}
	}

	// ✅ Lê o conteúdo atual
	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return changed, err
	}

	// ✅ Detecta se é binário
	binary := isBinary(filePath) || isBinaryByContent(filePath)

	// ✅ Calcula diff SÓ se for texto
	if !binary && oldContent != "" {
		lineChanges = CalculateDiff(oldContent, string(contentBytes))
	}

	// ✅ Comprime com GZIP se for binário OU se for maior que 10KB
	var finalContent []byte
	var compressed bool

	if binary || len(contentBytes) > 10240 {
		compressedBytes, err := CompressGzip(contentBytes)
		if err == nil {
			finalContent = compressedBytes
			compressed = true
		} else {
			finalContent = contentBytes
			compressed = false
		}
	} else {
		finalContent = contentBytes
		compressed = false
	}

	// ✅ Codifica em base64
	content := base64.StdEncoding.EncodeToString(finalContent)

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
		Compressed:  compressed,
		Binary:      binary,
		LineChanges: lineChanges,
	})

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
