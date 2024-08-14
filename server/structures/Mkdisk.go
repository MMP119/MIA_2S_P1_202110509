package structures

import (
	"fmt"
	"os"
	"path/filepath"
	util "server/util"
)

type MKDISK struct {
	Size int    
	Unit string
	Fit  string 
	Path string 
}


func CommandMkdisk(mkdisk *MKDISK) error {

	sizeBytes, err := util.ConvertToBytes(mkdisk.Size, mkdisk.Unit)
	if err != nil {
		fmt.Println("Error converting size:", err)
		return err
	}


	err = CreateDisk(mkdisk, sizeBytes)
	if err != nil {
		fmt.Println("Error creating disk:", err)
		return err
	}

	err = CreateMBR(mkdisk, sizeBytes)
	if err != nil {
		fmt.Println("Error creating MBR:", err)
		return err
	}

	return nil
}

func CreateDisk(mkdisk *MKDISK, sizeBytes int) error {

	err := os.MkdirAll(filepath.Dir(mkdisk.Path), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return err
	}

	file, err := os.Create(mkdisk.Path)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

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