package main

import (
	"bufio"
	"fmt"
	"os"
	analyzer "server/analyzer"
	structures "server/structures"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		// Si hay argumentos, procesar el archivo de carga masiva
		filePath := os.Args[1]
		err := processBatchFile(filePath)
		if err != nil {
			fmt.Println("Error:", err)
		}
	} else {
		// Si no hay argumentos, iniciar el modo interactivo
		interactiveMode()
	}
}

// Modo interactivo para comandos individuales
func interactiveMode() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">>> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		_, err := analyzer.Analyzer(input)
		if err != nil {
			fmt.Println("Error:", err)
		}else{
			verifyMBR("/home/mario/Escritorio/GitHub/MIA_2S_P1_202110509/server/disk.mia")
			// mkdisk -size=1 -unit=M -fit=WF -path=/home/mario/Escritorio/GitHub/MIA_2S_P1_202110509/server/disk.mia
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer:", err)
	}
	
}

// Procesar archivo de carga masiva
func processBatchFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") { // Ignorar líneas vacías y comentarios
			fmt.Println(line)
			continue
		}
		_, err := analyzer.Analyzer(line)
		if err != nil {
			fmt.Printf("Error al procesar comando '%s': %v\n", line, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error al leer el archivo: %v", err)
	}
	return nil
}


func verifyMBR(filePath string) {
	mbr, err := structures.ReadMBRFromFile(filePath)
	if err != nil {
		fmt.Println("Error al leer el MBR:", err)
		return
	}
	structures.PrintMBR(mbr)
}