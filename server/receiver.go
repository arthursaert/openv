package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileChange com conteúdo
type FileChange struct {
	Path     string `json:"path"`
	Hash     string `json:"hash"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
	Content  string `json:"content"`
}

// Commit structure
type Commit struct {
	ID      string       `json:"id"`
	Message string       `json:"message"`
	Files   []FileChange `json:"files"`
	Dir     string       `json:"dir"`
}

// ReceivedData representa os dados recebidos
type ReceivedData struct {
	RepositoryID string          `json:"repository_id"`
	Commits      []Commit        `json:"commits"`
	Config       json.RawMessage `json:"config"`
}

// ✅ NOVA FUNÇÃO: Extrai caminho relativo limpo
func cleanFilePath(fullPath string) string {
	// Remove prefixos comuns de caminhos absolutos
	prefixes := []string{
		"/home/",
		"/Users/",
		"C:/",
		"/root/",
	}

	clean := fullPath
	for _, prefix := range prefixes {
		if strings.HasPrefix(clean, prefix) {
			// Remove o prefixo e o próximo nível (nome do usuário)
			parts := strings.SplitN(strings.TrimPrefix(clean, prefix), "/", 3)
			if len(parts) >= 3 {
				clean = parts[2] // Pega só a parte relevante
			} else if len(parts) == 2 {
				clean = parts[1]
			}
			break
		}
	}

	// Garante que não começa com /
	clean = strings.TrimPrefix(clean, "/")

	// Se ainda estiver vazio, usa o nome do arquivo
	if clean == "" {
		clean = filepath.Base(fullPath)
	}

	return clean
}

func receiveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao ler dados: %v", err), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Gera ID único para este backup
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	repoID := timestamp

	// Cria diretório para este repositório
	repoDir := filepath.Join("repos", repoID)
	err = os.MkdirAll(repoDir, 0755)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao criar diretório: %v", err), http.StatusInternalServerError)
		return
	}

	// Parse dos dados
	var data ReceivedData
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao parsear JSON: %v", err), http.StatusBadRequest)
		return
	}

	// ✅ SALVA OS ARQUIVOS COM CAMINHOS RELATIVOS
	filesSaved := 0
	for _, commit := range data.Commits {
		for _, file := range commit.Files {
			// ✅ LIMPA O CAMINHO para não criar /home/arthur/...
			relativePath := cleanFilePath(file.Path)

			// Cria subdiretórios se necessário (agora com caminho limpo)
			filePath := filepath.Join(repoDir, relativePath)
			dir := filepath.Dir(filePath)
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Printf("⚠️  Erro ao criar dir %s: %v\n", dir, err)
				continue
			}

			// ✅ DECODIFICA E SALVA O CONTEÚDO
			content, err := base64.StdEncoding.DecodeString(file.Content)
			if err != nil {
				fmt.Printf("⚠️  Erro ao decodificar %s: %v\n", file.Path, err)
				continue
			}

			err = os.WriteFile(filePath, content, 0644)
			if err != nil {
				fmt.Printf("⚠️  Erro ao salvar %s: %v\n", relativePath, err)
				continue
			}

			filesSaved++
			fmt.Printf("📁 Salvo: %s (%d bytes)\n", relativePath, len(content))
		}
	}

	// Salva metadata do backup
	metadataPath := filepath.Join(repoDir, "metadata.json")
	err = os.WriteFile(metadataPath, body, 0644)
	if err != nil {
		fmt.Printf("⚠️  Erro ao salvar metadata: %v\n", err)
	}

	// Salva backup JSON também (para histórico)
	backupPath := fmt.Sprintf("backup_%s.json", timestamp)
	err = os.WriteFile(backupPath, body, 0644)
	if err != nil {
		fmt.Printf("⚠️  Erro ao salvar backup: %v\n", err)
	}

	fmt.Printf("📦 Recebido: Repo %s, %d commits, %d arquivos salvos\n",
		truncateID(data.RepositoryID), len(data.Commits), filesSaved)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Recebido: %d arquivos salvos em %s", filesSaved, repoDir)
}

func truncateID(id string) string {
	if len(id) >= 8 {
		return id[:8]
	}
	return id
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OpenV Server OK - %s", time.Now().Format(time.RFC3339))
}

func Start() {
	http.HandleFunc("/receive", receiveHandler)
	http.HandleFunc("/health", healthHandler)

	port := "0.0.0.0:8080"

	fmt.Println("🚀 OpenV Server rodando em http://localhost" + port)
	fmt.Println("📡 Endpoint de push: http://localhost" + port + "/receive")
	fmt.Println("📁 Arquivos salvos em: ./repos/ (com caminhos relativos)")
	fmt.Println("🛑 Pressione Ctrl+C para parar")

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("❌ Erro ao iniciar servidor: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	Start()
}
