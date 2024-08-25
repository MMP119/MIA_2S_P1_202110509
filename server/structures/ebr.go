package structures


import(
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

// Crea un nuevo EBR y lo serializa en la partición extendida
func CreateEBR(path string, size int32, extendedPartition *PARTITION, fdiskName string) (string, error) {
    var newEBR EBR

    // Establecer valores del EBR
    newEBR.Part_mount[0] = '0'         // Estado: inactivo (no montado)
    newEBR.Part_fit[0] = extendedPartition.Part_fit[0] // Fit heredado de la partición extendida
    newEBR.Part_size = size             // Tamaño de la partición lógica
    newEBR.Part_next = -1               // No hay siguiente EBR al principio

    // Asignar el nombre (en blanco por defecto)
    copy(newEBR.Part_name[:], fdiskName)

    // Calcular la posición del EBR dentro de la partición extendida
    ebrPosition := extendedPartition.Part_start
    msg, err := newEBR.SerializeEBR(path, ebrPosition)
    if err != nil {
        return msg, fmt.Errorf("error escribiendo el nuevo EBR: %v", err)
    }

    return "EBR creado exitosamente", nil
}


// Serializa el EBR y lo escribe en el disco en la posición dada
func (ebr *EBR) SerializeEBR(path string, position int32) (string, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return "error abriendo el archivo para serializacion el ebr",fmt.Errorf("error abriendo el archivo: %s", err)
	}
	defer file.Close()

	file.Seek(int64(position), 0)
	err = binary.Write(file, binary.BigEndian, ebr)
	if err != nil {
		return "error escribiendo el EBR",fmt.Errorf("error escribiendo el EBR: %s", err)
	}
	return "", nil
}

// Deserializa el EBR desde el disco en la posición dada
func (ebr *EBR) DeserializeEBR(path string, position int32) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "Error abriendo el archivo para deserialización del EBR", fmt.Errorf("error abriendo el archivo: %s", err)
	}
	defer file.Close()

	file.Seek(int64(position), 0)// Mover el puntero al inicio de la partición extendida
	err = binary.Read(file, binary.BigEndian, ebr)// Leer el EBR
	if err != nil {
		return "Error leyendo el EBR", fmt.Errorf("error leyendo el EBR: %s", err)
	}

	// Verificar si el EBR es válido (por ejemplo, si Part_size es mayor que 0)
	if ebr.Part_size <= 0 {
		return "EBR inválido o no encontrado", fmt.Errorf("EBR inválido o no encontrado")
	}

	return "", nil
}

// Imprime la información de todos los EBRs y sus particiones lógicas asociadas
func (ebr *EBR)PrintEBRsAndLogicalPartitions(path string, extendedPartition *PARTITION) error {
    var currentEBR EBR
    var logicalPartition PARTITION
    currentEBRPosition := extendedPartition.Part_start

    file, err := os.OpenFile(path, os.O_RDONLY, 0644)
    if err != nil {
        return fmt.Errorf("error abriendo el archivo: %v", err)
    }
    defer file.Close()

    for {
        // Moverse a la posición del EBR actual
        file.Seek(int64(currentEBRPosition), 0)

        // Leer el EBR en la posición actual
        err = binary.Read(file, binary.BigEndian, &currentEBR)
        if err != nil {
            return fmt.Errorf("error leyendo el EBR: %v", err)
        }

        // Imprimir la información del EBR
        fmt.Printf("EBR en la posición %d:\n", currentEBRPosition)
        fmt.Printf("  Estado: %c\n", currentEBR.Part_mount[0])
        fmt.Printf("  Ajuste: %c\n", currentEBR.Part_fit[0])
        fmt.Printf("  Inicio: %d\n", currentEBR.Part_start)
        fmt.Printf("  Tamaño: %d\n", currentEBR.Part_size)
        fmt.Printf("  Nombre: %s\n", string(currentEBR.Part_name[:]))
        fmt.Printf("  Siguiente EBR: %d\n", currentEBR.Part_next)

        // Leer la partición lógica asociada al EBR (ubicada en Part_start del EBR)
        file.Seek(int64(currentEBR.Part_start), 0)
        err = binary.Read(file, binary.BigEndian, &logicalPartition)
        if err != nil {
            return fmt.Errorf("error leyendo la partición lógica: %v", err)
        }

        // Imprimir la información completa de la partición lógica
        fmt.Println("Información de la partición lógica asociada:")
        logicalPartition.Print()

        // Si no hay más EBRs, terminamos el ciclo
        if currentEBR.Part_next == -1 {
            break
        }

        // Moverse al siguiente EBR
        currentEBRPosition = currentEBR.Part_next
    }

    return nil
}