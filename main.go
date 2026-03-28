package main

import (
	"fmt"
	"os"

	"github.com/arthur/openv/cmd"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		cmd.Init()
	case "commit":
		cmd.Commit(os.Args[2:])
	case "status":
		cmd.Status()
	case "log":
		cmd.Log()
	case "push":
		cmd.Push(os.Args[2:])
	case "restore": // ✅ NOVO
		cmd.Restore(os.Args[2:])
	case "--commit-on-m":
		cmd.CommitOnModified(os.Args[2:])
	case "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Comando desconhecido: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`🚀 OpenV - Sistema de Versionamento Leve

Uso:
  openv init                    Inicializa repositório OpenV
  openv commit "mensagem"       Commita arquivos staged
  openv --commit-on-m "msg" .   Commita apenas arquivos modificados
  openv status                  Mostra status dos arquivos
  openv log                     Mostra histórico com diffs
  openv restore <commit> <file> Restaura arquivo de commit antigo
  openv push <server_url>       Envia commits para servidor
  openv -h, --help              Mostra esta ajuda

Exemplos:
  openv init
  openv --commit-on-m "Primeiro commit" .
  openv log
  openv restore 0 teste.txt     # Restaura do commit 0
  openv push http://192.168.1.100:8080/receive`)
}
