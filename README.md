# OpenV

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)
[![Platforms](https://img.shields.io/badge/Platforms-Linux%20%7C%20Windows%20%7C%20macOS-brightgreen)](bin/)

Sistema de versionamento leve e rápido em Go.
Os binários do OpenV estão na pasta `bin/`, porém se quiser fazer uma modificação ou compilar, o tutorial está abaixo.

## Instalação

```bash
go build -o openv main.go cmd/*.go core/*.go
sudo cp openv /usr/local/bin/
```

## Uso

```bash
openv init
openv --commit-on-m "Primeiro commit" .
openv status
openv log
openv push http://192.168.1.100:8080/receive
```

## Tabela de comandos
|Comando|Função|Exemplo|
|-------|------|-------|
|`openv init`|Inicializa repositório OpenV|`openv init`|
|`openv commit "mensagem"`|Commita arquivos staged|`openv commit "Commit em arquivos staged"`|
|`openv --commit-on-m "msg" .`|Commita apenas arquivos modificados|`openv --commit-on-m "Commit em arquivos modificados" .`|
|`openv log`|Mostra histórico com diffs|`openv log`|
|`openv restore <commit> <file>`|Restaura arquivo de commit antigo|`openv restore 0 teste.txt`|
|`openv push <server_url>`|Envia commits para servidor|`openv push http://192.168.1.100:8080/receive`|
|`openv -h` ou `openv --help`|Mostra tela de ajuda|`openv -h`|

## Servidor (Minicomputador ou servidor)

```bash
go build -o openv-server server/receiver.go
./openv-server
```

## Estrutura .openv

Arquivo JSON com:
- repository_id: ID único do repositório
- version: Versão do OpenV
- commits: Lista de commits com arquivos modificados

**Eu tenho 9 anos, então erros podem acontecer!** caso encontre um erro, faça um issue no repositório oficial!

> Este projeto está licenciado sobre a licença Apache 2.0. Criado por Arthur.

## Espero que goste!