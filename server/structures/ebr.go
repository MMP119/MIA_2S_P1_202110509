package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

//structura del EBR

type EBR struct {
	Part_mount	[1]byte
	Part_fit	[1]byte
	Part_start	int32
	Part_size	int32
	Part_next	int32
	Part_name	[16]byte
}

func CreateEBR(path string, size int32, fdisk *FDISK, startEBR  int32) (string, error) {

    var newEBR EBR

    // Establecer valores del EBR
    newEBR.Part_mount[0] = '0'                              // Estado: inactivo (no montado)
    newEBR.Part_fit[0] =fdisk.Fit[0]                        // Fit heredado de la partición extendida
    newEBR.Part_start = startEBR                            // Inicio de la partición lógica
    newEBR.Part_size = size                                 // Tamaño de la partición lógica
    newEBR.Part_next = -1                                   // No hay siguiente EBR al principio
    copy(newEBR.Part_name[:], []byte(fdisk.Name))           // Convert the string to a [16]byte array and assign it to Part_name


    //imprimir la información del EBR
    fmt.Println("\n\n\nEBR creado exitosamente")
    fmt.Printf("  Estado: %c\n", newEBR.Part_mount[0])
    fmt.Printf("  Ajuste: %c\n", newEBR.Part_fit[0])
    fmt.Printf("  Inicio: %d\n", newEBR.Part_start)
    fmt.Printf("  Tamaño: %d\n", newEBR.Part_size)
    fmt.Printf("  Nombre: %s\n", string(newEBR.Part_name[:]))
    fmt.Printf("  Siguiente EBR: %d\n", newEBR.Part_next)
    fmt.Printf("\nFIN\n")
    
    return "EBR creado exitosamente", nil
}

// Serializa el EBR y lo escribe en el disco en la posición dada
func (ebr *EBR) SerializeEBR(path string, position int32) (string, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return "Error al abrir el archivo al serializar MBR",err
	}
	defer file.Close()

    _, err = file.Seek(int64(position), 0) // Moverse a la posición del EBR
    if err != nil {
        return "Error al moverse a la posición del EBR", fmt.Errorf("error al moverse a la posición: %s", err)
    }

	err = binary.Write(file, binary.LittleEndian, ebr)
	if err != nil {
		return "Error al escribir en el archivo al serializar MBR",err
	}

	return "",nil
}

// Deserializa el EBR desde el disco en la posición dada, si retorna un error es porque no existe un EBR en esa posición
func (ebr *EBR) DeserializeEBR(path string, position int32) (string, error) {
    file, err := os.Open(path)
    if err != nil {
        return "Error abriendo el archivo para deserialización del EBR", fmt.Errorf("error abriendo el archivo: %s", err)
    }
    defer file.Close()

    // Moverse a la posición del EBR
    _, err = file.Seek(int64(position), 0) // Moverse a la posición del EBR
    if err != nil {
        return "Error al moverse a la posición del EBR", fmt.Errorf("error al moverse a la posición: %s", err)
    }

    // Tamaño de la estructura EBR
    ebrSize := binary.Size(ebr)
    if ebrSize <= 0 {
        return "Tamaño inválido para el EBR", fmt.Errorf("tamaño inválido para el EBR: %d", ebrSize)
    }

    // Crear un buffer para leer el archivo
    buffer := make([]byte, ebrSize)
    n, err := file.Read(buffer)
    if err != nil {
        return "Error al leer el archivo al deserializar el EBR", fmt.Errorf("error al leer el archivo: %s", err)
    }

    // Verificar si se leyeron menos bytes de los esperados
    if n != ebrSize {
        return "Error: cantidad de bytes leídos no coincide con el tamaño del EBR", fmt.Errorf("se leyeron %d bytes, pero se esperaban %d", n, ebrSize)
    }

    reader := bytes.NewReader(buffer)
    err = binary.Read(reader, binary.LittleEndian, ebr)
    if err != nil {
        return "Error al leer el buffer al deserializar el EBR", fmt.Errorf("error al leer el buffer: %s", err)
    }

    // Validar si el EBR en esta posición es válido (e.g., Part_size > 0)
    if ebr.Part_size <= 0 {
        return "No existe un EBR válido en esta posición", fmt.Errorf("no existe un EBR válido en esta posición")
    }

    // Validar el Part_start y Part_next para asegurarse de que tienen valores lógicos
    if ebr.Part_start <= 0 || ebr.Part_next < -1 {
        return "Error: el EBR deserializado tiene valores inválidos", fmt.Errorf("datos inválidos en el EBR")
    }

    fmt.Printf("EBR deserializado exitosamente:\n  Estado: %c\n  Ajuste: %c\n  Inicio: %d\n  Tamaño: %d\n  Nombre: %s\n  Siguiente EBR: %d\n",
        ebr.Part_mount[0], ebr.Part_fit[0], ebr.Part_start, ebr.Part_size, string(ebr.Part_name[:]), ebr.Part_next)

    return "", nil
}