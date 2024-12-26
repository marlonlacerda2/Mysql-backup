package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Função para executar um comando no sistema operacional
func runCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	// Configurações de conexão com o MySQL
	user := ""     // Usuário do MySQL
	password := "" // Senha do MySQL
	host := ""     // Endereço do servidor MySQL
	port := ""     // Porta do servidor MySQL

	listDatabasesCmd := fmt.Sprintf("mysql -u%s -p%s -h%s -P%s -e 'SHOW DATABASES;' --skip-column-names", user, password, host, port)
	out, err := exec.Command("bash", "-c", listDatabasesCmd).Output()
	if err != nil {
		log.Fatalf("Erro ao listar os bancos de dados: %v", err)
	}

	databases := strings.Split(string(out), "\n")
	for _, db := range databases {
		if db == "information_schema" || db == "performance_schema" || db == "mysql" || db == "sys" {
			continue
		}

		backupFile := fmt.Sprintf("%s_backup.sql", db)

		mysqldumpCmd := fmt.Sprintf("mysqldump -u%s -p%s -h%s -P%s %s > %s", user, password, host, port, db, backupFile)

		fmt.Printf("Fazendo backup do banco de dados %s...\n", db)
		err = runCommand("bash", []string{"-c", mysqldumpCmd})
		if err != nil {
			log.Printf("Erro ao fazer o backup do banco %s: %v", db, err)
		} else {
			fmt.Printf("Backup do banco %s concluído. Arquivo salvo em %s\n", db, backupFile)
		}
	}
}
