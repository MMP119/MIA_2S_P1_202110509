package structures

import (
	"fmt"
	"os"
	"path/filepath"
	util "server/util"
)

// MKDISK estructura que representa el comando mkdisk con sus parámetros
type MKDISK struct {
	Size int    // Tamaño del disco
	Unit string // Unidad de medida del tamaño (K o M)
	Fit  string // Tipo de ajuste (BF, FF, WF)
	Path string // Ruta del archivo del disco
}


func CommandMkdisk(mkdisk *MKDISK) error {
	// Convertir el tamaño a bytes
	sizeBytes, err := util.ConvertToBytes(mkdisk.Size, mkdisk.Unit)
	if err != nil {
		fmt.Println("Error converting size:", err)
		return err
	}

	// Crear el disco con el tamaño proporcionado
	err = CreateDisk(mkdisk, sizeBytes)
	if err != nil {
		fmt.Println("Error creating disk:", err)
		return err
	}

	// Crear el MBR con el tamaño proporcionado
	err = CreateMBR(mkdisk, sizeBytes)
	if err != nil {
		fmt.Println("Error creating MBR:", err)
		return err
	}

	return nil
}

func CreateDisk(mkdisk *MKDISK, sizeBytes int) error {
	// Crear las carpetas necesarias
	err := os.MkdirAll(filepath.Dir(mkdisk.Path), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return err
	}

	// Crear el archivo binario
	file, err := os.Create(mkdisk.Path)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	// Escribir en el archivo usando un buffer de 1 MB
	buffer := make([]byte, 1024*1024) // Crea un buffer de 1 MB
	for sizeBytes > 0 {
		writeSize := len(buffer)
		if sizeBytes < writeSize {
			writeSize = sizeBytes // Ajusta el tamaño de escritura si es menor que el buffer
		}
		if _, err := file.Write(buffer[:writeSize]); err != nil {
			return err // Devuelve un error si la escritura falla
		}
		sizeBytes -= writeSize // Resta el tamaño escrito del tamaño total
	}
	return nil
}