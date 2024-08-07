package structures

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

//ruta donde se creará el disco
//const fullPath = "/home/mario/Escritorio/GitHub/MIA_2S_P1_202110509/server/disco.mia"

//CreateBinaryFile crea un archivo binario, del tamaño y unidad especificados y se escribe un MBR en el archivo
func CreateBinaryFile(size int, fit string, unit string, path string) error {

	// Convertir el tamaño a bytes
	sizeInBytes, err := convertToBytes(size, unit)
	if err != nil {
		return err
	}

	// Obtener el directorio de la ruta
	dir := filepath.Dir(path)

	// Verificar si el directorio existe, si no, crear las carpetas necesarias
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error: no se pudo crear el directorio: '%s'", err) // Si no se pudo crear el directorio, retornar el error
		}
	}

	// Crear el archivo en la ruta especificada
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error: no se pudo crear el disco: '%s'", err) // Si no se pudo crear el archivo, retornar el error
	}
	defer file.Close() // Cerrar el archivo

	//Crear y escribir el MBR en el archivo
	mbr := CrearMBR(int32(sizeInBytes), fit)
	err = writeMBRToFile(file, mbr)
	if err != nil {
		return err // Si no se pudo escribir el MBR, retornar el error
	}
	
	// Escribir en el archivo
	err = writeToFile(file, sizeInBytes-int(binary.Size(mbr)))
	if err != nil {
		return err
	}

	//fmt.Println("Disco creado exitosamente")
	return nil
}


// writeMBRToFile escribe el MBR en el archivo binario
func writeMBRToFile(file *os.File, mbr MBR) error {
	err := binary.Write(file, binary.LittleEndian, &mbr)
	if err != nil {
		return err // Si no se pudo escribir el MBR, retornar el error
	}
	return nil
}

//convertToBytes convierte el tamaño del disco y la unidad a bytes

func convertToBytes(size int, unit string) (int, error) {

	switch unit {
	case "K":
		return size * 1024, nil
	case "M":
		return size * 1024 * 1024, nil
	default:
		return 0, errors.New("error: la unidad de medida debe ser K o M")
	}
}

//writeToFile escribe en el archivo binario los bytes especificados

func writeToFile(file *os.File, sizeInBytes int) error {

	buffer := make([]byte, 1024*1024) // buffer de 1MB

	for sizeInBytes > 0 {

		writeSize := len(buffer)
		if sizeInBytes < writeSize{
			writeSize = sizeInBytes //ajusta el tamaño de escritura si es menor al buffer
		}

		if _, err := file.Write(buffer[:writeSize]); err != nil { // el :writeSize es para que solo escriba la cantidad de bytes que se necesita
			return err //retorna un error si no se pudo escribir en el archivo
		}
		sizeInBytes -= writeSize //resta el tamaño de escritura al tamaño total
	}
	fmt.Println("Disco creado exitosamente")
	return nil

}


//DeleteBinaryFile elimina un archivo binario
func DeleteBinaryFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("error: no se pudo eliminar el disco: '%s'", err)
	}
	fmt.Println("Disco eliminado exitosamente")
	return nil
}