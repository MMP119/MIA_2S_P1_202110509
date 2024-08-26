package structures

import (
	//"encoding/binary"
	"fmt"
	util "server/util"
)

// FDISK estructura que representa el comando fdisk con sus parámetros
type FDISK struct {
	Size int    // Tamaño de la partición
	Unit string // Unidad de medida del tamaño (B, K o M); por defecto K
	Path string // Ruta del archivo del disco
	TypE  string // Tipo de partición (P, E, L)
	Fit  string // Tipo de ajuste (BF, FF, WF); por defecto WF
	Name string // Nombre de la partición
}

func CommandFdisk(fdisk *FDISK) (string, error) {
	
	// Convertir el tamaño a bytes
	sizeBytes, err := util.ConvertToBytes(fdisk.Size, fdisk.Unit)
	if err != nil {
		fmt.Println("Error converting size:", err)
		return "Error converting size en Fdisk", err
	}

	var msg string

	// para crear la particion primaria
	if(fdisk.TypE == "P"){
		msg, err = CreatePrimaryPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creating primary partition:", err)
			return msg, err
		}
	}else if(fdisk.TypE == "E"){

		msg, err = CreatePrimaryPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creating extended partition:", err)
			return msg, err
		}


	}else if(fdisk.TypE == "L"){

		msg, err = CreateLogicalPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creating logical partition:", err)
			return msg, err
		}

	}
	return "",nil
}


func CreatePrimaryPartition(fdisk *FDISK, sizeBytes int)(string, error){
	
	var mbr MBR
	
	msg, err := mbr.DeserializeMBR(fdisk.Path)
	if err != nil {
		return msg, fmt.Errorf("error leyendo el MBR del disco: %s", err)
	}

	// Contar el número de particiones primarias y extendidas
	primaryCount := 0
	extendedExists := 0
	for _, partition := range mbr.Mbr_partitions {
		if partition.Part_status[0] != '2' {
			if partition.Part_type[0] == 'P' {
				primaryCount++
			} else if partition.Part_type[0] == 'E' {
				extendedExists++
			}
		}
	}

	// Verificar que no se exceda el límite de 4 particiones
	if primaryCount >4 {
		return "ERROR: No se pueden crear más particiones primarias", fmt.Errorf("límite de particiones primarias alcanzado")
	}

	if extendedExists >=2 {
		return "ERROR: No se pueden crear más particiones extendidas, ya existe una en el disco", fmt.Errorf("ya existe una partición extendida")
	}

	// se obtiene la primera particion libre
	particionDisponible, inicioParticion, indexParticion, msg:= mbr.GetFirstPartitionAvailable()
	if particionDisponible == nil {
		return msg, fmt.Errorf("no hay particiones disponibles")
	}

	// crear la particion con los parámetros proporcionados 
	particionDisponible.CreatePartition(inicioParticion, sizeBytes, fdisk.TypE, fdisk.Fit, fdisk.Name)

	// montar la particion
	mbr.Mbr_partitions[indexParticion] = *particionDisponible //asignar la particion al MBR

	// Serialiazar el MBR modificado
	msg, err = mbr.SerializeMBR(fdisk.Path)
	if err != nil {
		return msg, fmt.Errorf("error escribiendo el MBR al disco: %s", err)
	}	
	return "",nil
}


func CreateLogicalPartition(fdisk *FDISK, sizeBytes int) (string, error) {
    var mbr MBR
    // Deserializar el MBR del disco
    msg, err := mbr.DeserializeMBR(fdisk.Path)
    if err != nil {
        return msg, fmt.Errorf("error leyendo el MBR del disco: %s", err)
    }

    // Verificar que exista una partición extendida
    var extendedPartition *PARTITION
    for _, partition := range mbr.Mbr_partitions {
        if partition.Part_type[0] == 'E' {
            extendedPartition = &partition
            break
        }
    }
    if extendedPartition == nil {
        return "ERROR: no existe una partición extendida", fmt.Errorf("no existe una partición extendida")
    }

    // Intentar deserializar el primer EBR en la partición extendida
    var ebr EBR
    ebrPosition := extendedPartition.Part_start
    _, err = ebr.DeserializeEBR(fdisk.Path, ebrPosition)
    
    if err != nil || string(ebr.Part_name[:]) == ""|| ebr.Part_size <= 0 || ebr.Part_next < -1 || len(ebr.Part_name) == 0 || ebr.Part_name[0] == 0{
        // Si no existe EBR válido, crear el primer EBR
        fmt.Println("No existe un EBR válido en la partición extendida, creando el primer EBR.")
        
        // Crear el primer EBR
        _, err = CreateEBR(fdisk.Path, int32(sizeBytes), fdisk, ebrPosition)
        if err != nil {
            return "ERROR al crear el primer EBR", err
        }

        fmt.Println("Primer EBR creado exitosamente.")
        return "Primer EBR creado exitosamente", nil
    }

    // Si ya existe un EBR válido, buscar el último EBR
    fmt.Println("Ya existe un EBR, buscando el último EBR en la lista.")
    var lastEBR *EBR
    for {
        if ebr.Part_next == -1 {
            lastEBR = &ebr
            break
        }
        _, err = ebr.DeserializeEBR(fdisk.Path, ebr.Part_next)
        if err != nil {
            return "ERROR: no se pudo leer el siguiente EBR", fmt.Errorf("no se pudo leer el siguiente EBR")
        }
    }

    // Si ya se encontró el último EBR, se puede crear el siguiente EBR en la lista
    fmt.Println("Último EBR encontrado, creando el siguiente EBR.")
    newEBRStart := lastEBR.Part_start + lastEBR.Part_size // Calcular dónde empieza el nuevo EBR
    lastEBR.Part_next = newEBRStart // Actualizar el Part_next del último EBR

    // Serializar el nuevo EBR en el disco
    _, err = CreateEBR(fdisk.Path, int32(sizeBytes), fdisk, newEBRStart)
    if err != nil {
        return "ERROR al crear el nuevo EBR", err
    }

    fmt.Println("Nuevo EBR creado exitosamente.")
    return "Nuevo EBR creado exitosamente", nil
}


    // var ebr EBR
    // ebrPosition := extendedPartition.Part_start
    // _, err = ebr.DeserializeEBR(fdisk.Path, ebrPosition)
    // if err == nil {
    //     // Si no hay EBR en la partición extendida, creamos el primer EBR
    //     msg, err = CreateEBR(fdisk.Path, int32(sizeBytes), fdisk, ebrPosition)
    //     if err != nil {
    //         return msg, fmt.Errorf("error creando el EBR: %s", err)
    //     }
	// 	_, err = mbr.SerializeMBR(fdisk.Path)
	// 	if err != nil {
	// 		return "Error escribiendo el MBR al disco", fmt.Errorf("error escribiendo el MBR al disco: %s", err)
	// 	}

    //     // Crear la partición lógica justo después del EBR
    //     logicalPartitionStart := extendedPartition.Part_start + int32(binary.Size(ebr)) // EBR ocupa un espacio
    //     var logicalPartition PARTITION
    //     logicalPartition.CreatePartition(int(logicalPartitionStart), sizeBytes, "L", fdisk.Fit, fdisk.Name)
    // }
