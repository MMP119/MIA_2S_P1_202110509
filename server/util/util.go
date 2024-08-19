package utils

import (
	"errors"
	"fmt"
	"os"
)

func ConvertToBytes(size int, unit string) (int, error) {
	switch unit {
		case "B":
			return size, nil // Devuelve el tamaño en bytes
		case "K":
			return size * 1024, nil // Convierte kilobytes a bytes
		case "M":
			return size * 1024 * 1024, nil // Convierte megabytes a bytes
		default:
			return 0, errors.New("invalid unit") // Devuelve un error si la unidad es inválida
	}
}


//DeleteBinaryFile elimina un archivo binario
func DeleteBinaryFile(path string) (string, error) {
	err := os.Remove(path)
	if err != nil {
		return "", fmt.Errorf("error: no se pudo eliminar el disco: '%s'", err)
	}
	//fmt.Println("Disco eliminado exitosamente")
	return "Disco eliminado exitosamente", nil
}


// ConvertToFixedSizeArray convierte un string en un array de tamaño fijo
func ConvertToFixedSizeArray(input string, size int) [16]byte {
	var array [16]byte
	copy(array[:], input)
	return array
}

const Carnet string = "09"//202110509
var Alfabeto = []string {
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J","K", "L", "M", "N", 
	"O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
}

//map para almacenar la asignacion de letras a los path
var pathToLetter = make(map[string]string)

//indice para la siguiente letra disponible
var nextLetterIndex = 0

// GetLetter obtiene la letra asignada a un path
func GetLetter(path string)(string, error){
	//asignar letra si el path no tiene una asignada
	if _, exist := pathToLetter[path]; !exist{
		if nextLetterIndex < len(Alfabeto){
			pathToLetter[path] = Alfabeto[nextLetterIndex]
			nextLetterIndex++
		}else{
			return "No hay letras disponibles", errors.New("no hay letras disponibles")
		}
	}
	return pathToLetter[path], nil
}