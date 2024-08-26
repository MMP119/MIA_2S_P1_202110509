package structures

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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
    

    // Asignar el nombre (en blanco por defecto)
    copy(newEBR.Part_name[:], fdisk.Name)

    // Calcular la posición del EBR dentro de la partición extendida
    msg, err := newEBR.SerializeEBR(path, startEBR)
    if err != nil {
        return msg, fmt.Errorf("error escribiendo el nuevo EBR: %v", err)
    }

    return "EBR creado exitosamente", nil
}

// Serializa el EBR y lo escribe en el disco en la posición dada
func (ebr *EBR) SerializeEBR(path string, position int32) (string, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return "Error al abrir el archivo al serializar MBR",err
	}
	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, ebr)
	if err != nil {
		return "Error al escribir en el archivo al serializar MBR",err
	}

	return "",nil
}

// Deserializa el EBR desde el disco en la posición dada
func (ebr *EBR) DeserializeEBR(path string, position int32) (string, error) {
    file, err := os.Open(path)
    if err != nil {
        return "Error abriendo el archivo para deserialización del EBR", fmt.Errorf("error abriendo el archivo: %s", err)
    }
    defer file.Close()

    // Moverse a la posición del EBR
    _, err = file.Seek(int64(position), 0)
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
    bytesRead, err := file.Read(buffer)
    if err != nil {
        if err == io.EOF {
            return "Fin del archivo (EOF) alcanzado al leer el EBR", io.EOF
        }
        return "Error al leer el archivo al deserializar el EBR", fmt.Errorf("error leyendo el archivo: %s", err)
    }

    // Verificar que se haya leído el tamaño completo
    if bytesRead != ebrSize {
        return "Error: lectura incompleta del EBR", fmt.Errorf("lectura incompleta, se esperaban %d bytes, se leyeron %d", ebrSize, bytesRead)
    }

    // Convertir los bytes en la estructura EBR
    bufferReader := bytes.NewReader(buffer)
    err = binary.Read(bufferReader, binary.LittleEndian, ebr)
    if err != nil {
        return "Error al deserializar el EBR desde los bytes leídos", fmt.Errorf("error deserializando: %s", err)
    }

    // Aquí agregamos la validación: verificar si el EBR es "vacío"
    if ebr.Part_size == 0 && ebr.Part_next == -1 && len(ebr.Part_name) == 0 {
        return "No se encontró un EBR válido en esta posición", fmt.Errorf("EBR no inicializado")
    }

    return "", nil
}
