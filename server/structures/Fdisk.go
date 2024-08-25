package structures

import (
	"encoding/binary"
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

		msg, err = CreateExtendedPartition(fdisk, sizeBytes)
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
	extendedExists := false
	for _, partition := range mbr.Mbr_partitions {
		if partition.Part_status[0] != '2' {
			if partition.Part_type[0] == 'P' {
				primaryCount++
			} else if partition.Part_type[0] == 'E' {
				extendedExists = true
			}
		}
	}

	// Verificar que no se exceda el límite de 4 particiones
	if primaryCount >= 3 && !extendedExists {
		return "ERROR: No se pueden crear más de 3 particiones primarias si no hay una extendida", fmt.Errorf("límite de particiones primarias alcanzado")
	}

	if primaryCount >= 4 {
		return "ERROR: No se pueden crear más particiones primarias", fmt.Errorf("límite de particiones primarias alcanzado")
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


func CreateExtendedPartition(fdisk *FDISK, sizeBytes int)(string, error){
	
	var mbr MBR

	msg, err := mbr.DeserializeMBR(fdisk.Path)
	if err != nil {
		return msg, fmt.Errorf("error leyendo el MBR del disco: %s", err)
	}

	// verificar si ya existe una particion extendida
	for _, partition := range mbr.Mbr_partitions {
		if partition.Part_type[0] != '0' && partition.Part_type[0] == 'E' {
			return "ERROR: ya existe una particion extendida",fmt.Errorf("ya existe una partición extendida")
		}
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

    // Deserializar el primer EBR en la partición extendida
    var ebr EBR
    ebrPosition := extendedPartition.Part_start
    _, err = ebr.DeserializeEBR(fdisk.Path, ebrPosition)
    if err != nil {
        // Si no hay EBR en la partición extendida, creamos el primer EBR
        msg, err = CreateEBR(fdisk.Path, int32(sizeBytes), extendedPartition, fdisk.Name)
        if err != nil {
            return msg, fmt.Errorf("error creando el EBR: %s", err)
        }

        // Crear la partición lógica justo después del EBR
        logicalPartitionStart := extendedPartition.Part_start + int32(binary.Size(ebr)) // EBR ocupa un espacio
        var logicalPartition PARTITION
        logicalPartition.CreatePartition(int(logicalPartitionStart), sizeBytes, "L", fdisk.Fit, fdisk.Name)
        
        // Serializar la nueva partición lógica
        msg, err = logicalPartition.SerializePartition(fdisk.Path, logicalPartitionStart)
        if err != nil {
            return msg, fmt.Errorf("error escribiendo la partición lógica: %s", err)
        }

        return "Partición lógica creada exitosamente", nil
    }

    // Si ya existen EBRs, avanzar hasta el último EBR
    var lastEBR *EBR
    for {
        if ebr.Part_next == -1 {
            lastEBR = &ebr
            break
        }
        _, err = ebr.DeserializeEBR(fdisk.Path, ebr.Part_next)
        if err != nil {
            return "ERROR: no se pudo leer el último EBR", fmt.Errorf("no se pudo leer el último EBR")
        }
    }

    // Verificar si hay espacio para una nueva partición lógica
    newEBRStart := lastEBR.Part_start + lastEBR.Part_size
    if newEBRStart+int32(sizeBytes) > extendedPartition.Part_start+extendedPartition.Part_size {
        return "ERROR: no hay espacio para una nueva partición lógica", fmt.Errorf("no hay espacio suficiente")
    }

    // Crear un nuevo EBR en el espacio disponible
    lastEBR.Part_next = newEBRStart
    msg, err = CreateEBR(fdisk.Path, int32(sizeBytes), extendedPartition, fdisk.Name)
    if err != nil {
        return msg, fmt.Errorf("error creando el nuevo EBR: %s", err)
    }

    // Crear la nueva partición lógica justo después del nuevo EBR
    logicalPartitionStart := newEBRStart + int32(binary.Size(ebr))
    var logicalPartition PARTITION
    logicalPartition.CreatePartition(int(logicalPartitionStart), sizeBytes, "L", fdisk.Fit, fdisk.Name)

    // Serializar la nueva partición lógica
    msg, err = logicalPartition.SerializePartition(fdisk.Path, logicalPartitionStart)
    if err != nil {
        return msg, fmt.Errorf("error escribiendo la partición lógica: %s", err)
    }

    return "Partición lógica creada exitosamente", nil
}